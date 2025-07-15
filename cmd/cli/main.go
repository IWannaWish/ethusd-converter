package main

import (
	"fmt"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage ethusd-converter <ethereum_address>")
		log.Fatal("Error: no address provided")
	}

	address := os.Args[1]
	fmt.Println("Address provided:", address)

	// временно вместо настоящего вызова:
	output := core.GetAssets(address)
	log.Println(output)
}
