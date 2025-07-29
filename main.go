package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq" // Importing pq for PostgreSQL driver
	"github.com/tbirddv/gator/internal/config"
	"github.com/tbirddv/gator/internal/database"
)

type state struct {
	config  *config.Config
	queries *database.Queries
	args    []string
}

func main() {
	configData, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", configData.DBURL)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()
	queries := database.New(db)
	command := strings.ToLower(os.Args[1])
	args := make([]string, 0)
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	state := &state{
		config:  configData,
		queries: queries,
		args:    args,
	}
	commands := CommandInit(state)
	if command, exists := commands[command]; exists {
		err = command.Execute()
		if err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands:")
		for _, cmd := range commands {
			fmt.Printf("- %s: %s\n", cmd.Name, cmd.Description)
		}
		os.Exit(1)
	}
}
