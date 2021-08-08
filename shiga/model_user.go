package shiga

import "time"

//the table name for the shiga_user table
const ShigaUser = "shiga_user"

//User model collections for the user information
type User struct {
	ID          int64     `json:"id"`
	UserName    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	Suffix      string    `json:"suffix"`
	IsSuperUser bool      `json:"is_superuser"`
	IsAdmin     bool      `json:"is_admin"`
	LastLogin   time.Time `json:"last_login"`
	DateJoined  time.Time `json:"date_joined"`
	IsActive    bool      `json:"is_active"`
}
