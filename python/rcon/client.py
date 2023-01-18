import socket
import random
import rcon

# from rcon.error import RconError


class Rcon:

    host = None
    port = None
    password = None

    def __init__(self, host: str, port: int, password: str):
        self.host = host
        self.port = port
        self.password = password

        for attr in self.__dict__.values():
            if attr is None:
                raise ValueError("Missing argument")

        # create socket
        self.connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        # set time out
        self.connection.settimeout(5)

    def connect(self) -> None:

        try:
            # connect to server
            self.connection.connect((self.host, self.port))

            # send auth information
            self.connection.send(
                self.__make_packet(rcon.SERVERDATA_AUTH, self.password)
            )

        # 需要区分不同种类的错误并进行对应的处理
        except socket.error as e:
            raise e

        res: bytes = self.connection.recv(4096)
        res_dict = self.__unpack_packet(res)

        if (
            res_dict["packet_type"] != rcon.SERVERDATA_AUTH_RESPONSE
            or res_dict["id"] == -1
        ):
            raise ValueError("Authentication failed")

    def command(self, command: str) -> str:

        # send command
        try:
            self.connection.send(
                self.__make_packet(rcon.SERVERDATA_EXECCOMMAND, command)
            )

            # receive response
            res: bytes = self.connection.recv(4096)
        except socket.error as e:
            raise e

        res_dict = self.__unpack_packet(res)

        if res_dict["packet_type"] != rcon.SERVERDATA_RESPONSE_VALUE:
            raise ValueError("Invalid response")

        return res_dict["body"]

    def close(self) -> None:
        self.connection.close()

    def __make_packet(self, packet_type: int, body: str) -> bytearray:

        id: int = random.randint(0, 0x7FFFFFFF)

        # make packet
        packet = bytearray()
        packet += id.to_bytes(4, byteorder="little")
        packet += packet_type.to_bytes(4, byteorder="little")
        packet += bytearray(body, "utf-8")
        packet += b"\x00"

        length = len(packet)
        packet = length.to_bytes(4, byteorder="little") + packet

        return packet

    def __unpack_packet(self, packet: bytearray) -> dict:

        length = int.from_bytes(packet[0:4], byteorder="little")
        id = int.from_bytes(packet[4:8], byteorder="little")
        packet_type = int.from_bytes(packet[8:12], byteorder="little")
        body = packet[12:-1].decode("utf-8")

        return {"length": length, "id": id, "packet_type": packet_type, "body": body}

    def __repr__(self):
        return f"<Rcon {self.host}:{self.port}>"

    def __enter__(self):
        print("enter by __enter__")
        self.connect()
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        print("exit by __exit__")
        self.close()
