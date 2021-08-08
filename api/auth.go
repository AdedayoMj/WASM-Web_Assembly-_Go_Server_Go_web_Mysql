package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gowebapp/config"
	"gowebapp/shiga"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/sulat"
	"github.com/itrepablik/tago"
)

// SGC initialize this variable globally sulat.SendGridConfig{}
var SGC = sulat.SGC{}

func init() {
	// Set shiga SMTP option
	SGC = shiga.SetSendGridAPI(sulat.SGC{
		SendGridAPIKey: config.SendGridAPIKey,
	})
}

// AuthRouters are the collection of all URLs for the Auth App.
func AuthRouters(r *mux.Router) {
	r.HandleFunc("/api/v1/user/login", LoginUserEndpoint).Methods("POST")
	r.HandleFunc("/api/v1/user/register", SignUpUserEndpoint).Methods("POST")
	r.HandleFunc("/login", Login).Methods("GET")
	r.HandleFunc("/signup", Signup).Methods("GET")
	r.HandleFunc("/about", About).Methods("GET")
	r.HandleFunc("/chat", Chat).Methods("GET")
}

// Login function is to render the login page.
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(config.SiteRootTemplate+"front/login.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "Login - " + config.SiteShortName,
		"PageMetaDesc": config.SiteShortName + " account, sign in to access your account.",
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
		"Settings":     config.SiteSettings,
	}
	tmpl.Execute(w, data)
}

// Register function is to render the registration page.
func Signup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(config.SiteRootTemplate+"front/signup.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "Signup - " + config.SiteShortName,
		"PageMetaDesc": config.SiteShortName + " create new user account.",
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
		"Settings":     config.SiteSettings,
	}
	tmpl.Execute(w, data)
}

// About function is to render the home page.
func About(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(config.SiteRootTemplate+"front/about.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "About - " + config.SiteShortName,
		"PageMetaDesc": config.SiteShortName + " profile, find out about the developer.",
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
		"Settings":     config.SiteSettings,
	}
	tmpl.Execute(w, data)
}

// Chat function is to render the homepage page.
func Chat(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(config.SiteRootTemplate+"demo/react/chat.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "Chat - " + config.SiteShortName,
		"PageMetaDesc": config.SiteShortName + " profile, find out chat the developer.",
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
		"Settings":     config.SiteSettings,
	}
	tmpl.Execute(w, data)
}

type jsonResponse struct {
	IsSuccess  string `json:"isSuccess"`
	AlertTitle string `json:"alertTitle"`
	AlertMsg   string `json:"alertMsg"`
	AlertType  string `json:"alertType"`
}

