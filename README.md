# Bottle

Bottle is a simple message queue. This code itsn't ready for production and will never be!

Bottle protocol is very simple, here is a list of the commands accepted by the server:

| Command | Description |
| --------|------------ |
| USE [QUEUE] | Connects to the QUEUE. It must be the fist command on a new connection |
| PUT "{\"message\":\"hello\"}" | Adds a message to the QUEUE |
| GET | Get the next message of the QUEUE |
| PURGE | Clean the QUEUE |
| QUIT | Close the current connection |


Into this repository, there is a Python scripts that show how to connect and interact with Bottle.

