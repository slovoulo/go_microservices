package helpers

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)




//A function to create an encrypted hash of users password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err!=nil{
        log.Println("An error occcured encrypting password")
        return "",err
    }
    log.Println("Password encrypted successfully")
    return string(bytes), nil
}






// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func  WriteErrorJSON(w http.ResponseWriter,  errorCode  int, errorMessage string)  {

    //Error code example: http.StatusInternalServerError
    w.WriteHeader(errorCode)
    w.Write([]byte(errorMessage))





	
}
