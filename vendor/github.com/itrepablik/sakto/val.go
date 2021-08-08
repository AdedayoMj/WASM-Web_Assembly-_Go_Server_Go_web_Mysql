package sakto

import (
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/crypto/bcrypt"
)

const _defaultTimeZone = "UTC"

// HashAndSalt is to hash user's password using bycrypt in Go.
func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares the two passwords between the stored hash and the raw plain text password.
func CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetCurDT gets the current system date and time with specified local timezone, e.g timeZone = "Asia/Manila".
func GetCurDT(t time.Time, timeZone string) time.Time {
	loc, _ := time.LoadLocation(timeZone)
	return time.Now().In(loc)
}

// IsUserNameValid ensure that only allowed special characters like @, ., +, -, and _.
func IsUserNameValid(userName string) bool {
	var verifyUserName = regexp.MustCompile(`^[a-zA-Z0-9_@.+-]*$`).MatchString
	return verifyUserName(userName)
}

// IsEmailValid ensure the email is a valid email address.
func IsEmailValid(email string) bool {
	const emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	var verifyEmail = regexp.MustCompile(emailRegexString).MatchString
	return verifyEmail(email)
}

// IsURLValid ensure that only a valid URL pattern will be accepted.
func IsURLValid(strURL string) bool {
	u, err := url.Parse(strURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// TimeParseToString parses the time.Time with specific date and time format.
func TimeParseToString(strDT, dateTimeFormat string) (string, error) {
	i, err := strconv.ParseInt(strDT, 10, 64)
	if err != nil {
		return "", err
	}
	tm := time.Unix(i, 0)
	fCreatedDate, _ := time.Parse(dateTimeFormat, tm.Format(dateTimeFormat))
	return fCreatedDate.Format(dateTimeFormat), nil
}

// ParseFloatToString parses any float value, options: bitSize either (32 or 64) with decimal value.
// Ideal for any pricing or financial amount values.
func ParseFloatToString(strInt string, bitSize, decimals int) (string, error) {
	numVal, err := strconv.ParseFloat(strInt, bitSize)
	if err != nil {
		return "", err
	}
	return humanize.CommafWithDigits(numVal, decimals), nil
}

// TrimQ removes the trailing double quotes ("") from a string.
func TrimQ(cleanStr string) string {
	cleanStr = strings.TrimSpace(cleanStr)
	if len(cleanStr) > 0 && cleanStr[0] == '"' {
		cleanStr = cleanStr[1:]
	}
	if len(cleanStr) > 0 && cleanStr[len(cleanStr)-1] == '"' {
		cleanStr = cleanStr[:len(cleanStr)-1]
	}
	return cleanStr
}

// LocalNow gets the current local time with localize timezone
func LocalNow(tz string) time.Time {
	if len(strings.TrimSpace(tz)) == 0 {
		tz = _defaultTimeZone
	}
	return GetCurDT(time.Now(), tz)
}

// IsFileExist checks if the file existed or not
func IsFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
