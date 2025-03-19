import time
import json
from channels.generic.websocket import WebsocketConsumer

class TextEditorConsumer(WebsocketConsumer):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
    
    def connect(self):
        # accept connection
        self.accept()
    
    def disconnect(self, code):
        return super().disconnect(code)
    
    def receive(self, text_data=None, bytes_data=None):

        text_data_json = json.loads(text_data)
        content = """
            <div hx-swap-oob="beforeend:#content">
            <p>{time}: {message}</p>
            </div>
        """
        self.send(text_data=content.format(time=time.time(), message=text_data_json['chat_message']))
        


