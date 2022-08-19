import requests

import api
import cfg

def test(message):
  print('simple.handle.test');
  message.respond(
    f'Test: {message.text}'
  );

def ip(message):
  print('simple.handle.ip');
  response = requests.get('https://api.ipify.org').text;
  message.respond(
    response.removeprefix('http://')
  );

def handle(message):
  command = message.command;
  if command == 'test': test(message);
  elif command == 'ip' and message.user.id in cfg.chatIds: ip(message);
  else: print(f'handleMessage unknown command {message.id}'); #TODO: Generic response
