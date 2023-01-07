package main

import (
  "fmt"
)

var Types = []string{ "Breakfast", "Lunch", "Merienda", "Dinner", "Extra" } //TODO: Const

func fmtEntry(entry *Entry) string {
  date := entry.Date.Format("Mon 02/01/2006, 03:04")
  if entry.Skipped {
    return fmt.Sprintf("- Date: %v\n.\t- Type: %v.\n\t- Skipped.", date, Types[entry.Type])
  } else {
    tags := []string{}
    if entry.Meat {
      tags = append(tags, "Meat")
    }
    if entry.Veggies {
      tags = append(tags, "Veggies")
    }
    if entry.Fruit {
      tags = append(tags, "Fruit")
    }
    return fmt.Sprintf("- Date: %v\n\t- Type: %v.\n\t- Tags:%v\n\t- Description: %v.", date, entry.Type, tags, entry.Description)
  }
}
