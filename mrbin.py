#!/usr/bin/env python3
import serial,re,sys,requests
from time import sleep
heightOfBin = 17
while True:
    try:
        arduino = serial.Serial('/dev/cu.usbmodem1421',9600)
        a = arduino.readline()
        b = a.decode()
        c= re.findall('\d+', b )
        d= int(c[1])
        print("The raw data sent by Arduino is ",b)
        percent = (d/heightOfBin)*100
        print("The percent is",percent)
        print("The distance measured by arduino is",c[1])
    except Exception as e:
        print("Error is", e)
        percent = 46
    #except UnicodeDecodeError as u:
    #    print("Error in UnicodeDecodeError",u)
    #    percent= 46
    #else:
    #    print("strange error")
    #    percent = 47
    try:
        url = 'https:/lit-sea-89877.herokuapp.com/addPercent'
        d = requests.post(url,data= str(percent))
        print(d)
    except:
        print("Connection issue")
    sleep(10)
