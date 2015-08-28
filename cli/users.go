package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/netlify/netlify-git-api/userdb"
	"golang.org/x/crypto/ssh/terminal"
)

func promptString(name string) (string, error) {
	fmt.Println(fmt.Sprintf("%v: ", name))
	var input string
	fmt.Scanln(&input)
	return input, nil
}

func promptPassword() (string, error) {
	fmt.Println("Password: ")
	bytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	return fmt.Sprintf("%s", bytes), err
}

// ListUsers lists all users
func ListUsers(dbPath string) {
	db, err := userdb.Read(dbPath)
	if err != nil {
		log.Fatalf("Error: failed to read db %v: %v\n", dbPath, err)
	}

	if len(db.Users) == 0 {
		log.Printf("No users found in %v\n", dbPath)
	} else {
		for _, user := range db.Users {
			log.Printf("%v: %v <%v>\n", user.ID, user.Name, user.Email)
		}
	}
}

// AddUser ads a new user
func AddUser(dbPath, email, name, pw string) {
	var err error

	db, err := userdb.Read(dbPath)
	if err != nil {
		log.Fatalf("Error: failed to read db %v: %v\n", dbPath, err)
	}

	if email == "" {
		email, err = promptString("Email")
		if err != nil || email == "" {
			log.Fatalf("Error: Could not read email: %v\n", err)
		}
	}

	if name == "" {
		name, err = promptString("Name")
		if err != nil || name == "" {
			log.Fatalf("Error: Could not read name: %v\n", err)
		}
	}

	if pw == "" {
		pw, err = promptPassword()
		if err != nil {
			panic(err)
		}
	}

	if _, err := db.Add(email, name, pw); err != nil {
		log.Fatalf("Error: Could not add user: %v: %v", email, err)
	}

	if err := db.Write(); err != nil {
		log.Fatalf("Error: Could not write db %v: %v", dbPath, err)
	}

	log.Printf("User %v added", email)
}

// DeleteUser deletes an existing user
func DeleteUser(dbPath, email string) {
	db, err := userdb.Read(dbPath)
	if err != nil {
		log.Fatalf("Error: failed to read db %v: %v\n", dbPath, err)
	}

	if email == "" {
		email, err = promptString("Email")
		if err != nil || email == "" {
			log.Fatalf("Error: Could not read email: %v\n", err)
		}
	}

	db.Delete(email)

	if err := db.Write(); err != nil {
		log.Fatalf("Error: Could not write db %v: %v", dbPath, err)
	}

	log.Printf("User %v deleted", email)
}
