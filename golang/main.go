package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	helper "github.com/FurkanSamaraz/IsEmpty"
	emailControl "github.com/FurkanSamaraz/emailControl"
	_ "github.com/lib/pq"
)

var uname, pwd, email string

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "172754"
	dbname   = "postgres"
)

type UserModel struct {
	Id       int
	Username string
	Password string
	Email    string
}
type LoginModel struct {
	Username string
	Password string
	Email    string
}

func openConnention() *sql.DB {

	psq := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psq)

	if err != nil {
		helper.IsEmpty(err.Error())
	}
	err = db.Ping()
	if err != nil {
		helper.IsEmpty(err.Error())
	}

	return db
}

var userr UserModel

func login(w http.ResponseWriter, r *http.Request) {

	var people []UserModel
	db := openConnention()
	r.ParseForm()
	uname := r.FormValue("username")
	pwd := r.FormValue("password")
	email := r.FormValue("email")
	rows, _ := db.Query("SELECT * FROM userr")
	for rows.Next() {
		rows.Scan(&userr.Id, &userr.Username, &userr.Password, &userr.Email)
		people = append(people, userr)

	}
	if uname == userr.Username && pwd == userr.Password && email == userr.Email {
		fmt.Fprintf(w, "Login successful\n")
		fmt.Fprintln(w, "Hello", uname)
		peopleByte, _ := json.MarshalIndent(userr, "", "\t")
		w.Write(peopleByte)

	}

	defer db.Close()
}

func bosSignup(w http.ResponseWriter, r *http.Request) {
	unameCheck := helper.IsEmpty(uname)
	pwdCheck := helper.IsEmpty(pwd)
	mailCheck := helper.IsEmpty(email)

	if unameCheck || pwdCheck || mailCheck {

	} else {
		fmt.Fprintf(w, "Error Empty! \n")
	}
}

func isValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
func register(w http.ResponseWriter, r *http.Request) {
	db := openConnention()
	r.ParseForm()
	var lg LoginModel

	lg.Username = r.FormValue("username")
	lg.Password = r.FormValue("password")
	lg.Email = r.FormValue("email")

	rows, _ := db.Query("SELECT * FROM userr")
	for rows.Next() {
		rows.Scan(&userr.Id, &userr.Username, &userr.Password, &userr.Email)

	}

	uCheck := strings.Contains(lg.Password, lg.Username)
	eCheck := strings.Contains(lg.Password, lg.Email)

	if uCheck == true || eCheck == true {
		fmt.Fprintf(w, "Password must not contain username or email.")
	} else {

		if isValid(lg.Password) != true {
			fmt.Fprintf(w, "Use special characters, numbers, upper and lower case letters in the password.")
		} else {

			if lg.Username == "" || lg.Password == "" || lg.Email == "" {
				fmt.Fprintf(w, "cannot be empty")
			} else {
				if userr.Username == lg.Username {
					fmt.Fprintf(w, "username is used")
				} else {
					if emailControl.CheckEmail(lg.Email) == true {

						db.Exec("INSERT INTO userr(username,password,eposta) VALUES($1,$2,$3)", lg.Username, lg.Password, lg.Email)

						peopleByte, _ := json.MarshalIndent(lg, "", "\t")

						w.Header().Set("Content-Type", "application/json")

						w.Write(peopleByte)

						defer db.Close()
						bosSignup(w, r)
						db.Close()
					} else {
						fmt.Fprintln(w, "record failed error email!! ", uname)
					}
				}
			}
		}
	}
}
func update(w http.ResponseWriter, r *http.Request) {

	db := openConnention()
	r.ParseForm()

	userr.Username = r.FormValue("username")
	userr.Password = r.FormValue("password")
	userr.Email = r.FormValue("email")
	db.Exec("UPDATE userr SET username=$1,password=$2,eposta=$3 WHERE id=$4 ", userr.Username, userr.Password, userr.Email, userr.Id)

	peopleByte, _ := json.MarshalIndent(userr, "", "\t")

	w.Header().Set("Content-Type", "application/json")

	w.Write(peopleByte)

	defer db.Close()
	bosSignup(w, r)
	db.Close()
}
func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/update", update)
	http.ListenAndServe(":8080", mux)
}
