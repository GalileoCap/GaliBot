import api

from utils import *
import handlers.simple as simple
import handlers.todo as todo

def handle(message) -> None:
  command : str = message.command;
  if command == 'todo': todo.handle(message); #A: TODO list
  else: simple.handle(message); #A: Default

if __name__ == '__main__':
  updates : list[api.Update] = api.getUpdates(lastUpdateId())
  if len(updates) > 0:
    print(updates)
    for update in updates:
      if update.message:
        handle(update.message)
      #TODO: Else other type of updates
    saveLastUpdate(updates[-1].id);
  else: print('No updates'); 
