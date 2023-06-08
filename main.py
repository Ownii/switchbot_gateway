from switchbotmeter import DevScanner
import requests
import configparser

config = configparser.ConfigParser()
config.read("config.ini")

current_devices = next(DevScanner())
for device in current_devices:
    data = {
        'mac': device.mac,
        'model': device.model,
        'mode': device.mode,
        'date': str(device.date),
        'temp': device.temp,
        'humidity': device.humidity,
    }
    requests.post(config["API"]["DATA_ENDPOINT"], json=data)
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
