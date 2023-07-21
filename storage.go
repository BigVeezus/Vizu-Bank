package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostgresStore, error){
	connStr := "user=postgres dbname=postgres sslmode=disable"
	db,err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	} 

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}


func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(100) NOT NULL,
		last_name varchar(100) NOT NULL,
		acc_number serial NOT NULL,
		balance serial NOT NULL,
		created_at timestamp NOT NULL
	)`

	_,err := s.db.Exec(query)
	return err
}


func (s *PostgresStore) CreateAccount(acc *Account) error {
	
	query := `insert into account 
	(first_name, last_name, acc_number, balance, created_at)
	values ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.AccNumber,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	// fmt.Println(resp)

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil,nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {

	query := `SELECT * FROM account`

	rows, err := s.db.Query(
		query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName, 
			&account.LastName, 
			&account.AccNumber, 
			&account.Balance, 
			&account.CreatedAt);
		

		if err != nil {
			return nil,err
		}

		accounts = append(accounts, account) 
	}

	return accounts, nil
}
