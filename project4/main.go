package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("README.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(file.Name())
}
