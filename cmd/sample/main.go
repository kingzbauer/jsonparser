package main

import (
	"fmt"
	"os"

	"github.com/kingzbauer/jsonparser"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Expected filename")
		os.Exit(1)
	}

	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if err := jsonparser.Parse(src); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
