package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/pub", "./pub")
	router.Static("/sub", "./sub")
	router.Static("/assets", "./assets")
	//router.Static("/", "./_src")
	// Listen and serve on 0.0.0.0:8080
	router.Run(":58080")
}
