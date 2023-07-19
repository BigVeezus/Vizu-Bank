package main

import "math/rand"

type Account struct {
	ID 			int		`json:"id"`
	FirstName   string	`json:"firstName"`
	LastName 	string	`json:"lastName"`
	AccNumber 	int64	`json:"accountNum"`
	Balance 	int64	`json:"accountBalance"`
}

func NewAccount(firstName, lastName string) *Account{
	return &Account{
		ID: rand.Intn(100000),
		FirstName: firstName,
		LastName: lastName,
		AccNumber: int64(rand.Intn(10000000)),
		
	}
}