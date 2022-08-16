import api
import requests
import cfg

def lastUpdateId():
  try:
    with open('./lastUpdate.log', 'r') as fin:
      return int(fin.readline())
  except FileNotFoundError:
    #TODO: Log
    return 0 #A: Default to 0

def saveLastUpdate(_id):
  with open('./lastUpdate.log', 'w') as fin:
    fin.write(str(_id - 1))

def test(message):
  print('handle.test');
  api.sendMessage(message.chat.id, f'Test: {message.text}', message.id);

def ip(message):
  print('handle.ip')
  api.sendMessage(message.chat.id, requests.get('https://api.ipify.org').text, message.id);

def handleMessage(message):
  command = message.command;
  if command == 'test': test(message);
  elif command == 'ip' and message.user.id in cfg.chatIds: ip(message);

if __name__ == '__main__':
  updates = api.getUpdates(lastUpdateId() + 1)
  if len(updates) > 0:
    print(updates)
    for update in updates:
      handleMessage(update.message)
    saveLastUpdate(updates[-1].id)
  else: print('No updates'); 
