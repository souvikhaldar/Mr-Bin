#!/usr/bin/env python3
import serial,re,sys
from time import sleep
arduino = serial.Serial('/dev/cu.usbmodem1421',9600)
a = arduino.readline()
b = a.decode()
c= re.findall('\d+', b )
print(c[1])
