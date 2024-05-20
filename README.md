# Raspberry Pi Facial Detection Project
## Overview

This project involves setting up a Raspberry Pi with a motion sensor and camera to capture images when motion is sensed. The captured images are then sent to a GoLang server running on a Google Cloud Platform (GCP) Compute Engine instance via gRPC. The server utilizes the OpenCV and TensorFlow to predict whether or not the face is an authorized person. A more detailed look into the project can be found in the slides [here](https://owlssouthernct-my.sharepoint.com/:p:/g/personal/hillsj3_southernct_edu/Ea9vDnWmG7JMpVtWc0n_JlkBxaOFLcqpfA_gY5P61IPeRQ?e=PifgCm).

## GCP Server Setup

### 1. Create a GCP Compute Engine Virtual Machine

Create a virtual machine hosted on the Google Cloud running Debian linux

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

Create a `.env` file in the root of your project directory and set the following environment variables(Replace with your own credentials):

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

## Example Outputs

![e7ab3732f5d687b9ab0d6b33e3987c861539385660af7865643346585f7e36a5](https://github.com/dsc333/dsc333-final-project-submissions-joeymhills/assets/69769618/efbe6fc7-2915-4672-99e8-054bb9cfdd53)

![d4bf1f60d1d9fbe5b520284a6220926b2a0cda22fd25fe633fd1e77e3e08434a](https://github.com/dsc333/dsc333-final-project-submissions-joeymhills/assets/69769618/ce8b8ca4-623a-405b-9815-392e9a84380d)

./rpi-facial-detection
```
The go build -tags rpi command builds the program specifically for Raspberry Pi, and ./rpi-facial-detection runs the compiled program.
