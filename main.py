from switchbotmeter import DevScanner
import requests

for current_devices in DevScanner():
    for device in current_devices:
        requests.post('http://192.168.178.89:8091/switchbot', device)
        print("===========")
        print(f'mac: {device.mac}')
        print(f'model: {device.model}')
        print(f'mode: {device.mode}')
        print(f'date: {device.date}')
        print(f'temp: {device.temp}')
        print(f'humidity: {device.humidity}')
        print("===========")
        print()
        print()
