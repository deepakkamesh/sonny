#!/bin/sh
echo "Pinging controller"
./cli ping
#sleep 2

echo "Testing PIR sensor"
./cli pir
#sleep 2

echo "Testing ultrasonic sensor"
./cli dist
#sleep 2

echo "Testing Acceleration"
./cli accel
#sleep 2

echo "Testing Heading"
./cli head
#sleep 2

echo "Testing temperature/humidity"
./cli temp
#sleep 2

echo "Testing battery voltage"
./cli batt
#sleep 2

echo "Testing light level"
./cli ldr

echo "Blinking LED 5 times at 500ms"
./cli blink -d 500 -t 5
sleep 5

echo "Blinking LED 5 times at 1000ms"
./cli blink -t 5




echo "Testing servos"
./cli servo -s 1 -a 70
./cli servo -s 2 -a 70
sleep 2
./cli servo -s 1 -a 90
./cli servo -s 2 -a 90


