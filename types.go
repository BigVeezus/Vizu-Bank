package main

import (
	"math/rand"
	"time"
)

type CreateAccountRequest struct{
	FirstName	string	`json:"firstName"`
	LastName 	string	`json:"lastName"`
}

type Account struct {
	ID 			int		`json:"id"`
	FirstName   string	`json:"firstName"`
	LastName 	string	`json:"lastName"`
	AccNumber 	int64	`json:"accountNum"`
	Balance 	int64	`json:"accountBalance"`
	CreatedAt	time.Time	`json:"createdAt"`
}

type TransferRequest struct {
	ToAccountNum	int		`json:"toAccountNum"`
	Amount		int 	`json:"amount"`
	Description	string 	`json:"description"`
}

func NewAccount(firstName, lastName string) *Account{
	return &Account{
		FirstName: firstName,
		LastName: lastName,
		AccNumber: RandomAza(10000000),
		CreatedAt: time.Now().UTC(),
		
	}
}

func RandomAza(num int) int64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// Call the resulting `rand.Rand` just like the
	// functions on the `rand` package.
	
	aza := r1.Intn(num)
	return int64(aza)
}