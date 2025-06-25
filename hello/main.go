package main

import (
	"fmt"

	"rsc.io/quote/v4"
)

func main() {
	name := ""
	fmt.Print("Please my friend, tell me your name: ")
	fmt.Scanln(&name)

	age := 0
	fmt.Print("Now your age: ")
	fmt.Scanln(&age)

	var age_cond string
	if age >= 18 {
		age_cond = "old"
	} else {
		age_cond = "young"
	}

	content := "Hello my %s friend! It's really nice see you here :-)\n\n"
	content += quote.Glass() + "\n"
	content += "I don't know what that means, but here we are %s"

	fmt.Println(fmt.Sprintf(content, age_cond, name))
}
