package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	email := flag.String("email", "", "admin email")
	phone := flag.String("phone", "", "e.g. +2348030000000")
	first := flag.String("first", "", "first name")
	last := flag.String("last", "", "last name")

	flag.Parse()

	if *email == "" || *phone == "" {
		log.Fatal("email and phone are required")
	}

	if *first == "" || *last == "" {
		log.Fatal("first and last name are required")
	}

	fmt.Print("Password (at least 12 characters): ")

	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if len(password) < 12 {
		log.Fatal("password must be at least 12 characters")
	}

	admin := Admin{
		Email:    strings.ToLower(*email),
		Phone:    *phone,
		First:    *first,
		Last:     *last,
		Password: password,
	}

	err := save(admin)
	if err != nil {
		log.Fatal("could not create admin: ", err)
	}

	fmt.Println("Platform admin created. Use staff login.")
}
