package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "access denied"})
}

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtIjo4NjQ0MDQ1LCJleHBpcmVzQXQiOjE1MTYyMzkwMjJ9.GI7VSt6Iidf-FYstYgPGDNb-EM13Y085LA_6Wb2Wyxk"
	
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Calling JWT middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)

		if err != nil {
			permissionDenied(w)
			return 
		}
		if !token.Valid {
			permissionDenied(w)
			return 
		}

		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
	
		if err != nil {
			WriteJSON(w, http.StatusNotFound, ApiError{Error: "invalid id"})
			return 
		}

		account, err := s.GetAccountByID(id)
		

		if err != nil {
			WriteJSON(w, http.StatusNotFound, ApiError{Error: "invalid id"})
			return 
		}

		claims := token.Claims.(jwt.MapClaims)
		// panic(reflect.TypeOf(claims["accountNum"]))
		// fmt.Println(claims["accountNum"])
		// fmt.Println(account.AccNumber)


		if account.AccNumber != int64((claims["accountNum"]).(float64)) {
			 permissionDenied(w)
			 return
		}


	handlerFunc(w,r)
	}
}

func createJWT(acc *Account) (string, error) {

	// Create the Claims
claims := &jwt.MapClaims{
	"expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
	"accountNum":  int64(acc.AccNumber),
}
secret := os.Getenv("JWT_SECRET")
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 return token.SignedString([]byte(secret))


}

func validateJWT(tokenString string) (*jwt.Token ,error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

type apiFunc func (http.ResponseWriter, *http.Request) error  

type ApiError struct {
	Error string	`json:"error"`
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

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount) )
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetorDeleteAccountByID) , s.store)) 
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer) )

	log.Println("API listening on: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}
 // 4211409
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST"{
		return fmt.Errorf("method not allowed! %s", r.Method)
	}
	
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc,err := s.store.GetAccountByNum(int(req.AccNumber))
	if err != nil {
		return err;
	}

	if !acc.ValidatePassword(req.Password) {
		return fmt.Errorf("Invalid user details")
	}

	token, err := createJWT(acc)
	if err != nil {
		return nil
	}

	resp := LoginResp{
		Token: token,
		AccNumber: acc.AccNumber,
	}

	fmt.Printf("%+v\n", acc)

	return WriteJSON(w, http.StatusOK, resp);
}
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if(r.Method == "GET"){
		return s.handleGetAccount(w,r)
	}
	if(r.Method == "POST"){
		return s.handleCreateAccount(w,r)
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


func (s *APIServer) handleGetorDeleteAccountByID(w http.ResponseWriter, r *http.Request) error {
	
	if r.Method == "GET"  {
	
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return fmt.Errorf("invalid Id: %s", idStr)
	}

	account, err := s.store.GetAccountByID(id)

	if err != nil {
		return err
	}


	return WriteJSON(w, http.StatusOK, account)

	}
	if r.Method == "DELETE" {
		
	 return s.handleDeleteAccount(w,r)
	}

	return fmt.Errorf("method: %s not valid", r.Method)
}


func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
createAccReq := new(CreateAccountRequest)
if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
	return err
}
	account,err := NewAccount(createAccReq.FirstName,createAccReq.LastName, createAccReq.Password)
	if err != nil {
		return err;
	}
	
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}



	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return fmt.Errorf("invalid Id: %s", idStr)
	}

	 err = s.store.DeleteAccount(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "POST" {

	
	transferReq := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()



	return WriteJSON(w, http.StatusOK, transferReq)

}
return fmt.Errorf("method: %s not valid", r.Method)
}