package main;

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "encoding/json"
  "strings"
  "strconv"
  "fmt"
  "log"
  "errors"
)

type TODOEntry struct {
  ID int64;
  Title string `json:",omitempty"`;
  Description string `json:",omitempty"`;
  Urgency string `json:",omitempty"`;
  Length string `json:",omitempty"`;
  Tags []string `json:",omitempty"`;

  MessageID int64;
  Step int;
};

func cmdTodo(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = "This is your TODO hub"; //TODO: Better description
  reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("New Entry", "todoNew;"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("List All", "todoList;"),
      tgbotapi.NewInlineKeyboardButtonData("Random Entry", "todoRandom;"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("Edit Entry", "todoEdit;"),
      tgbotapi.NewInlineKeyboardButtonData("Configure", "todoConfig;"),
    ),
  );
  return nil;
}

func todoGetEntries(user *User, page int) ([]TODOEntry, error) {
  entries, present := TODODb[user.ID];
  if !present {
    return []TODOEntry{}, nil;
  }
  return entries[min(page * 6, len(entries)) : min((page + 1) * 6, len(entries))], nil; //TODO: Simplify
}

func todoGetEntry(user *User, id int64) TODOEntry {
  entries, _ := todoGetEntries(user, 0); //TODO
  return entries[0];
}

func todoAddEntry(user *User, entry *TODOEntry) error {
  entries, present := TODODb[user.ID];
  if !present {
    TODODb[user.ID] = []TODOEntry{ *entry };
  }
  TODODb[user.ID] = append(entries, *entry);
  //TODO: Entry ID

  return nil;
}

//************************************************************
//S: todoNew

func todoNewGetEntry(user *User) TODOEntry {
  var entry TODOEntry;
  value, present := UserCache[user.ID];
  if !present {
    value = []byte("{}");
  }
  json.Unmarshal(value, &entry);
  return entry;
}

func todoNewSetEntry(user *User, entry *TODOEntry) {
  bytes, _ := json.Marshal(entry);
  UserCache[user.ID] = bytes;
}

func todoNewEdit(user *User, entry *TODOEntry) tgbotapi.EditMessageTextConfig {
  return tgbotapi.NewEditMessageTextAndMarkup(user.ID, int(entry.MessageID), "", tgbotapi.NewInlineKeyboardMarkup());
}

func cbTodoNew(user *User, data string, msg *tgbotapi.Message, edit *tgbotapi.EditMessageTextConfig) error {
  //TODO: Reply with title -> description ? -> Tags ? -> Urgency -> Length -> Edit/Send, Back in every step
  entry := todoNewGetEntry(user);

  switch parts := strings.Split(data, ";"); parts[0] {
  case "": //A: Start
    entry = TODOEntry{MessageID: int64(msg.MessageID)}; //A: Reset

  case "edit":
    entry.Step = 0;
  case "back":
    entry.Step--;
  case "next":
    entry.Step++;
  case "urgency":
    entry.Urgency = parts[1];
    entry.Step++;
  case "length":
    entry.Length = parts[1];
    entry.Step++;
  case "confirm":
    if err := todoAddEntry(user, &entry); err != nil {
      return err;
    }
    entry.Step++;
  }
  todoNewSetEntry(user, &entry); //A: Cache the entry

  todoNewFoo(user, &entry, edit);
  return nil;
}

