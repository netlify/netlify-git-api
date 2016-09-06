package cli

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/netlify/netlify-git-api/api"
	"github.com/netlify/netlify-git-api/repo"
	"github.com/netlify/netlify-git-api/userdb"
	"github.com/pborman/uuid"
)

type userWrapper struct {
	dbUser *userdb.User
}

func (u *userWrapper) Name() string {
	return u.dbUser.Name
}

func (u *userWrapper) Email() string {
	return u.dbUser.Email
}

func (u *userWrapper) HasPermission(_ string, _ string) bool {
	return true
}

type resolver struct {
	db       *userdb.UserDB
	repoPath string
	tokens   map[string]string
	sync     bool
}

func (r *resolver) GetRepo(req *http.Request) (*repo.Repo, error) {
	authHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return nil, nil
	}
	id, ok := r.tokens[parts[1]]
	if !ok {
		return nil, nil
	}

	user := r.db.Get(id)
	if user == nil {
		return nil, nil
	}

	currentRepo, err := repo.Open(&userWrapper{dbUser: user}, r.repoPath, r.sync)
	if err != nil {
		panic(fmt.Sprintf("Unable to open git repository in %v: %v", r.repoPath, err))
	}

	return currentRepo, nil
}

func (r *resolver) Authenticate(email, pw string) (string, error) {
	user := r.db.LookupByEmail(email)
	if user == nil {
		return "", nil
	}
	if !user.Authenticate(pw) {
		return "", nil
	}

	token := uuid.New()

	r.tokens[token] = user.ID

	return token, nil
}

// Serve starts a new REST API server
func Serve(dbPath, host, port string, sync bool) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working dir: %v\n", err)
	}

	userDB, err := userdb.Read(dbPath)
	if err != nil {
		log.Fatalf("Error reading user db %v: %v\n", dbPath, err)
	}
	if len(userDB.Users) == 0 {
		log.Fatalf("Error - no users in user db %v\n", dbPath)
	}

	resolver := &resolver{db: userDB, repoPath: cwd, tokens: map[string]string{}, sync: sync}

	api := api.NewAPI(resolver)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), api))
}
