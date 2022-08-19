import requests
import cfg

class Update:
  def __init__(self, obj):
    self.id = obj.get('update_id')
    self.message = Message(obj['message']) if 'message' in obj else None
    #TODO: Everything else as needed
  
  def __repr__(self):
    return f'{{id: {self.id}, message: {self.message}}}'

class Message:
  def __init__(self, obj):
    self.id = obj.get('message_id')
    self.user = User(obj['from']); #A: Messages have to have a user
    self.chat = Chat(obj['chat']); #A: Messages have to have a chat

    raw = obj.get('text');
    parts = raw.split() if raw else [];

    isCommand = len(parts) > 0 and parts[0][0] == '/';
    self.command = parts[0][1:] if isCommand else None;
    self.subcommand = parts[1] if isCommand and len(parts) > 1 else None;
    self.text = ' '.join(parts[1:]) if isCommand else raw;
    #TODO: Everything else as needed

  def __repr__(self):
    return f'{{id: {self.id}, user: {self.user}, command: {self.command}, text: {self.text}}}'

  def respond(self, text):
    sendMessage(self.chat.id, text, self.id);

class User:
  def __init__(self, obj):
    self.id = obj.get('id')
    self.username = obj.get('username')
    #TODO: Everything else as needed

  def __repr__(self):
    return f'{{id: {self.id}, username: {self.username}}}'

class Chat:
  def __init__(self, obj):
    self.id = obj.get('id')
  
  def __repr__(self):
    return f'{{id: {self.id}}}'

def getUpdates(lastUpdateId = -1): #TODO: Get better params
  r = requests.get(
    f'https://api.telegram.org/bot{cfg.API_TOKEN}/getUpdates',
    params = {'offset': lastUpdateId + 1}, #A: Only new updates 
  ).json()
  if r['ok'] == 'False': print('getUpdates ERROR') #TODO: Write log and throw error
  return [Update(obj) for obj in r['result']]

def sendMessage(chatId, text, reply_to): #TODO: Get better params
  r = requests.post(
    f'https://api.telegram.org/bot{cfg.API_TOKEN}/sendMessage',
    params = {'chat_id': chatId, 'text': text, 'reply_to_message_id': reply_to},
  )
  if r.json()['ok'] == 'False': print('sendMessage ERROR') #TODO: Write long and throw error

#TODO: Better handling of missing parameters in classes
