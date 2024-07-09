package main

import (
	"citadel/config"
	"citadel/database"
	"fmt"
	"os"
	"strings"

	"github.com/caesar-rocks/orm"
)

const (
	possibleArgs = "serve migrations:run migrations:rollback migrations:reset db:seed list:routes"
)

func main() {
	args := os.Args
	if len(args) != 2 || !strings.Contains(possibleArgs, args[1]) {
		fmt.Printf("Usage: go run caesar.go [%s]\n", possibleArgs)
		os.Exit(1)
	}

	switch args[1] {
	case "serve":
		env := config.ProvideEnvironmentVariables()
		config.ProvideApp(env).Run()
	case "list:routes":
		listRoutes()
	case "migrations:run":
		getDB().Migrate(database.GetMigrations())
	case "migrations:rollback":
		getDB().Rollback(database.GetMigrations())
	case "migrations:reset":
		getDB().Reset(database.GetMigrations())
	case "db:seed":
		getDB().Seed()
	}
}

func getDB() *orm.Database {
	env := config.ProvideEnvironmentVariables()
	db := config.ProvideDatabase(env)
	return db
}

func listRoutes() {
	env := config.ProvideEnvironmentVariables()
	app := config.ProvideApp(env)
	router := app.RetrieveRouter()

	fmt.Println("Method\tPattern")
	fmt.Println("------\t-------")

	for _, route := range router.Routes {
		fmt.Printf("%s\t%s\n", route.Method, route.Pattern)
	}
}
