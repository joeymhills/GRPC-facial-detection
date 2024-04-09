package client

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib -lwiringPi
#include <wiringPi.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

#define COMMAND "libcamera-still -o img/temp.jpg"

int motionSensor() {
  // Initialize WiringPi library
  if (wiringPiSetup() == -1) {
    fprintf(stderr, "Failed to initialize WiringPi\n");
    return 1;
  }
  //Sets pin 22 to input mode
  pinMode(22, INPUT);

  //waits for motion to be detected
  while(digitalRead(22) != 0) {
    sleep(.1);
  }
  //Once the sensor is triggered then we execute the bash command to take a picture
  system(COMMAND);

  return 0;
}

*/
import "C"
import (
  "log"
  "os"
  "time"
  "net/http"
  "mime/multipart"
  "bytes"
  "io"
)

//TODO: probably not idiomatic
func sendImage() {

  //address for google vm
  addr := "https://34.68.52.223:80"
  imagePath := "img/temp.jpg"

  // Open the image file
  file, err := os.Open(imagePath)
  if err != nil {
      log.Println("Error opening image file:", err)
      return
  }
	defer file.Close()
  // Create a new HTTP request
  req, err := http.NewRequest("POST", addr, nil)
  if err != nil {
    log.Println("Error creating request:", err)
    return
  }

  // Create a new form data body
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)

  // Add the image file to the form data body
  part, err := writer.CreateFormFile("image", "image.jpg")
  if err != nil {
    log.Println("Error creating form file:", err)
    return
  }
  _, err = io.Copy(part, file)
  if err != nil {
    log.Println("Error copying file to form:", err)
    return
  }

  // Close the form data writer
  err = writer.Close()
  if err != nil {
    log.Println("Error closing form writer:", err)
    return
  }

  // Set the content type header
  req.Header.Set("Content-Type", writer.FormDataContentType())

  // Set the request body
  req.Body = io.NopCloser(body)

  // Send the request
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Println("Error sending request:", err)
    return
  }
  defer resp.Body.Close()

  // Print the response status code
  log.Println("Response status code:", resp.StatusCode)
}

func WaitForMotion() {
  //Calls C code that waits for motion
  if C.motionSensor() == 0 {
    //Once motion is sensed we take a picture
    log.Println("motion sensed")
    sendImage()

    //Sleep to prevent taking too many pictures
    time.Sleep(time.Second*2)
    
    //Recursively call WaitForMotion() to reinstate motion detection mode
    WaitForMotion()
  }
}
