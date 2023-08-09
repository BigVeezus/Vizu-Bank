package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResp struct{
	AccNumber	int64	`json:"accountNum"`
	Token		string		`json:"token"`
}

type LoginRequest struct{
	AccNumber	int64	`json:"accountNum"`
	Password	string	`json:"password"`
}

type CreateAccountRequest struct{
	FirstName	string	`json:"firstName"`
	LastName 	string	`json:"lastName"`
	Password	string	`json:"password"`
}

type Account struct {
	ID 			int		`json:"id"`
	FirstName   string	`json:"firstName"`
	LastName 	string	`json:"lastName"`
	AccNumber 	int64	`json:"accountNum"`
	EncryptedPassword	string `json:"."`
	Balance 	int64	`json:"accountBalance"`
	CreatedAt	time.Time	`json:"createdAt"`
}

type TransferRequest struct {
	ToAccountNum	int		`json:"toAccountNum"`
	Amount		int 	`json:"amount"`
	Description	string 	`json:"description"`
}

func (a *Account) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(firstName, lastName, password string) (*Account, error){
	encpw,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil
	}

	return &Account{
		FirstName: firstName,
		LastName: lastName,
		EncryptedPassword: string(encpw),
		AccNumber: RandomAza(10000000),
		CreatedAt: time.Now().UTC(),
		
	}, nil
}

func RandomAza(num int) int64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// Call the resulting `rand.Rand` just like the
	// functions on the `rand` package.
	
	aza := r1.Intn(num)
	return int64(aza)
}