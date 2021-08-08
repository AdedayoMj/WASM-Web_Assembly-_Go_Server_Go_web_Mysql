package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	// "net/url"
	"strconv"
	"strings"
	"syscall/js"
	"time"
	// _ "github.com/go-sql-driver/mysql"
)

// login exposed function to JS interface{}
func login(this js.Value, args []js.Value) interface{} {
	// Get the JS objects here
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get document object", "error");`)
	}
	username := jsDoc.Call("getElementById", "username")
	if !username.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get username object", "error");`)
	}
	password := jsDoc.Call("getElementById", "password")
	if !password.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get password object", "error");`)
	}

	isSiteKeepMe := jsDoc.Call("getElementById", "isSiteKeepMe")
	if !isSiteKeepMe.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get remember option", "error");`)
	}

	csrfToken := jsDoc.Call("getElementById", "csrfToken")
	if !csrfToken.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get CSRF token", "error");`)
	}

	//convert input fields to a proper field types
	var pUserName string = username.Get("value").String()
	var pPassword string = password.Get("value").String()
	pIsSiteKeepMe, _ := strconv.ParseBool(isSiteKeepMe.Get("value").String())
	var pCSRFToken string = csrfToken.Get("value").String()

	if len(strings.TrimSpace(pUserName)) == 0 {
		return js.Global().Call("eval", `Swal.fire("Username is Required!", "Please enter your username", "error");`)
	}

	if len(strings.TrimSpace(pPassword)) == 0 {
		return js.Global().Call("eval", `Swal.fire("Password is Required!", "Please enter your password", "error");`)
	}

	//compose JSON post payload to the API endpoint
	payLoad := map[string]interface{}{
		"username":     pUserName,
		"password":     pPassword,
		"isSiteKeepMe": pIsSiteKeepMe,
	}

	bytesRepresentation, err := json.Marshal(payLoad)
	if err != nil {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Something went wrong with your login credentials", "error");`)
	}

	//HTTP new request
	siteHost := "http://127.0.0.1:8081/api/v1/user/login"
	client := &http.Client{}
	// data := url.Values{}

	req, err := http.NewRequest("POST", fmt.Sprintf(siteHost), bytes.NewBuffer(bytesRepresentation))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-TOKEN", pCSRFToken)
	if err != nil {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Something went wrong with your connection, please try again", "error");`)
	}

	var isSuccess bool = false

	//Get the response from the htpp.NewRequest post method from a channel

	c1 := make(chan map[string]interface{}, 1)
	var result2 map[string]interface{}
	go func() {
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&result2)
		c1 <- result2 //send response data to our channel name 'c'
	}()
	var result map[string]interface{}
	go func() interface{} {
		for {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(time.Second * 10)
				timeout <- true
			}()

			select {
			case result = <-c1:
				// fmt.Println("result 1: ", result)
				i := result["IsSuccess"] //You must return with JSON format value of either 'true' or 'false' only.
				mStatus := fmt.Sprint(i)
				isSuccess, _ = strconv.ParseBool(mStatus)
				msg := ""
				if !isSuccess {
					msg = `Swal.fire("Login Unsuccessful", "Either your username or password is wrong.", "error");`
				} else {
					msg = `Swal.fire("Login Successful", "Your account has been verified and it's successfully logged-in.", "success");`
				}
				return APIResponse(isSuccess, msg)
				// return testResult(isSuccess)
				break
			case <-timeout:
			}
		}
	}()

	return nil
}

