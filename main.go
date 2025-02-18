package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("charge eth and memo for test or product")

	router := gin.Default()

	router.GET("/chargetest", chargetestHandler)
	router.GET("/chargepro", chargeproHandler)

	router.Run(":8003")
}

// 获取所有用户
func chargetestHandler(c *gin.Context) {
	// 从请求中获取参数
	addr := c.Query("addr")
	if addr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	// 构造并执行第一条命令
	cmd1 := exec.Command(
		"./mefs",
		"transfer",
		"eth",
		"--endPoint=https://testchain.metamemo.one:24180",
		"--sk=0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4",
		addr,
		"1000gwei",
	)

	fmt.Println("charging eth..")

	// 执行第一条命令并捕获输出
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		log.Printf("Error executing first mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute first mefs command"})
		return
	}

	fmt.Printf("result:\n%s\n", string(output1))

	// 构造并执行第二条命令
	cmd2 := exec.Command(
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

	fmt.Println("charging memo..")

	// 执行第二条命令并捕获输出
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		log.Printf("Error executing second mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute second mefs command"})
		return
	}

	fmt.Printf("result:\n%s\n", string(output2))

	// 返回两条命令的输出
	c.JSON(http.StatusOK, gin.H{
		"output1": string(output1),
		"output2": string(output2),
	})
}

// 获取所有用户
func chargeproHandler(c *gin.Context) {
	// 从请求中获取参数
	addr := c.Query("addr")
	if addr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	// 构造并执行第一条命令
	cmd1 := exec.Command(
		"./mefs",
		"transfer",
		"eth",
		"--endPoint=https://chain.metamemo.one:8501",
		"--sk=0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4",
		addr,
		"1000gwei",
	)

	fmt.Println("charging eth..")

	// 执行第一条命令并捕获输出
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		log.Printf("Error executing first mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute first mefs command"})
		return
	}

	fmt.Printf("result:\n%s\n", string(output1))

	// 构造并执行第二条命令
	cmd2 := exec.Command(
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

	fmt.Println("charging memo..")

	// 执行第二条命令并捕获输出
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		log.Printf("Error executing second mefs command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute second mefs command"})
		return
	}

	fmt.Printf("result:\n%s\n", string(output2))

	// 返回两条命令的输出
	c.JSON(http.StatusOK, gin.H{
		"output1": string(output1),
		"output2": string(output2),
	})
}
