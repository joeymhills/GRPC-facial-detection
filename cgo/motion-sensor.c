#include <wiringPi.h>
#include <stdio.h>

#define SENSOR_PIN 2

int main() {
    // Initialize WiringPi library
    if (wiringPiSetup() == -1) {
        fprintf(stderr, "Failed to initialize WiringPi\n");
        return 1;
    }

    // Set pin mode to input
    pinMode(SENSOR_PIN, INPUT);

    while (1) {
        // Read the state of the sensor pin
        int sensorState = digitalRead(SENSOR_PIN);

        // Print the state
        printf("Motion sensor state: %d\n", sensorState);

        // Add some delay to avoid reading too frequently
        delay(1000);  // 1 second delay
    }

    return 0;
}
