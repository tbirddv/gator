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

	commands["agg"] = Command{
		Name:        "agg",
		Description: "Aggregate RSS feeds, Usage: agg <time_between_requests>",
		Execute: func() error {
			return HandleAgg(state)
		},
	}

	commands["addfeed"] = Command{
		Name:        "addfeed",
		Description: "Add a new RSS feed",
		Execute: func() error {
			return HandleCreateFeed(state)
		},
	}

	commands["feeds"] = Command{
		Name:        "feeds",
		Description: "List all RSS feeds",
		Execute: func() error {
			return HandleGetFeeds(state)
		},
	}

	commands["follow"] = Command{
		Name:        "follow",
		Description: "Follow an RSS feed",
		Execute: func() error {
			return HandleFollow(state)
		},
	}

	commands["following"] = Command{
		Name:        "following",
		Description: "List all RSS feeds followed by the current user",
		Execute: func() error {
			return HandleGetFollows(state)
		},
	}

	commands["unfollow"] = Command{
		Name:        "unfollow",
		Description: "Unfollow an RSS feed for the current user",
		Execute: func() error {
			return HandleUnfollow(state)
		},
	}

	commands["browse"] = Command{
		Name:        "browse",
		Description: "Browse posts from Current User's followed feeds. Usage: browse [Number of Posts to Browse]",
		Execute: func() error {
			return HandleBrowse(state)
		},
	}

	return commands
}
