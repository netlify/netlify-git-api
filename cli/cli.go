package cli

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app    = kingpin.New("netlify-git-api", "Get a REST API for a Git repository")
	dbPath = app.Flag("db", "File path to the user db").Default(".users.yml").String()

	serve = app.Command("serve", "Start a local Git API server")
	port  = serve.Flag("port", "Port to listen to").Short('p').Default("8080").String()
	host  = serve.Flag("host", "IP to bind to").Short('h').Default("127.0.0.1").IP()
	//sync  = serve.Flag("sync", "Push and pull to the origin remote ()").Short('s').Bool()

	users = app.Command("users", "List users")

	usersList        = users.Command("list", "List all users")
	usersAdd         = users.Command("add", "Add a new user")
	usersAddName     = usersAdd.Flag("name", "Name of the new user").String()
	usersAddEmail    = usersAdd.Flag("email", "Email of new user").String()
	usersAddPassword = usersAdd.Flag("password", "Password of new user").String()
	usersDel         = users.Command("del", "Remove a user")
	usersDelEmail    = usersDel.Arg("email", "Email of the user").String()
)

// Run the cli tool
func Run() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serve.FullCommand():
		fmt.Printf("Starting server on %v:%v\n", *host, *port)
		Serve(*dbPath, host.String(), *port, false)
	case usersList.FullCommand():
		ListUsers(*dbPath)
	case usersAdd.FullCommand():
		AddUser(*dbPath, *usersAddEmail, *usersAddName, *usersAddPassword)
	case usersDel.FullCommand():
		DeleteUser(*dbPath, *usersDelEmail)
	}
}
