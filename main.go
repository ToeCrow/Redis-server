package main

import (
	"fmt"
	"os"

	"github.com/thokro/redis-server/internal/server"
)

func main() {
	if err := server.Run(":6379"); err != nil {
		fmt.Println("Error at start:", err)
		os.Exit(1)
	}
}
