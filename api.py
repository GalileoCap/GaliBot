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
		self.user = User(obj['from']) if 'from' in obj else None
		self.chat = Chat(obj['chat']) if 'chat' in obj else None
		self.text = obj.get('text')
		#TODO: Everything else as needed

	def __repr__(self):
		return f'{{id: {self.id}, user: {self.user}, text: {self.text}}}'

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
