package main

import (
	"fmt"

	_ "github.com/1024XEngineer/xinfra/server/docs"
	"github.com/1024XEngineer/xinfra/server/internal/router"
)

func main() {
	r := router.Setup()

	fmt.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
