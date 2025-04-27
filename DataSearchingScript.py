import socket
import time

def find_nearest_device(key):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        s.settimeout(1)  # Небольшой таймаут для recvfrom

        # Отправляем запрос
        request = f"REQUEST:{key}"
        s.sendto(request.encode(), ("255.255.255.255", 37020))
        print("Поиск устройств...")

        start_time = time.time()
        nearest_ip = None
        min_delay = float('inf')

        while True:
            try:
                data, addr = s.recvfrom(1024)
                message = data.decode().strip()
                
                if message.startswith("RESPONSE:"):
                    parts = message.split(":")
                    ip = parts[1]
                    delay = time.time() - start_time
                    
                    print(f"Устройство {ip} ответило за {delay:.3f} сек")
                    
                    # Если это первый ответ или быстрее предыдущего
                    if delay < min_delay:
                        min_delay = delay
                        nearest_ip = ip
                        print(f"Найдено устройство: {nearest_ip}")
                        break  # Выход после первого успешного ответа

            except socket.timeout:
                if nearest_ip:
                    print(f"Превышен таймаут, ближайшее устройство: {nearest_ip}")
                    return nearest_ip
                else:
                    print("Устройства не найдены.")
                    return None

if __name__ == "__main__":
    key = input("Введите ключ для поиска: ")
    ip = find_nearest_device(key)
    if ip:
        print(f" IP найденного устройства: {ip}")