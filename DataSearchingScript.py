import socket
import time
import ast  # Для преобразования строки с данными в словарь

def find_nearest_device(key):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        s.settimeout(1)

        request = f"REQUEST:{key}"
        s.sendto(request.encode(), ("255.255.255.255", 37020))
        print("Поиск устройств...")

        start_time = time.time()
        nearest_data = None  # Будет хранить данные, а не только IP

        while True:
            try:
                data, addr = s.recvfrom(1024)
                message = data.decode().strip()
                
                if message.startswith("RESPONSE:"):
                    parts = message.split(":", 2)  # Разделяем только на 3 части
                    ip = parts[1]
                    received_data = ast.literal_eval(parts[2])  # Безопасное преобразование строки в dict
                    delay = time.time() - start_time
                    
                    print(f"Устройство {ip} ответило за {delay:.3f} сек")
                    print(f"Данные: {received_data}")  # Выводим содержимое!
                    
                    if not nearest_data or delay < min_delay:
                        nearest_data = {
                            "ip": ip,
                            "data": received_data,
                            "delay": delay
                        }
                        break  # Выход после первого ответа

            except socket.timeout:
                if nearest_data:
                    print(f"Превышен таймаут. Получены данные: {nearest_data['data']}")
                    return nearest_data
                else:
                    print("Устройства не найдены.")
                    return None

if __name__ == "__main__":
    key = input("Введите ключ для поиска: ")
    result = find_nearest_device(key)
    if result:
        print(f"IP: {result['ip']}, Данные: {result['data']}")
