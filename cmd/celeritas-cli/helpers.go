package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
)

func setup(arg1 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}

		cel.RootPath = path
		cel.DB.DataType = os.Getenv("DATABASE_TYPE")
	}
}

func getDSN() string {
	dbType := cel.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}

		return dsn
	}

	return "mysql://" + cel.BuildDSN()
}

func showHelp() {
	color.Yellow(`Available commands:

	help					- show the help commands
	version 				- print application version
	make migration <name>	- create two new up and down migrations in the migration folder 
	make auth				- creates and runs migrations for authenticates tables, and creates models and middleware
	make handler <name>		- creates a new stub handler in handlers directory
	make model <name>		- creates a new stub model in data directory
	make session 			- creates a table in the database as a session store
	make mail <name>		- creates two starter templates in the mail directory
	migrate					- runs all up migrations that have no been run previously
	migrate down			- reverses the most recent migration
	migrate reset 			- runs all down migrations in reverse order, and then all up migrations
   `)
}
