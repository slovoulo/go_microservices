package models

import "gorm.io/gorm"

type User struct{
    gorm.Model

    Username string
    Email string
    PasswordHash string
    //Documents []Document
}