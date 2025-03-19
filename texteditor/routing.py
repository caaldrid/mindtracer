from django.urls import path
from . import consumer

websocket_patterns = [
    path(r'ws/text/', consumer.TextEditorConsumer.as_asgi(), name="websocket"),
]