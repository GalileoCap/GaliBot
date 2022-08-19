import os

import cfg

def forcePath(path):
  os.makedirs(path, exist_ok = True);
  return path;

def lastUpdateId():
  try:
    with open(f'{forcePath(cfg.logDir)}/{cfg.lastUpdateFile}', 'r') as fin:
      return int(fin.readline());
  except FileNotFoundError:
    #TODO: Log
    return 0; #A: Default to 0

def saveLastUpdate(_id):
  with open(f'{forcePath(cfg.logDir)}/{cfg.lastUpdateFile}', 'w') as fin:
    fin.write(str(_id));
