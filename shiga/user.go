package shiga

import (
	"database/sql"
	"errors"

	"gowebapp/config"
	"strings"
	"time"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/sakto"
	"github.com/itrepablik/sulat"
)

// CreateUser add a new user to the users collection
func CreateUser(dbCon *sql.DB, u User, e EmailConfig, confirmPassword string, terms bool) (int64, error) {
	// Check if username is empty
	if len(strings.TrimSpace(u.UserName)) == 0 {
		return 0, errors.New("Username is Required")
	}

	// Check if the username is available or not
	if !IsUserNameExist(dbCon, u.UserName) {
		return 0, errors.New("Username already exist, please try again")
	}

	// Check if email is empty
	if len(strings.TrimSpace(u.Email)) == 0 {
		return 0, errors.New("Email is Required")
	}

	// Check if email address is valid or not
	if !sakto.IsEmailValid(u.Email) {
		return 0, errors.New("Invalid Email Address, please try again")
	}
	// Check if the email address is available or not
	if !IsUserEmailExist(dbCon, u.Email) {
		return 0, errors.New("Email already exist, please try again")
	}

	// Check if password is empty
	if len(strings.TrimSpace(u.Password)) == 0 {
		return 0, errors.New("Password is Required")
	}

	//Match both passwords
	if strings.TrimSpace(confirmPassword) != strings.TrimSpace(u.Password) {
		return 0, errors.New("Passwords didn't match, please try again")
	}

	// Check if Terms of service has been checked
	if !terms {
		return 0, errors.New("Terms of Service is Required, By joining " + config.SiteShortName + ", you're agreeing to our terms and conditions.")
	}

	if !u.IsActive {
		// Check if from email address is empty
		if len(strings.TrimSpace(e.From)) == 0 {
			return 0, errors.New("From email address is required")
		}

		// Check if to email address is empty
		if len(strings.TrimSpace(e.To)) == 0 {
			return 0, errors.New("To email address is required")
		}

		// Check if subject is empty
		emailSubject := "Activate your new account account"
		if len(strings.TrimSpace(e.Subject)) > 0 {
			emailSubject = e.Subject
		}

		// Check if HTML Header template has been customized
		emailHTMLHeader := ShigaHTMLHeader // default to Shiga HTML Header
		if len(strings.TrimSpace(e.CustomizeHeaderTemplate)) > 0 {
			emailHTMLHeader = e.CustomizeHeaderTemplate
		}

		// Check if HTML Body template has been customized
		emailHTMLBody := NewUserActivation(e.EmailConfirmationURL, u.UserName, e.SiteName, e.SiteSupportEmail) // default to Yabi HTML Body
		if len(strings.TrimSpace(e.CustomizeBodyTemplate)) > 0 {
			emailHTMLBody = e.CustomizeBodyTemplate
		}

		// Check if HTML Footer template has been customized
		emailHTMLFooter := ShigaHTMLFooter // default to Yabi HTML Footer
		if len(strings.TrimSpace(e.CustomizeFooterTemplate)) > 0 {
			emailHTMLFooter = e.CustomizeFooterTemplate
		}

		// Send an email confirmation now, prepare the HTML email content first
		mailOpt := &sulat.SendMail{
			Subject: emailSubject,
			From:    sulat.NewEmail(e.FromAlias, e.From),
			To:      sulat.NewEmail(e.ToAlias, e.To),
			CC:      sulat.NewEmail(e.CCAlias, e.CC),
			BCC:     sulat.NewEmail(e.BCCAlias, e.BCC),
		}
		htmlContent, err := sulat.SetHTML(&sulat.EmailHTMLFormat{
			IsFullHTML: false,
			HTMLHeader: emailHTMLHeader,
			HTMLBody:   emailHTMLBody,
			HTMLFooter: emailHTMLFooter,
		})
		_, err = sulat.SendEmailSG(mailOpt, htmlContent, &SGC)
		if err != nil {
			itrlog.Error("SendGrid error: ", err)
		}
	}

	// Hash and salt your plain text password
	hsPassword, err := sakto.HashAndSalt([]byte(u.Password))
	if err != nil {
		return 0, err
	}

	// Insert the new user's information here
	ins, err := dbCon.Prepare("INSERT INTO " + ShigaUser + " (username, password, email, first_name, " +
		"middle_name, last_name, suffix, is_superuser, is_admin, date_joined, is_active) VALUES" +
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}

	// Pass on all the parameter values here
	ins.Exec(u.UserName, hsPassword, u.Email, u.FirstName, u.MiddleName, u.LastName, u.Suffix, u.IsSuperUser,
		u.IsAdmin, time.Now(), u.IsActive)

	// Get the lastest inserted id
	lid, err := GetLastInsertedID(dbCon, "id", ShigaUser)
	defer ins.Close()

	return lid, nil

}

// GetLastInsertedID gets the latest inserted id for any specified table and it's auto_increment field
func GetLastInsertedID(dbCon *sql.DB, autoIDFieldName, tableName string) (int64, error) {
	var id int64 = 0
	err := dbCon.QueryRow("SELECT " + autoIDFieldName + " FROM " + tableName + " ORDER BY " + autoIDFieldName + " DESC LIMIT 1").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// IsUserNameExist check from the user's collection if it's existed or not, we don't allow to have a
// duplicate username, it must be a unique value
func IsUserNameExist(dbCon *sql.DB, userName string) bool {
	var id int64 = 0
	err := dbCon.QueryRow("SELECT id FROM "+ShigaUser+" WHERE username = ?", userName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true // returned no rows, the username is not found from the yabi table
		}
		return false
	}
	return false
}

// IsUserEmailExist check from the user's collection if it's existed or not, we don't allow to have a
// duplicate email, it must be a unique value
func IsUserEmailExist(dbCon *sql.DB, email string) bool {
	var id int64 = 0
	err := dbCon.QueryRow("SELECT id FROM "+ShigaUser+" WHERE email = ?", email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true // returned no rows, the email is not found from the yabi table
		}
		return false
	}
	return false
}
