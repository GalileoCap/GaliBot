module github.com/GalileoCap/tgGenBot/plugins/ip

go 1.19

replace github.com/GalileoCap/tgGenBot/api => ../../../tgGenBot/api

require (
	github.com/GalileoCap/tgGenBot/api v0.0.0-00010101000000-000000000000
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
)

require (
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/muesli/cache2go v0.0.0-20221011235721-518229cd8021 // indirect
)
