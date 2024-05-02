# Raspberry Pi Facial Detection Project

This project involves setting up a Raspberry Pi with a motion sensor and camera to capture images when motion is sensed. The captured images are then sent to a GoLang server running on a Google Cloud Platform (GCP) Compute Engine instance via gRPC. The server utilizes the Google AI Vision API to detect facial landmarks and determine whether the person in the image is a stranger or not.

## GCP Server Setup

### 1. Create a GCP Compute Engine Virtual Machine

In order to make use of the GCP AI Vision API you need to be on a virtual machine hosted on the Google Cloud running Debian linux

### 2. Install Dependencies

Ensure that Go programming language (>v1.22) is installed. Go installation instructions [here](https://go.dev/doc/install)

Also install [opencv](https://docs.opencv.org/4.x/d7/d9f/tutorial_linux_install.html) along with [gocv](https://gocv.io/getting-started/linux/)

Create a python virtual environment

pip install tensorflow opencv-python numpy



### 3. Clone the Repository

Clone the GitHub repository for this project using the following command:

```bash
git clone https://github.com/joeymhills/GRPC-facial-detection.git
```
### 4. Set Environment Variables

Create a `.env` file in the root of your project directory and set the following environment variables(Replace with addess and port to your GCP virtual machine):

```plaintext
GCP_ADDR="34.66.85.133"
GCP_PORT=":8080"
GCP_BUCKET_NAME="dsc333-hw2"
SQL_ADDR=35.193.123.84
SQL_USER=joey
SQL_PASS=2654
SQL_NAME=db
```

### 5. Build and Run the Program

Use the following commands to build and run the program on your Raspberry Pi:

```plaintext
go build -tags server
./rpi-facial-detection
```
The go build -tags server command builds the program specifically for the virtual machine, and ./rpi-facial-detection runs the compiled program.


## Raspberry Pi Setup

### 1. Install Go and C Compiler

Ensure that Go programming language and a C compiler are installed on your Raspberry Pi.

### 2. Install wiringPi

Install wiringPi library on your Raspberry Pi. You can find the wiringPi GitHub repository [here](https://github.com/WiringPi/WiringPi).

### 3. Clone the Repository

Clone the GitHub repository for this project using the following command:

```bash
git clone https://github.com/joeymhills/GRPC-facial-detection.git
```

### 4. Set Environment Variables

Create a `.env` file in the root of your project directory and set the following environment variables(Replace with addess and port to your GCP virtual machine):

```plaintext
GCP_ADDR=127.0.0.1
GCP_PORT=:8080
```

### 5. Build and Run the Program

Use the following commands to build and run the program on your Raspberry Pi:

```plaintext
go build -tags rpi
./rpi-facial-detection
```
The go build -tags rpi command builds the program specifically for Raspberry Pi, and ./rpi-facial-detection runs the compiled program.
