# GaliBot

My Telegram bot, implemented on GO using [tgbotapi](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5)

## CURRENTLY REFACTORING

<!-- TODO: Usage -->

## Process

1. Init db
1. Init bot
1. Register commands
1. Start recurrent processes (eg. IP)
1. Receive updates
    1. Get user from db (create if new)
        * user id
        * Permissions
        * Current mode
    1. Handle
        * Message
            * Command
                1. /cancel previous mode ? (eg. BotFather)
                1. Apply
            * Text
                1. Mode handler (ignore by default)
        * CallbackQuery
            * Mode
            * Data
            1. Send data to the mode handler
        * Other types later
    1. Cache any db updates and flush them now
    1. React to update (by replying or editing a message)
