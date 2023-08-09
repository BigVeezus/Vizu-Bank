package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string ) *Account{
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new acc number ->", acc.AccNumber)

	return acc;

}

func seedAccounts(s Storage){
	seedAccount(s, "anthony", "bam", "hello")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")

	flag.Parse();

	store, err  := NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	
	if *seed {
	// seed stuff
	fmt.Println("seeding the db")
	seedAccounts(store);
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}