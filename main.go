package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"gritt/ride"
)

func main() {
	addr := flag.String("addr", "localhost:4502", "Dyalog RIDE address")
	flag.Parse()

	fmt.Printf("Connecting to %s...\n", *addr)
	client, err := ride.Connect(*addr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("Connected!")

	// Test with 1+1
	output, err := client.Execute("1+1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1+1 = %s", strings.Join(output, ""))
}
