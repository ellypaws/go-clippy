# Clippy-Go

This Discord bot is designed to assist users in retrieving and managing function details across different platforms, such as Excel and Google Sheets. With an integrated points system, it encourages community collaboration and problem-solving.
It uses [DiscordGo](https://github.com/bwmarrin/discordgo) for the bot, and [Bingo](https://github.com/nokusukun/bingo) for the database

### Setup

To run the bot, you need the Bot token from the [Discord Developer Portal](https://discord.com/developers/applications). Once you have the token, you can run the bot using the following command:

- **token** (`token`): The bot access token.
    ```
    token YOUR_BOT_ACCESS_TOKEN
    ```
- **guild** (`guild`): The Guild ID. If not provided, the bot registers commands globally.
    ```
    guild GUILD_ID
    ```
- **rmcmd** (`rmcmd`): This flag determines whether to remove all commands after shutting down. This is particularly useful for cleanup. By default, this is set to `false`.
    ```
    rmcmd=true
    ```
  
### Example:
```
go-clippy.exe token YOUR_BOT_ACCESS_TOKEN
```

The token can also be set by setting the `BOT_TOKEN` environment variable. You can also use the .env file in the working directory.

| Platform       | Command                    |
|----------------|----------------------------|
| Powershell     | `$env:token = "BOT_TOKEN"` |
| Bash           | `export token=BOT_TOKEN`   |
| Command Prompt | `set token=BOT_TOKEN`      |

### Running the bot

Once the bot is running, you can use the following commands to interact with it:

   | Key                            | Action           |
   |--------------------------------|------------------|
   | <kbd>Ctrl</kbd> + <kbd>C</kbd> | Quit the program |

1. [x] TUI to be implemented

### Slash Commands

The bot supports a variety of slash commands to ease the functionality:

#### 1. `/lookup`

Fetches details of a function for a given platform.

Usage:
```
/lookup function:FUNCTION_NAME platform:PLATFORM_NAME
```

Example:
```
/lookup function:sum platform:Excel
```

#### 2. `/search`

Similar to `/lookup`, but has an autocomplete feature. You can also search by description.

Usage:
```
/search function:FUNCTION_NAME description:DESCRIPTION platform:PLATFORM_NAME
```

Example:
```
/search function:sum platform:Excel
/search description:filter platform:Google Sheets
```

#### 3. `/solved`

Marks a thread as finished and optionally awards points to a user for solving it.

Usage:
```
/solved @User #channel
```

Parameters:
- `@User` (optional): Mention the user you wish to award points to.
- `#channel` (optional): Mention the channel where the thread you wish to mark as finished is located. If the command is run inside a thread, the bot assumes it's the current thread.

Example:
```
/solved @JohnDoe #excel-help
```

### Todo

1. [x] Actually run the bot
2. [x] Implement [bubbletea](https://github.com/charmbracelet/bubbletea.git) TUI
3. [ ] Leaderboard
4. [x] Points system


<img src="https://go.dev/images/gophers/ladder.svg" width="48" alt="Go Gopher climbing a ladder." align="right">

### Contribution

---

Feel free to contribute to the bot's development or suggest new features. Your feedback is always appreciated.