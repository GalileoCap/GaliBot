import pandas as pd
import re
import db

dfName = 'todo';

def add(message):
  print('todo.handle.add');
  df = db.getDf(dfName);
  tags = re.search(r'^add\s*(\[[\w,]+\])', message.text);
  text = re.search(r'^add\s*(?:\[[\w,]+])?\s*"([^("$)]*)"$', message.text);
  data = {
    'userId': message.user.id,
    'tags': tags.group(1) if tags else '[]',
    'text': text.group(1) if text else '',
  };
  df = pd.concat([df, pd.DataFrame([data])], ignore_index = True).drop_duplicates(); #TODO: Don't reset index
  db.saveDf(df, dfName);
  message.respond('added'); #TODO: React

def select(message):
  print('todo.handle.select');
  df = db.getDf(dfName);
  message.respond(str(df)); #TODO: Fancier #TODO: Divide in pages

def drop(message):
  print('todo.handle.drop');
  df = db.getDf(dfName);
  _, ids, _ = re.split(r'^drop\s*(.*)', message.text)
  toDrop = [int(x) for x in ids.split(',') if int(x) in df.index];
  badIdx = [int(x) for x in ids.split(',') if not int(x) in df.index];
  df.drop(toDrop, inplace = True);
  db.saveDf(df, dfName);
  message.respond('done'); #TODO: React
  if len(badIdx) > 0:
    message.respond(f'invalid indeces {badIdx}');


def handle(message):
  _, subcommand, _ = re.split(r'^(\w+)', message.text); subcommand = subcommand.lower();
  if subcommand == 'add': add(message);
  elif subcommand == 'select': select(message);
  elif subcommand == 'drop': drop(message);
  else: pass #TODO: Unknown