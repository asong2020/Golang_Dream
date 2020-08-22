package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `validate:"min=6,max=10"`
	Age      uint8  `validate:"gte=6,lte=10"`
	Sex      string `validate:"oneof=female male"`
}

func main() {
	validate := validator.New()

	user1 := User{Username: "asong", Age: 11, Sex: "null"}
	err := validate.Struct(user1)
	if err != nil {
		fmt.Println(err)
	}

	user2 := User{Username: "asong111", Age: 8, Sex: "male"}
	err = validate.Struct(user2)
	if err != nil {
		fmt.Println(err)
	}

}
