import socket
device_data = {
    "key": "abc123",
    "data": {"value": 42, "info": "Sample data"}
}

def listen_for_requests():
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        s.bind(("0.0.0.0", 37020))

        print("Ожидание запросов...")
        while True:
            data, addr = s.recvfrom(1024)
            message = data.decode().strip()
            
            if message.startswith("REQUEST:"):
                key = message.split(":")[1]
                if key == device_data["key"]:
                    print(f"Запрос от {addr[0]}")
                    response = f"RESPONSE:{addr[0]}:{device_data['data']}"
                    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as reply_s:
                        reply_s.sendto(response.encode(), addr)
                else:
                    print(f"Неверный ключ от {addr[0]}")

if __name__ == "__main__":
    listen_for_requests()