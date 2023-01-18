# -*- coding: utf-8 -*-
import socket
import random
import time

host = "192.168.0.30"
port = 25575
address = (host, port)
password = "logitech"

SERVERDATA_AUTH = 3
SERVERDATA_AUTH_RESPONSE = 2
SERVERDATA_EXECCOMMAND = 2
SERVERDATA_RESPONSE_VALUE = 0


def make_packet(packet_type: int, body: str) -> bytearray:

    id: int = random.randint(0, 0x7FFFFFFF)

    # make packet
    packet = bytearray()
    packet += id.to_bytes(4, byteorder="little")
    packet += packet_type.to_bytes(4, byteorder="little")
    packet += bytearray(body, "utf-8")
    packet += b"\x00"

    length = len(packet)
    packet = length.to_bytes(4, byteorder="little") + packet

    print(id, packet_type, body, length)
    print(packet)

    return packet


if __name__ == "__main__":
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(address)

    # send auth information
    s.send(make_packet(SERVERDATA_AUTH, password))

    msg = s.recv(2048)

    s.send(make_packet(SERVERDATA_EXECCOMMAND, "say hello world"))

    print(msg)

    time.sleep(10000)

    s.close()
