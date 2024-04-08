#include <wiringPi.h>
#include <stdio.h>

int64_t motionSensor() {
    // Initialize WiringPi library
    if (wiringPiSetup() == -1) {
        fprintf(stderr, "Failed to initialize WiringPi\n");
        return 1;
    }

    // Set pin mode to input
    pinMode(22, INPUT);
    int i = 0;
    //Waits for the motion sensor to read 0 and then terminates the function
    while (digitalRead(22) != 0) {
        delay(10);
    }
    return 0;
}
