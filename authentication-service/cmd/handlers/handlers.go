package handlers

import (
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slovojoe/authentication-service/cmd/database"
	"github.com/slovojoe/authentication-service/cmd/models"
	"github.com/slovojoe/authentication-service/helpers"
	"golang.org/x/crypto/bcrypt"
)

type User models.User
type u *User

type UserInput struct{
	Email string
	Password string
	Username string
}

type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
	payload:=jsonResponse{
		Error: false,
		Message: "Welcome to Go Auth Services",
	}

	out,_:=json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)

	

}

func CreateUser(w http.ResponseWriter, r *http.Request){
	// get the body of the  POST request
	// unmarshal this into a new Userinput struct
	// append this to the Users array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userinput UserInput
	json.Unmarshal(reqBody, &userinput)

	//extract the password from the struct and hash it
	hashedPassword,hasherr:=helpers.HashPassword(userinput.Password)

	//If there is no error hashing the password
	if(hasherr==nil){
	//Create a new user based on the hashed password
	var newUser = User{Email:userinput.Email, PasswordHash: hashedPassword,Username: userinput.Username}
	if result := database.Db.Create(&newUser); result.Error != nil {
		fmt.Println(result.Error)
	}
	json.NewEncoder(w).Encode(newUser)

	}

}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
// func (u *UserInput) PasswordMatches(plainText string) (bool, error) {
// 	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
// 			// invalid password
// 			return false, nil
// 		default:
// 			return false, err
// 		}
// 	}

// 	return true, nil
// }
 
//Fetch a user based on their email
func FetchUserEmail(email string)(*User,error){
	

	var singleUser User
	// SELECT * FROM users WHERE email = email;
	result:=database.Db.First(&singleUser, "email = ?", email )
	if result.Error !=nil{
		log.Printf("An error occured fetching user %s",result.Error)
		return nil,result.Error
	}
	return  &singleUser,nil

	
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func  PasswordMatches(plainText string,  PasswordHash  string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(PasswordHash), []byte(plainText))
	if err != nil {
		
			return false, err
		
	}

	return true, nil
}

//A function that authenticates user provided credentials based on records stored in db
func  Authenticate(w http.ResponseWriter, r *http.Request) {
	// get the body of the  POST request
	// unmarshal this into a new Userinput struct

	reqBody, _ := ioutil.ReadAll(r.Body)
	var userinput UserInput
	json.Unmarshal(reqBody, &userinput)

	

	// validate the user against the database
	user, err := FetchUserEmail(userinput.Email)
	if err != nil {
		log.Println("Invalid email")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid email \n"))
		w.Write([]byte(err.Error()))
		return
	}

	valid, perr := PasswordMatches(userinput.Password,user.PasswordHash)
	if perr != nil || !valid {
		log.Println("Invalid password")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid password \n"))
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

	log.Println("user authenticated proceeding to log ")

	//Call logger service to log authentication
	logErr:=logRequest(w,"Authentication log",fmt.Sprintf("%s logged in", userinput.Email))
	if logErr!=nil{
		helpers.WriteErrorJSON(w,http.StatusBadRequest,err.Error())
	}

	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	json.NewEncoder(w).Encode(payload)
}

func logRequest(w http.ResponseWriter, name string, data string) error{
	log.Println("Auth starting log service")
	type entry struct{
		Name string	`json:"name"`
		Data string	`json:"data"`
	}

	newEntry:=entry{Name: name,Data: data}

	jsonData,_:= json.Marshal(newEntry)
	log.Printf("Auth json data as %s",jsonData)
	logServiceURL := "http://logger-service:9090/writelog"

	request,err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err!=nil{
		log.Printf("Auth service couldnt connect to the logger service, %s",err)
		//http.Error(w, "Internal server error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Auth service couldnt connect to the logger service"))
		return err

	}
	log.Println("Auth service connected to log service")

	client:=&http.Client{}
	_, clierr := client.Do(request)
	if err!=nil{
		log.Printf("Log service client error: , %s",err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Log service client error"))
		return clierr

	}
	return nil

}