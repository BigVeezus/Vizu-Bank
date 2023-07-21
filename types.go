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

func NewAccount(firstName, lastName string) *Account{
	return &Account{
		FirstName: firstName,
		LastName: lastName,
		AccNumber: int64(rand.Intn(10000000)),
		CreatedAt: time.Now().UTC(),
		
	}
}