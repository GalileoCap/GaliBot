import pandas as pd
import re
import db

dfName = 'todo';

def add(message):
  print('todo.handle.add');
  df = db.getDf(dfName);
  _, tags, text, _ = re.split(r'^add\s*(\[[^\]]*\])\s*"([^("$)]*)"$', message.text);
  data = {
    'userId': message.user.id,
    'tags': tags,
    'text': text,
  };
  ndf = pd.concat([df, pd.DataFrame([data])], ignore_index = True).drop_duplicates(); #TODO: Don't reset index
  db.saveDf(ndf, dfName);
  message.respond('added');

def select(message):
  print('todo.handle.select');
  df = db.getDf(dfName);
  message.respond(str(df)); #TODO: Fancier #TODO: Divide in pages

def handle(message):
  _, subcommand, _ = re.split(r'^(\w+)', message.text); subcommand = subcommand.lower();
  if subcommand == 'add': add(message);
  elif subcommand == 'select': select(message);
  else: pass #TODO: Unknown
