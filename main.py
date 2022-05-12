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
		fin.write(str(_id))

def handleUpdate(update):
	if not update.message.text is None:
		if update.message.text == '/test':
			api.sendMessage(update.message.chat.id, 'Test', update.message.id)
		if update.message.text == '/ip' and update.message.user.id in cfg.chatIds:
			api.sendMessage(update.message.chat.id, requests.get('https://api.ipify.org').text, update.message.id)

if __name__ == '__main__':
	updates = api.getUpdates(lastUpdateId() + 1)
	if len(updates) > 0:
		print(updates)
		for update in updates:
			handleUpdate(update)
		saveLastUpdate(updates[-1].id)
	else: pass #TODO: Log no new updates
