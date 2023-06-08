from switchbotmeter import DevScanner

for current_devices in DevScanner():
    for device in current_devices:
        print(device)
        print(f'{device.mac} -> {device.temp}')