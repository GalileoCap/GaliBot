/*
 * TODO: Use memory-only database for testing
 */
package main

import (
  _ "github.com/mattn/go-sqlite3"
  "database/sql"
	"github.com/muesli/cache2go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "time"
	"log"
);

type ModeT int;
const (
  Normal ModeT = iota
);

type User struct {
  // This block is stored in the db
  ID int64;
  FirstName string;
  LastName string `json:",omitempty"`;
  UserName string `json:",omitempty"`;
  Permissions int8; // 0: admin, 1: tester, 2: user, 3: block //TODO: Enum

  // This block only exists on the cache
  Mode string;
  FOODDOEntry *FOODDOEntry;
};
var DB *sql.DB;
var Users *cache2go.CacheTable;

func dbInit(dbPath string) {
  var err error;
  DB, err = sql.Open("sqlite3", dbPath);
  if err != nil {
    log.Fatalf("[dbInit] %v", err);
  }
  //TODO: Check/Create tables

  Users = cache2go.Cache("users");
  Users.SetDataLoader(dbUsersDataLoader);
  Users.SetAboutToDeleteItemCallback(dbUsersDelete)
}

func dbUsersDelete(item *cache2go.CacheItem) {
  user := item.Data().(*User)
  cancelMode(user)
}

func dbUsersDataLoader(key interface{}, args ...interface{}) *cache2go.CacheItem {
  err := args[1].(*error);

  //A: Load them from the database
  var user User;
  *err = DB.QueryRow("SELECT * FROM users WHERE id = ?", key.(int64)).Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Permissions)
  if *err != nil {
    return nil;
  }

  return cache2go.NewCacheItem(user.ID, time.Duration(Config.CacheLifespan) * time.Second, &user);
}

func dbGetUser(requestUser *tgbotapi.User) (*User, error) {
  var err error;
  item, _ := Users.Value(requestUser.ID, requestUser, &err); //NOTE: The error returned is useless
  if err == nil { //A: They are cache'd
    return item.Data().(*User), nil; 
  }

  if err != sql.ErrNoRows { //A: Error other than not being in the db
    return nil, err
  }

  //A: They're a new user
  user, err := dbNewUser(requestUser);
  if err != nil {
    return nil, err;
  }
  return user, nil;
}

func dbSaveUser(user *User) error {
  Users.Add(user.ID, time.Duration(Config.CacheLifespan) * time.Second, user);
  
  _, err := DB.Exec("INSERT OR REPLACE INTO users (id, firstname, lastname, username, permissions) VALUES (?,?,?,?,?);", user.ID, user.FirstName, user.LastName, user.UserName, user.Permissions);
  return err;
}

func dbNewUser(requestUser *tgbotapi.User) (*User, error) {
  user := &User{
    ID: requestUser.ID,
    FirstName: requestUser.FirstName,
    LastName: requestUser.LastName,
    UserName: requestUser.UserName,
    Permissions: 2, //A: Default to user
  };

  err := dbSaveUser(user);
  if err != nil {
    return nil, err;
  }

  log.Printf("[dbNewUser] New user: %+v", user);
  return user, nil;
}
