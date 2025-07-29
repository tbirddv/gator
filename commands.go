package main

type Command struct {
	Name        string
	Description string
	Execute     func() error
}

func CommandInit(state *state) map[string]Command {
	commands := make(map[string]Command)

	commands["login"] = Command{
		Name:        "login",
		Description: "Set the current user for gator",
		Execute: func() error {
			return HandleLogin(state)
		},
	}

	commands["help"] = Command{
		Name:        "help",
		Description: "Display help information for commands",
		Execute: func() error {
			return HandleHelp(commands, state)
		},
	}

	commands["register"] = Command{
		Name:        "register",
		Description: "Register a new user with gator",
		Execute: func() error {
			return HandleRegister(state)
		},
	}

	commands["reset"] = Command{
		Name:        "reset",
		Description: "Reset all users in the database",
		Execute: func() error {
			return HandleResetUsers(state)
		},
	}

	commands["users"] = Command{
		Name:        "users",
		Description: "List all users in the database",
		Execute: func() error {
			return HandleGetUsers(state)
		},
	}

	return commands
}
