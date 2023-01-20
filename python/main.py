import rcon

host = "192.168.0.30"
port = 25575
password = "logitech"


if __name__ == "__main__":
    try:
        with rcon.Rcon(host, port, password) as client:
            while True:
                command = input("Command: ")
                print(client.command(command))

    except KeyboardInterrupt:
        pass
