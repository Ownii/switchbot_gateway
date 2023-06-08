from switchbotmeter import DevScanner

for current_devices in DevScanner():
    print("test")
    print(current_devices[0].mac)
#    for device in current_devices:
        #print(device)
        #print("found device")
        #print(device.mac)
        #print(device.temp)
        #print(f'{device.mac} -> {device.temp}')