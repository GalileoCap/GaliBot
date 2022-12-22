package main;

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "encoding/json"
  "strings"
  "fmt"
  "log"
  //"errors"
)

type TODONewEntry struct {
  Title string `json:",omitempty"`;
  Description string `json:",omitempty"`;
  Urgency string `json:",omitempty"`;
  Length string `json:",omitempty"`;
  Tags []string `json:",omitempty"`;

  Step int;
};

func cmdTodo(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = "This is your TODO hub"; //TODO: Better description
  reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("New Entry", "todoNew;"),
    ),
    //tgbotapi.NewInlineKeyboardRow(
      //tgbotapi.NewInlineKeyboardButtonData("List All", "todoList"),
      //tgbotapi.NewInlineKeyboardButtonData("Random Entry", "todoRandom"),
    //),
  );
  return nil;
}

func todoGetEntries(user *User) ([]TODONewEntry, error) {
  entries, present := TODODb[user.ID];
  if !present {
    return []TODONewEntry{}, nil;
  }
  return entries, nil;
}

func todoAddEntry(user *User, entry *TODONewEntry) error {
  entries, present := TODODb[user.ID];
  if !present {
    TODODb[user.ID] = []TODONewEntry{ *entry };
  }
  TODODb[user.ID] = append(entries, *entry);

  return nil;
}

func todoNewGetEntry(msg *tgbotapi.Message) TODONewEntry {
  var entry TODONewEntry; //TODO: Get entry from text
  parts := strings.Split(msg.Text, "\n");
  json.Unmarshal([]byte(parts[len(parts) - 1]), &entry);
  return entry;
}

func cbTodoNew(user *User, data string, msg *tgbotapi.Message, edit *tgbotapi.EditMessageTextConfig) error {
  //TODO: Reply with title -> description ? -> Tags ? -> Urgency -> Length -> Edit/Send, Back in every step
  entry := todoNewGetEntry(msg);

  switch parts := strings.Split(data, ";"); parts[0] {
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
    err := todoAddEntry(user, &entry);
    //TODO: Change message
    return err;
  }

  foo(user, &entry, edit);
  return nil;
}

func foo(user *User, entry *TODONewEntry, edit *tgbotapi.EditMessageTextConfig) {
  var action string;
  backButton := tgbotapi.NewInlineKeyboardButtonData("Back", "todoNew;back");
  nextButton := tgbotapi.NewInlineKeyboardButtonData("Next", "todoNew;next");

  switch entry.Step {
  case -1: //A: Back to hub
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
    entryB, _ := json.Marshal(entry);
    edit.Text = fmt.Sprintf("Creating new TODO entry (final step!).\nPlease confirm or edit the information:\n\t- Title: \"%s\"\n\t- Description: \"%s\"\n\t- Tags: %v\n\t- Urgency: %s\n\t- Length: %s\n%s", entry.Title, entry.Description, entry.Tags, entry.Urgency, entry.Length, string(entryB));
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Edit", "todoNew;edit"),
        tgbotapi.NewInlineKeyboardButtonData("Confirm", "todoNew;confirm"),
      ),
    );
    return;
  }

  entryB, _ := json.Marshal(entry);
  edit.Text = fmt.Sprintf("Creating new TODO entry (step %v/5).\n\n%v\n\nYou'll be able to edit it at the end, also you can ignore this message at any time to cancel.\n%s", entry.Step, action, string(entryB));
  //TODO: Warn about illegal characters
}

func modeTodoNewTitle(user *User, msg *tgbotapi.Message) {
  if msg.ReplyToMessage == nil {
    //TODO: Didn't reply
    return;
  }
  //TODO: Check length and characters

  entry := todoNewGetEntry(msg.ReplyToMessage);
  log.Printf("aaa %+v", entry);
  entry.Title = msg.Text;
  entry.Step++;

  edit := newEdit(user, msg.ReplyToMessage);
  foo(user, &entry, &edit);

  if err := sendEdit(edit); err != nil {
    log.Printf("[modeTodoNewTitle/sendEdit] Error: %v, %+v", err, edit);
  }
}

func modeTodoNewDescription(user *User, msg *tgbotapi.Message) {
  if msg.ReplyToMessage == nil {
    //TODO: Didn't reply
    return;
  }
  //TODO: Check length and characters

  entry := todoNewGetEntry(msg.ReplyToMessage);
  entry.Description = msg.Text;
  entry.Step++;

  edit := newEdit(user, msg.ReplyToMessage);
  foo(user, &entry, &edit);

  if err := sendEdit(edit); err != nil {
    log.Printf("[modeTodoNewDescription/sendEdit] Error: %v, %+v", err, edit);
  }
}

func modeTodoNewTags(user *User, msg *tgbotapi.Message) {
  if msg.ReplyToMessage == nil {
    //TODO: Didn't reply
    return;
  }
  //TODO: Check length and characters

  entry := todoNewGetEntry(msg.ReplyToMessage);
  entry.Tags = strings.Split(msg.Text, ",");
  entry.Step++;

  edit := newEdit(user, msg.ReplyToMessage);
  foo(user, &entry, &edit);

  if err := sendEdit(edit); err != nil {
    log.Printf("[modeTodoNewTags/sendEdit] Error: %v, %+v", err, edit);
  }
}
