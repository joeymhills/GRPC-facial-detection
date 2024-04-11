package client

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib -lwiringPi
#include <wiringPi.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

#define COMMAND "libcamera-jpeg -o img/temp.jpg"

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
  "context"
  
  pb "github.com/joeymhills/rpi-facial-detection/proto"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
)

type imageClient struct{
  pb.UnimplementedImageServiceServer
}

func sendImage() {

  //address for google vm
  addr := os.Getenv("GCP_ADDR")+os.Getenv("GCP_PORT")
  imagePath := "img/temp.jpg"

  //Reads data from imagePath
  imageData, err := os.ReadFile(imagePath)
  if err != nil {
    log.Println("error reading image data:", err) 
  }

  creds, err := credentials.NewClientTLSFromFile("server/server.crt", "")
  if err != nil {
	 log.Fatalf("Failed to load server cert:", err)
  }

  // Dial the gRPC server with TLS credentials
  conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
  if err != nil {
	 log.Fatalf("Failed to dial server: %v", err)
  }
  log.Println("gRPC server dialed successfully")

  defer conn.Close()

  //Creates client gRPC client
  client := pb.NewImageServiceClient(conn)

  log.Println("gRPC client ready for requests")

  ctx := context.Background()
  req := &pb.ImageRequest{
    ImageData: imageData,
  }
  resp, err := client.UploadImage(ctx, req)
  if err != nil {
    log.Fatalln("error in sending image to server:", err)
  }

  log.Println("Response from server:", resp)
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
