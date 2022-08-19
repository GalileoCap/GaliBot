import requests

import api
import cfg

def test(message):
  print('simple test');
  api.sendMessage(message.chat.id, f'Test: {message.text}', message.id);

def ip(message):
  print('simple ip')
  api.sendMessage(message.chat.id, requests.get('https://api.ipify.org').text, message.id);

def handle(message):
  command = message.command;
  if command == 'test': test(message);
  elif command == 'ip' and message.user.id in cfg.chatIds: ip(message);
  else: print(f'handleMessage unknown command {message.id}'); #TODO: Generic response
