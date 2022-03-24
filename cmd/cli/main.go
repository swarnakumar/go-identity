package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/db/users"
)

func getUserFromCli(reader *bufio.Reader, ctx context.Context, client *db.Client) string {
	fmt.Printf("Enter User's Email: ")
	emailBytes, _, _ := reader.ReadLine()
	email := strings.TrimSpace(string(emailBytes))
	exists := client.Users.CheckEmailExists(ctx, email)
	if !exists {
		return email
	}

	fmt.Printf("The user already exists. Do you want to enter a new email? (y/n) ")
	answerBytes, _, _ := reader.ReadLine()
	if (strings.ToLower(string(answerBytes)) == "n") || (strings.ToLower(string(answerBytes)) == "no") {
		os.Exit(1)
	}

	return getUserFromCli(reader, ctx, client)
}

func getPwdFromCli(reader *bufio.Reader) string {
	fmt.Printf("Enter User's Password: ")
	passwordBytes, _, _ := reader.ReadLine()
	password := strings.TrimSpace(string(passwordBytes))

	ok, _ := users.CheckPassword(password)
	if !ok {
		fmt.Printf("The password is too simple. Try a more complex one.")
		return getPwdFromCli(reader)
	}

	return password
}

func getSuperUserFlagFromCli(reader *bufio.Reader) bool {
	fmt.Printf("Is this going to be a superuser (y/N): ")
	superUserResp, _, _ := reader.ReadLine()
	isSuperUser := strings.ToLower(strings.TrimSpace(string(superUserResp)))

	return strings.HasPrefix(isSuperUser, "y")

}

func CreateFromCli() {
	ctx := context.Background()
	client := db.New(ctx)

	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	email := getUserFromCli(reader, ctx, client)
	password := getPwdFromCli(reader)
	isAdmin := getSuperUserFlagFromCli(reader)

	createdBy := "cli"
	user, err := client.Users.Create(ctx, email, password, isAdmin, &createdBy)

	if err != nil {
		fmt.Printf("failed creating user: %s", err)
	} else {
		log.Println("user was created: ", user)
	}
}

func main() {
	CreateFromCli()
}
