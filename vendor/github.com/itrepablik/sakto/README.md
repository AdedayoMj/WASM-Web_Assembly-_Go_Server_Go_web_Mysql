![sakto](https://user-images.githubusercontent.com/58651329/80955641-2fd1a980-8e32-11ea-91b3-f83263a9b15b.png)
The **sakto** package is the simplified common input validators for your Go projects.

# Installation
```
go get -u github.com/itrepablik/sakto
```

# Usage
These are some of the examples on how you can use this package.
```
package main

import (
	"fmt"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/sakto"
)

func main() {
	// These are some few examples that you can use for your Go project.
	// HashAndSalt usage: hash any plain text password
	plainTextPassword := "hEllo_World!"
	hsPassword, err := sakto.HashAndSalt([]byte(plainTextPassword))
	if err != nil {
		itrlog.Fatal(err)
	}
	fmt.Println("hsPassword: ", hsPassword)

	// CheckPasswordHash usage: compare the plain text password vs hashed password stored from your database.
	isPassHashMatch, err := sakto.CheckPasswordHash(plainTextPassword, hsPassword)
	if isPassHashMatch {
		fmt.Println("Password match, login successful!")
	} else {
		fmt.Println("Invalid password, please try again!")
	}
}
```

# Subscribe to Maharlikans Code Youtube Channel:
Please consider subscribing to my Youtube Channel to recognize my work on this package. Thank you for your support!
https://www.youtube.com/channel/UCdAVUmldU9Jn2VntuQChHqQ/

# License
Code is distributed under MIT license, feel free to use it in your proprietary projects as well.