func register(this js.Value, args []js.Value) interface{} {
	// Get the JS objects here
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get document object", "error");`)
	}

	username := jsDoc.Call("getElementById", "username")
	if !username.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get username object", "error");`)
	}
	email := jsDoc.Call("getElementById", "email")
	if !email.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get email object", "error");`)
	}
	password := jsDoc.Call("getElementById", "password")
	if !password.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get password object", "error");`)
	}
	confirmpassword := jsDoc.Call("getElementById", "confirmpassword")
	if !confirmpassword.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get confirm password object", "error");`)
	}

	chkAgree := jsDoc.Call("getElementById", "chkAgree")
	if !chkAgree.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get terms and service object", "error");`)
	}

	//Get the csrf token value from the client side
	csrfToken := jsDoc.Call("getElementById", "csrfToken")
	if !csrfToken.Truthy() {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Unable to get password", "error");`)
	}

	//convert input fields to a proper field types
	var pCSRFToken string = csrfToken.Get("value").String()
	var pUserName string = username.Get("value").String()
	var pEmail string = email.Get("value").String()
	var pPassword string = password.Get("value").String()
	var pConfirmPassword string = confirmpassword.Get("value").String()
	var pTerms bool = chkAgree.Get("checked").Bool()

	//compose JSON post payload to the API endpoint
	payLoad := map[string]interface{}{
		"username":        pUserName,
		"email":           pEmail,
		"password":        pPassword,
		"confirmpassword": pConfirmPassword,
		"terms":           fmt.Sprint(pTerms),
		"isActive":        "false",
	}

	bytesRepresentation, err := json.Marshal(payLoad)
	if err != nil {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Something went wrong with your user's registration", "error");`)
	}

	//HTTP new request
	siteHost := "http://127.0.0.1:8081/api/v1/user/register"
	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf(siteHost), bytes.NewBuffer(bytesRepresentation))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-TOKEN", pCSRFToken)
	if err != nil {
		return js.Global().Call("eval", `Swal.fire("Oops!, Error", "Something went wrong with your connection, please try again", "error");`)
	}

	var isSuccess bool = false

	//Get the response from the htpp.NewRequest post method from a channel

	c1 := make(chan map[string]interface{}, 1)
	var result2 map[string]interface{}
	go func() {
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&result2)
		c1 <- result2 //send response data to our channel name 'c'
	}()
	var result map[string]interface{}
	go func() interface{} {
		for {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(time.Second * 10)
				timeout <- true
			}()

			select {
			case result = <-c1:
				fmt.Println("result 1: ", result)
				i := result["IsSuccess"] //You must return with JSON format value of either 'true' or 'false' only.
				mStatus := fmt.Sprint(i)
				isSuccess, _ = strconv.ParseBool(mStatus)
				alertTitle := fmt.Sprint(result["AlertTitle"])
				alertMsg := fmt.Sprint(result["AlertMsg"])
				alertType := fmt.Sprint(result["AlertType"])

				msg := ""
				if !isSuccess {
					msg = `Swal.fire("` + alertTitle + `", "` + alertMsg + `", "` + alertType + `");` //error
				} else {
					msg = `Swal.fire("` + alertTitle + `", "` + alertMsg + `", "` + alertType + `");` //success
				}
				return APIResponse(isSuccess, msg)
				// return testResult(isSuccess)
				break
			case <-timeout:
			}
		}
	}()

	return nil
}

// func testResult(isSuccess bool) interface{} {
// 	if !isSuccess {
// 		return js.Global().Call("eval", `Swal.fire("Login Unsuccessful", "Either your username or password is wrong.")`)
// 	}
// 	return js.Global().Call("eval", `Swal.fire("Login Successful", "Your account has been verified and it's successfully logged-in.")`)
// }

func APIResponse(isSuccess bool, msg string) interface{} {
	if !isSuccess {
		return js.Global().Call("eval", msg)
	}
	return js.Global().Call("eval", msg)
}

// func getRequestPOST(r *http.Request, client *http.Client, c chan map[string]interface{}) {
// 	go func() {
// 		resp, _ := client.Do(r)
// 		resp.Body.Close()
// 		var result map[string]interface{}
// 		json.NewDecoder(resp.Body).Decode(&result)
// 		c <- result //send response data to our channel name 'c'
// 		close(c)
// 	}()
// }
func exposeGoFuncJS() {
	// Start exposing the following Go functions to JS client side
	js.Global().Set("login", js.FuncOf(login))
	js.Global().Set("register", js.FuncOf(register))

}

func main() {
	fmt.Println("Welcome to Adegboye WASM Go Testing")
	c := make(chan bool, 1)

	// Initializes all your exposable Go's functions to JS
	exposeGoFuncJS()

	<-c
}
