#include <wiringPi.h>
#include <stdio.h>

int main() {
    // Initialize WiringPi library
    if (wiringPiSetup() == -1) {
        fprintf(stderr, "Failed to initialize WiringPi\n");
        return 1;
    }

    // Set pin mode to input
    pinMode(22, INPUT);

    while (1) {
        //Reads for motion sensor input
        if(digitalRead(22) == 0){
            printf("Motion Detected!\n");
        }
            // Add some delay to avoid reading too frequently
            delay(10);
        }

        return 0;
    }
}
