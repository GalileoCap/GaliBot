def lastUpdateId():
  try:
    with open('./lastUpdate.log', 'r') as fin:
      return int(fin.readline());
  except FileNotFoundError:
    #TODO: Log
    return 0; #A: Default to 0

def saveLastUpdate(_id):
  with open('./lastUpdate.log', 'w') as fin:
    fin.write(str(_id));
