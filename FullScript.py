# Unified_Data_Exchange.py
import socket
import time
import ast
from functools import reduce

# --- Configuration ---
PORT = 37020
BROADCAST_ADDR = "255.255.255.255"
TIMEOUT = 1  # seconds

# --- Cipher Engine ---
class FibCipher:
    @staticmethod
    def _fib(n):
        return reduce(lambda x, n: [x[1], x[0]+x[1]], range(n), [0,1])[0]

    @classmethod
    def encrypt(cls, data):
        data = list(data)
        for i in range(len(data)):
            data[i] = chr(ord(data[i]) + cls._fib(i) + 73)
        return "".join(data)

    @classmethod
    def decrypt(cls, data):
        data = list(data)
        for i in range(len(data)):
            data[i] = chr(ord(data[i]) - cls._fib(i) - 73)
        return "".join(data)

# --- Network Operations ---
def send_response(addr, data):
    encrypted_data = FibCipher.encrypt(str(data))
    response = f"RESPONSE:{addr[0]}:{encrypted_data}"
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as reply_s:
        reply_s.sendto(response.encode(), addr)

def receive_data(key):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        s.settimeout(TIMEOUT)

        request = f"REQUEST:{key}"
        s.sendto(request.encode(), (BROADCAST_ADDR, PORT))
        print("[Searching] Broadcasting request...")

        start_time = time.time()
        nearest_data = None

        while True:
            try:
                data, addr = s.recvfrom(1024)
                message = data.decode().strip()
                
                if message.startswith("RESPONSE:"):
                    parts = message.split(":", 2)
                    ip = parts[1]
                    encrypted_data = parts[2]
                    received_data = ast.literal_eval(FibCipher.decrypt(encrypted_data))
                    delay = time.time() - start_time
                    
                    print(f"[Found] Device {ip} responded in {delay:.3f}s")
                    print(f"[Data] {received_data}")
                    
                    nearest_data = {
                        'ip': ip,
                        'data': received_data,
                        'delay': delay
                    }
                    break

            except socket.timeout:
                if nearest_data:
                    print(f"[Timeout] Final data: {nearest_data['data']}")
                    return nearest_data
                print("[Error] No devices found")
                return None

def listen_for_requests(key, shared_data):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        s.bind(("0.0.0.0", PORT))

        print("[Giver] Waiting for requests...")
        while True:
            data, addr = s.recvfrom(1024)
            message = data.decode().strip()
            
            if message.startswith("REQUEST:"):
                request_key = message.split(":")[1]
                if request_key == key:
                    print(f"[Giver] Valid request from {addr[0]}")
                    send_response(addr, shared_data)
                else:
                    print(f"[Giver] Invalid key from {addr[0]}")

# --- Main Program ---
def main():
    print("=== Secure P2P Data Exchange ===")
    print("1. Share data (Giver mode)")
    print("2. Search for data (Searcher mode)")
    choice = input("Select mode (1/2): ")

    key = input("Enter access key: ")
    
    if choice == "1":
        data_value = input("Enter value to share: ")
        data_info = input("Enter description: ")
        shared_data = {
            "value": data_value,
            "info": data_info
        }
        print("\n[System] Starting in Giver mode...")
        listen_for_requests(key, shared_data)
    else:
        print("\n[System] Starting in Searcher mode...")
        result = receive_data(key)
        if result:
            print("\n=== Received Data ===")
            print(f"From: {result['ip']}")
            print(f"Value: {result['data']['value']}")
            print(f"Info: {result['data']['info']}")
            print(f"Response time: {result['delay']:.3f}s")

if __name__ == "__main__":
    main()