package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)
	
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func (http.ResponseWriter, *http.Request) error  

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w,r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()} )
		}
	}
}

type APIServer struct {
	listenAddr string
	store		Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: 		store,
	}
}

func (s *APIServer) Run(){
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount) )
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID) )

	log.Println("API listening on: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if(r.Method == "GET"){
		return s.handleGetAccount(w,r)
	}
	if(r.Method == "POST"){
		return s.handleCreateAccount(w,r)
	}
	if(r.Method == "DELETE"){
		return s.handleDeleteAccount(w,r)
	}
	
	return fmt.Errorf("method not allowed! %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts,err := s.store.GetAccounts() 
	if err != nil {
		
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}


func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	// account := NewAccount("Elvis", "Osuji")

	return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
createAccReq := new(CreateAccountRequest)
if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
	return err
}
	account := NewAccount(createAccReq.FirstName,createAccReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}