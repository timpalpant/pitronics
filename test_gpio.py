# External module imports
import RPi.GPIO as GPIO
import time

pins = [2, 3, 4, 17]
GPIO.setmode(GPIO.BCM)
for pin in pins:
    GPIO.setup(pin, GPIO.OUT)

try:
    while True:
        for pin in pins:
            GPIO.output(pin, GPIO.LOW)
            time.sleep(0.1)
        time.sleep(1)
        for pin in pins:
            GPIO.output(pin, GPIO.HIGH)
            time.sleep(0.1)
        time.sleep(1)
except KeyboardInterrupt:
    GPIO.cleanup() # cleanup all GPIO
