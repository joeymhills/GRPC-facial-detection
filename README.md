# Raspberry Pi Facial Detection Project

This project involves setting up a Raspberry Pi with a motion sensor and camera to capture images when motion is sensed. The captured images are then sent to a GoLang server running on a Google Cloud Platform (GCP) Compute Engine instance via gRPC. The server utilizes the Google AI Vision API to detect facial landmarks and determine whether the person in the image is a stranger or not.

## Raspberry Pi Setup

### 1. Install Go and C Compiler

Ensure that Go programming language and a C compiler are installed on your Raspberry Pi.

### 2. Install wiringPi

Install wiringPi library on your Raspberry Pi. You can find the wiringPi GitHub repository [here](https://github.com/WiringPi/WiringPi).

### 3. Set Environment Variables

Create a `.env` file in your project directory and set the following environment variables:

```plaintext
GCP_ADDR=127.0.0.1
GCP_PORT=:8080
```

### 4. Build and Run the Program

Use the following commands to build and run the program on your Raspberry Pi:

```plaintext
go build -tags rpi
./rpi-facial-detection
```
The go build -tags rpi command builds the program specifically for Raspberry Pi, and ./rpi-facial-detection runs the compiled program.

