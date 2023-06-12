package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/todopeer/backend/orm"
)

func main() {
	db, err := gorm.Open("sqlite3", orm.DBPATH)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	userORM := orm.NewUserORM(db)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter email:")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Println("Enter name:")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	user := &orm.User{
		Email: email,
		Name:  &name,
	}

	fmt.Println("Enter password:")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	err = user.SetPassword(strings.TrimSpace(password))
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	err = userORM.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("User %s created successfully\n", email)
}
