package main

import (
	"os"
	"shuttle/databases"
	_ "github.com/lib/pq"
	"github.com/fatih/color"
	"github.com/pressly/goose/v3"
)

func main() {
	color.Yellow("Connecting to Database...")

	db, err := databases.PostgresConnection()
	if err != nil {
		color.Red("Failed to connect to PostgreSQL:", err)
		os.Exit(1)
	}
	defer db.Close()

	sqlDB := db.DB

	err = os.Chdir("databases")
	if err != nil {
		color.Red("Failed to change directory to databases:", err)
		os.Exit(1)
	}

	color.Yellow("Running migrations...")

	// Get param from command line
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "down":
			color.Yellow("Rolling back migrations...")
			err = goose.Down(sqlDB, "./migrations")
		case "redo":
			color.Yellow("Redoing migrations...")
			err = goose.Redo(sqlDB, "./migrations")
		case "status":
			color.Yellow("Checking migration status...")
			err = goose.Status(sqlDB, "./migrations")
		case "downall":
			color.Yellow("Rolling back all migrations...")
			err = goose.DownTo(sqlDB, "./migrations", 0)
		case "version":
			color.Yellow("Checking migration version...")
			err = goose.Version(sqlDB, "./migrations")
		default:
			color.Red("Unknown command:", os.Args[1])
			os.Exit(1)
		}
	} else {
		color.Yellow("Running migrations...")
		err = goose.Up(sqlDB, "./migrations")
	}

	if err != nil {
		color.Red("Migration failed:", err)
		os.Exit(1)
	}

	color.Green("Migration successful")
}