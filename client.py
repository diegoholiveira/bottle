# -*- encoding: utf-8 -*-
"""
    Bottle Client

    This files contains an example of a Bottle client
"""

import socket

class Bottle(object):
    """Connects into the bottle server
    """

    def __init__(self, host="localhost", port=42000):
        self.server = (host, port,)
        self.sock = None

    def _connect(self):
        if self.sock is None:
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.sock.connect(self.server)

    def _send(self, message):
        self._connect()

        self.sock.send(message + "\n")
        response = self.sock.recv(4096)
        return response.strip()

    def use(self, queue):
        msg = "USE {}".format(queue)
        response = self._send(msg)
        if response == "OK":
            return True
        raise Exception(response)

    def put(self, data):
        msg = "PUT {}".format(data)
        response = self._send(msg)
        if response == "OK":
            return True
        raise Exception(response)

    def get(self):
        response = self._send("GET")
        if response == "NULL":
            return None
        return response
    
    def close(self):
        if self.sock is not None:
            self.sock.close()
            self.sock = None

conn = Bottle()

conn.use("emails")
for v in range(1, 99999):
    conn.put("mensagem de teste {}".format(v))

while True:
     data = conn.get()
     if data is None:
         break
     print "Data from server: {}".format(data)
 
conn.close()

