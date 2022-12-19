package main

import (
  //"database/sql"
  //_ "github.com/go-sql-driver/mysql"
)

type DBUser struct {
  ChatId int64;
  Name string;
  Permissions string; 
};

func getUser(chatId int64) (DBUser, error) {
  var user DBUser;
  err := DB.QueryRow("SELECT * FROM users WHERE chatid = ?", chatId).Scan(&user.ChatId, &user.Name, &user.Permissions);
  return user, err;
}
