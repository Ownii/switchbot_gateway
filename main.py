from switchbotmeter import DevScanner

for current_devices in DevScanner():
    for device in current_devices:
        #print(device)
        #print("found device")
        print(f'{device.mac}: {device.temp}')
        #print(device.temp)
        #print(f'{device.mac} -> {device.temp}')