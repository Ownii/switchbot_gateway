from switchbotmeter import DevScanner

for current_devices in DevScanner():
    for device in current_devices:
        print("===========")
        print(f'mac: {device.mac}')
        print(f'modal: {device.modal}')
        print(f'mode: {device.mode}')
        print(f'date: {device.date}')
        print(f'temp: {device.temp}')
        print(f'humidity: {device.humidity}')
        print("===========")
        print()
        print()