// LoginUserEndpoint is to validate the user's login credential
func LoginUserEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, errBody := ioutil.ReadAll(r.Body)
	if errBody != nil {
		itrlog.Error(errBody)
		panic(errBody.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	userName := strings.TrimSpace(keyVal["username"])
	password := keyVal["password"]
	isSiteKeepMe, _ := strconv.ParseBool(keyVal["isSiteKeepMe"])

	itrlog.Info("userName: ", userName)
	itrlog.Info("password: ", password)
	itrlog.Info("isSiteKeepMe: ", isSiteKeepMe)

	// Check if username is empty
	if len(strings.TrimSpace(userName)) == 0 {
		w.Write([]byte(`{ "IsSuccess": "false", "AlertTitle": "Username is Required", "AlertMsg": "Please enter your username.", "AlertType": "error" }`))
		return
	}

	// Check if password is empty
	if len(strings.TrimSpace(password)) == 0 {
		w.Write([]byte(`{ "IsSuccess": "false", "AlertTitle": "Password is Required", "AlertMsg": "Please enter your password.", "AlertType": "error" }`))
		return
	}

	// Set the cookie expiry in days.
	expDays := "1" // default to expire in 1 day.
	if isSiteKeepMe == true {
		expDays = config.UserCookieExp
	}

	// Encrypt the username value to store it from the user's cookie.
	encryptedUserName, err := tago.Encrypt(userName, config.MyEncryptDecryptSK)
	if err != nil {
		itrlog.Error(err)
	}

	// // Initialize the database connection
	// dbCon, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/carsier?parseTime=true&charset=utf8mb4,utf8")
	// if err != nil {
	// 	itrlog.Error(err)
	// }
	// defer dbCon.Close()

	// // Now, insert the new user's information here
	// ins, err := dbCon.Prepare("INSERT INTO shiga_user (username, password, email, first_name, " +
	// 	"middle_name, last_name, suffix, is_superuser, is_admin, date_joined, is_active) VALUES" +
	// 	"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	// if err != nil {
	// 	itrlog.Error(err)
	// }

	// // Pass on all the parameter values here
	// ins.Exec(userName, password, "politz@live.com", "P O L", "D.", "Peligro", "Jr.", 1, 0, time.Now(), 0)
	// defer ins.Close()

	w.Write([]byte(`{ "isSuccess": "true", "alertTitle": "Login Successful", "alertMsg": "Your account has been verified and it's successfully logged-in.",
		"alertType": "success", "redirectTo": "` + config.SiteBaseURL + `dashboard", "eUsr": "` + encryptedUserName + `", "expDays": "` + expDays + `" }`))
}

// SignUpUserEndpoint is to register the user's login credential
func SignUpUserEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, errBody := ioutil.ReadAll(r.Body)
	if errBody != nil {
		itrlog.Error(errBody)
		panic(errBody.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	userName := strings.TrimSpace(keyVal["username"])
	email := keyVal["email"]
	password := keyVal["password"]
	confirmPassword := keyVal["confirmpassword"]
	terms, _ := strconv.ParseBool(strings.TrimSpace(keyVal["terms"]))
	isActive, _ := strconv.ParseBool(strings.TrimSpace(keyVal["isActive"]))

	// Open the MySQL DSB Connection
	dbShiga, err := sql.Open("mysql", DBConStr(""))
	if err != nil {
		itrlog.Error(err)
	}
	defer dbShiga.Close()

	// New shiga user information
	newUser := shiga.User{
		UserName:    userName,
		Password:    password,
		Email:       email,
		IsSuperUser: false,
		IsAdmin:     false,
		IsActive:    isActive,
	}

	// Email confirmation
	newUserEmailConfirmation := shiga.EmailConfig{
		From:            config.SiteEmail,
		FromAlias:       "Support Team",
		To:              email,
		Subject:         "Activate your " + config.SiteShortName + " account",
		DefaultTemplate: shiga.EmailFormatNewUser,
		// EmailConfirmationURL: confirmURL,
	}

	// Check if "isActive" is true, then we don't send an email confirmation to activate the new user's account.
	if isActive {
		_, err := shiga.CreateUser(dbShiga, newUser, newUserEmailConfirmation, confirmPassword, terms)
		if err != nil {
			itrlog.Error(err)
			w.Write([]byte(`{ "IsSuccess": "false", "AlertTitle": "New User Creation Failed!", 
			"AlertMsg": "` + err.Error() + `", "AlertType": "error"}`))
			return
		}

		// Response back to the user about the succcessful user's registration
		w.Write([]byte(`{ "IsSuccess": "true", "AlertTitle": "New User", 
		"AlertMsg": "You've successfully created a new ` + config.SiteShortName + `'s account.",
		"AlertType": "success" }`))
		// "RedirectURL": "` + shiga.YB.BaseURL + `account_activation_sent"
		// http.Redirect(w, r, "https://127.0.0.1/welcome", http.StatusOK )
	} else {
		// Insert the new user's registration here
		_, err := shiga.CreateUser(dbShiga, newUser, newUserEmailConfirmation, confirmPassword, terms)
		if err != nil {
			itrlog.Error(err)
			w.Write([]byte(`{ "IsSuccess": "false", "AlertTitle": "New User Creation Failed!", 
			"AlertMsg": "` + err.Error() + `", "AlertType": "error" }`))
			return
		}

		// Response back to the user about the succcessful user's registration with auto-redirect to a successful page
		w.Write([]byte(`{ "IsSuccess": "true", "AlertTitle": "Registration is Successful", 
		"AlertMsg": "You've successfully created your new ` + config.SiteShortName + `'s user account ",
		"AlertType": "success", "RedirectTo": "` + config.SiteBaseURL + `" }`))
	}
}

// DBConStr is the connection string for your database
func DBConStr(dbName string) string {
	db := dbName
	if len(strings.TrimSpace(dbName)) == 0 {
		db = config.DBName
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8", config.DBUserName,
		config.DBPassword, config.DBHostName, db)
}
