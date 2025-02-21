package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/rockiecn/charge-svr/kv"
)

var g_db *kv.KV

func main() {
	fmt.Println("charge eth and memo for test or product")

	// init db
	db, err := kv.NewBadgerDb("./db")
	if err != nil {
		panic(err)
	}
	g_db = db

	// init gin
	router := gin.Default()

	// handlers
	router.GET("/charge", chargeHandler)
	router.GET("/query", queryHandler)

	// run server
	router.Run("0.0.0.0:8003")
}

// charge eth and memo
func chargeHandler(c *gin.Context) {
	// get chain
	chain := c.Query("chain")
	if chain != "test" && chain != "product" {
		log.Printf("Error chain must set to test or product in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "chain must be set to test or product"})
		return
	}

	// get address
	addr := c.Query("addr")
	if addr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	// check exist in db
	// key = addr+chain
	exist, err := g_db.Exists(addr + chain)
	if err != nil {
		log.Printf("Error check exist: %v, address: %s, chain: %s", err, addr, chain)
		c.JSON(http.StatusInternalServerError, gin.H{"call db.exists error": err, "address": addr, "chain": chain})
		return
	}

	// address already charged
	if exist {
		log.Printf("Error address %s has been charged in chain %s", addr, chain)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "address has been charged", "address": addr, "chain": chain})
		return
	}

	// 构造并执行第一条命令
	var cmd1 *exec.Cmd
	var cmd2 *exec.Cmd
	switch chain {
	case "test":
		cmd1 = exec.Command(
			"./mefs",
			"transfer",
			"eth",
			"--endPoint=https://testchain.metamemo.one:24180",
			"--sk=0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4",
			addr,
			"1000gwei",
		)
		cmd2 = exec.Command(
			"./mefs",
			"transfer",
			"memo",
			"--endPoint=https://testchain.metamemo.one:24180",
			"--instanceAddr=0x66e92976548a7C959DE81B7396Fe619a8C99E05c",
			"--sk=0ba0403047a1c0cb08c31c0f5432b2df93b032164d83c27377e5da31cdacc0d0",
			"--version=2",
			addr,
			"1",
		)
	case "product":
		cmd1 = exec.Command(
			"./mefs",
			"transfer",
			"eth",
			"--endPoint=https://chain.metamemo.one:8501",
			"--sk=0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4",
			addr,
			"1000gwei",
		)
		cmd2 = exec.Command(
			"./mefs",
			"transfer",
			"memo",
			"--endPoint=https://chain.metamemo.one:8501",
			"--instanceAddr=0xbd16029A7126C91ED42E9157dc7BADD2B3a81189",
			"--sk=0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4",
			"--version=2",
			addr,
			"1",
		)
	}

	fmt.Println("charging eth..")

	// 执行第一条命令并捕获输出
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		log.Printf("Error executing first mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute first mefs command"})
		return
	}

	fmt.Printf("charge eth result:\n%s\n", string(output1))

	fmt.Println("charging memo..")

	// 执行第二条命令并捕获输出
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		log.Printf("Error executing second mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute second mefs command"})
		return
	}

	fmt.Printf("charge memo result:\n%s\n", string(output2))

	// set this address to be charged
	g_db.Set(addr+chain, "true")

	// 返回两条命令的输出
	c.JSON(http.StatusOK, gin.H{
		"output1": string(output1),
		"output2": string(output2),
	})
}

// query size of a provider
func queryHandler(c *gin.Context) {
	// mefs provider id
	id := c.Query("id")
	if id == "" {
		log.Printf("no id in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be set in request"})
		return
	}

	// get chain
	chain := c.Query("chain")
	if chain != "test" && chain != "product" {
		log.Printf("Error chain must set to test or product in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "chain must be set to test or product"})
		return
	}

	// 构造命令
	cmd := exec.Command(
		"./contractsv2",
		"get",
		fmt.Sprintf("--ep=%s", chain),
		"settleInfo",
		id,
		"0",
	)

	fmt.Println("querying provider size..")

	// 执行查询命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing first mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute first mefs command"})
		return
	}

	fmt.Printf("query size result:\n%s\n", string(output))

	strOut := string(output)
	var strSize string

	// 使用正则表达式匹配 "Size": 后面的数字
	re := regexp.MustCompile(`"Size":\s*(\d+)`)
	matches := re.FindStringSubmatch(strOut)

	if len(matches) > 1 {
		strSize = matches[1] // 匹配到的数字部分

	} else {
		fmt.Println("Size not found")
	}

	// 返回两条命令的输出
	c.JSON(http.StatusOK, gin.H{
		"size": strSize,
	})
}