func todoNewFoo(user *User, entry *TODOEntry, edit *tgbotapi.EditMessageTextConfig) {
  var action string;
  backButton := tgbotapi.NewInlineKeyboardButtonData("Back", "todoNew;back");
  nextButton := tgbotapi.NewInlineKeyboardButtonData("Next", "todoNew;next");

  switch entry.Step {
  case -1: //A: Back to hub
    edit.ReplyMarkup = nil;
    //TODO

  case 0: //A: Title
    action = "Please reply to this message with the entry's title.";
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        backButton,
      ),
    );
    //TODO: Next if already set

    enterMode(user, "todoNewTitle");
    //TODO: Force reply

  case 1: //A: Description
    action = "Please reply to this message with the entry's description, or press the button to leave it empty.";
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        backButton,
        nextButton,
      ),
    );

    enterMode(user, "todoNewDescription");
    //TODO: Force reply

  case 2: //A: Tags
    action = "Please reply to this message with the entry's tags separated by commas, or press the button to leave them empty.";
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        backButton,
        nextButton,
      ),
    );

    enterMode(user, "todoNewTags");
    //TODO: Force reply

  case 3: //A: Urgency
    action = "Please select the entry's urgency.";
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Urgent", "todoNew;urgency;urgent"),
        tgbotapi.NewInlineKeyboardButtonData("Medium", "todoNew;urgency;medium"),
        tgbotapi.NewInlineKeyboardButtonData("Low", "todoNew;urgency;low"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        backButton,
      ),
    );

    cancelMode(user);

  case 4: //A: Length
    action = "Please select the entry's expected length.";
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Long", "todoNew;length;long"),
        tgbotapi.NewInlineKeyboardButtonData("Medium", "todoNew;length;medium"),
        tgbotapi.NewInlineKeyboardButtonData("Short", "todoNew;length;short"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        backButton,
      ),
    );

  case 5: //A: Edit/Send
    edit.Text = fmt.Sprintf("Creating new TODO entry (final step!).\nPlease confirm or edit the information:\n\t- Title: \"%s\"\n\t- Description: \"%s\"\n\t- Tags: %v\n\t- Urgency: %s\n\t- Length: %s", entry.Title, entry.Description, entry.Tags, entry.Urgency, entry.Length);
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Edit", "todoNew;edit"),
        tgbotapi.NewInlineKeyboardButtonData("Confirm", "todoNew;confirm"),
      ),
    );

    return;

  case 6: //A: Confirmed
    edit.Text = fmt.Sprintf("Successfully created new TODO entry!\nInformation:\n\t- Title: \"%s\"\n\t- Description: \"%s\"\n\t- Tags: %v\n\t- Urgency: %s\n\t- Length: %s\nYou can use the unique ID %v to edit/delete it later.", entry.Title, entry.Description, entry.Tags, entry.Urgency, entry.Length, entry.ID);
    edit.ReplyMarkup = nil;

    return;
  }

  edit.Text = fmt.Sprintf("Creating new TODO entry (step %v/5).\n\n%v\n\nYou'll be able to edit it at the end, also you can ignore this message at any time to cancel.", entry.Step+1, action);
  //TODO: Warn about illegal characters

}

func modeTodoNewTitle(user *User, msg *tgbotapi.Message) {
  //TODO: Check length and characters

  entry := todoNewGetEntry(user);
  entry.Title = msg.Text;
  entry.Step++;
  todoNewSetEntry(user, &entry); //A: Cache the entry

  edit := todoNewEdit(user, &entry);
  todoNewFoo(user, &entry, &edit);

  if err := sendEdit(edit); err != nil {
    log.Printf("[modeTodoNewTitle/sendEdit] Error: %v, %+v", err, edit);
  }
}

func modeTodoNewDescription(user *User, msg *tgbotapi.Message) {
  //TODO: Check length and characters

  entry := todoNewGetEntry(user);
  entry.Description = msg.Text;
  entry.Step++;
  todoNewSetEntry(user, &entry); //A: Cache the entry

  edit := todoNewEdit(user, &entry);
  todoNewFoo(user, &entry, &edit);

  if err := sendEdit(edit); err != nil {
    log.Printf("[modeTodoNewDescription/sendEdit] Error: %v, %+v", err, edit);
  }
}

