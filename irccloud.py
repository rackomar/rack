import requests
from os import environ
import heroku3
from time import sleep
import asyncio
import ssl
import websockets

ssl_context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
ssl_context.load_verify_locations(requests.certs.where()) # use requests certificate store for the websocket connection

async def activate_socket():
    uri = "wss://api.irccloud.com/"
    async with websockets.connect(
        uri, ssl=ssl_context
    ) as websocket:

        await websocket.recv()
            
auth_formtoken = requests.post("https://www.irccloud.com/chat/auth-formtoken").json()["token"]

session_id = requests.post("https://www.irccloud.com/chat/login", data={"email": environ.get("IRCCLOUD_USERNAME"), "password": environ.get("IRCCLOUD_PASSWORD"), "token": auth_formtoken})

asyncio.get_event_loop().run_until_complete(activate_socket())

heroku3.from_key(environ['heroku-key']).apps()[environ['heroku-app-name']].scale_formation_process('worker', 0)
sleep(1)
