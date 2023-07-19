package main

import (
	"fmt"
	"log"
)


func main() {
	store, err := NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)
	// server := NewAPIServer(":3000")
	// server.Run()
}