func modeTodoNewTags(user *User, msg *tgbotapi.Message) {
  //TODO: Check length and characters

  entry := todoNewGetEntry(user);
  entry.Tags = strings.Split(msg.Text, ",");
  entry.Step++;
  todoNewSetEntry(user, &entry); //A: Cache the entry

  edit := todoNewEdit(user, &entry);
  todoNewFoo(user, &entry, &edit);

  if err := sendEdit(edit); err != nil {
    log.Printf("[modeTodoNewTags/sendEdit] Error: %v, %+v", err, edit);
  }
}

//************************************************************
//S: todoList

type TODOListCache struct {
  Page int;
};

func todoListGetCache(user *User) TODOListCache { //TODO: Generic read from cache
  var cache TODOListCache;
  value, present := UserCache[user.ID];
  if !present {
    value = []byte("{}");
  }
  json.Unmarshal(value, &cache);
  return cache;
}

func todoListEntry(entry TODOEntry) tgbotapi.InlineKeyboardButton {
  return tgbotapi.NewInlineKeyboardButtonData(
    fmt.Sprintf("%v (%v, %v)", entry.Title, entry.Urgency, entry.Length),
    fmt.Sprintf("todoEntry;%v", entry.ID),
  );
}

func todoListRow(entries []TODOEntry, i int) []tgbotapi.InlineKeyboardButton {
  if 2*i+1 < len(entries) {
    return tgbotapi.NewInlineKeyboardRow(
      todoListEntry(entries[2*i]), todoListEntry(entries[2*i+1]),
    );
  } else {
    return tgbotapi.NewInlineKeyboardRow(
      todoListEntry(entries[2*i]),
    );
  }
}

func cbTodoList(user *User, data string, msg *tgbotapi.Message, edit *tgbotapi.EditMessageTextConfig) error {
  cache := todoListGetCache(user);

  parts := strings.Split(data, ";");
  switch parts[0] {
  case "": //A: First time
    cache.Page = 0;
  case "pageBack":
    cache.Page = max(cache.Page - 1, 0);
  case "pageNext":
    cache.Page++; //TODO: Limit when you've reached the end
  }
  
  entries, _ := todoGetEntries(user, cache.Page);

  edit.Text = fmt.Sprintf("Page %v/X", cache.Page+1); //TODO: Total number of pages
  //TODO: Description and query

  markup := make([][]tgbotapi.InlineKeyboardButton, 0, 5); 
  for i := 0; i < (len(entries) / 2); i++ {
    markup = append(markup, todoListRow(entries, i));
  }
  if len(entries) % 2 != 0 { //A: There's a row of only one element
    markup = append(markup, todoListRow(entries, len(entries) / 2));
  }
  markup = append(markup, tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("Previous Page", "todoList;pageBack"), //TODO: Not on page 0
    tgbotapi.NewInlineKeyboardButtonData("Edit Query", "todoList;query"),
    tgbotapi.NewInlineKeyboardButtonData("Next Page", "todoList;pageNext"),
  ));
  markup = append(markup, tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("Back", "todoHub;"),
  ));

  *edit.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: markup};

  //TODO: Paging

  return nil;
}

//************************************************************
//S: todoEntry

func cbTodoEntry(user *User, data string, msg *tgbotapi.Message, edit *tgbotapi.EditMessageTextConfig) error {
  parts := strings.Split(data, ";");
  if len(parts) != 1 {
    return errors.New(fmt.Sprintf("[cbTodoEntry] Invalid data: %v", data));
  }
  id, err := strconv.ParseInt(parts[0], 10, 64);
  if err != nil {
    return err;
  }
  entry := todoGetEntry(user, id);

  edit.Text = fmt.Sprintf(
    "- Title: %v\n- Description: %v\n- Tags: %v\n- Urgency: %v\n- Length: %v\n- ID: %v\n",
    entry.Title, entry.Description, entry.Tags, entry.Urgency, entry.Length, entry.ID,
  );
  *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("Delete", "todoEntry;delete"),
      tgbotapi.NewInlineKeyboardButtonData("Edit", "todoEntry;edit"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("Back", "todoBack;"),
    ),
  );

  return nil;
}
