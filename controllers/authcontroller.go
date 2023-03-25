package controllers

import (
	"errors"
	"go-auth/config"
	"go-auth/entities"
	"go-auth/libraries"
	"go-auth/models"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string
	Password string
}

var userModel = models.NewUserModel()
var validation = libraries.NewValidation()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			temp, _ := template.ParseFiles("views/index.html")
			temp.Execute(w, nil)
		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		errorMessages := validation.Struct(UserInput)

		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
			}

			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)
		} else {
			var user entities.User
			userModel.Where(&user, "username", UserInput.Username)

			var message error
			if user.Username == "" {
				// не нашел в базе
				message = errors.New("неправильное имя пользователя или пароль")
			} else {
				// проверка пароля
				errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
				if errPassword != nil {
					message = errors.New("неправильное имя пользователя или пароль")
				}
			}
			if message != nil {
				data := map[string]interface{}{
					"error": message,
				}

				temp, _ := template.ParseFiles("views/login.html")
				temp.Execute(w, data)
			} else {
				session, _ := config.Store.Get(r, config.SESSION_ID)

				session.Values["loggedIn"] = true
				session.Values["email"] = user.Email
				session.Values["username"] = user.Username

				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		user := entities.User{
			Email:     r.Form.Get("email"),
			Username:  r.Form.Get("username"),
			Password:  r.Form.Get("password"),
			Cpassword: r.Form.Get("cpassword"),
		}

		errorMessage := validation.Struct(user)

		if errorMessage != nil {
			data := map[string]interface{}{
				"validation": errorMessage,
				"user":       user,
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {
			// insert database
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashPassword)

			userModel.Create(user)

			data := map[string]interface{}{
				"message": "регистрация прошла успешно",
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)

		}

	}
}
