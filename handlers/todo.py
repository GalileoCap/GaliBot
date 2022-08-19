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
  ndf = pd.concat([df, pd.DataFrame([data])]);
  db.saveDf(ndf, dfName);
  message.respond('added');

def handle(message):
  _, subcommand, _ = re.split(r'^(\w+)', message.text);
  if subcommand == 'add': add(message);
  else: pass #TODO: Unknown
