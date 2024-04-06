package client

import (
	"log"
	"os"
	"os/exec"
    "context"

	"google.golang.org/grpc"

	pb "github.com/joeymhills/rpi-facial-detection/proto"
)

type imageClient struct{
  pb.UnimplementedImageServiceServer
}

func SendImage() {
  imagePath := "img/temp.jpg"

  //TODO: Change command to "libcamera-still -o img/temp.jpg" 
  //executes libcamera to capture image 
  cmd := exec.Command("ls", "-l")
  err := cmd.Run()
  if err != nil {
    log.Println("error with cli", err)
    return
  }

  //Reads data from imagePath
  imageData, err := os.ReadFile(imagePath)
  if err != nil {
    log.Println("error reading image data:", err) 
  }

  // Set up a connection to the server
  conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithInsecure())
  if err != nil {
    log.Fatalln("Failed to dial server:", err)
  }
  defer conn.Close()

  //Creates client gRPC client
  client := pb.NewImageServiceClient(conn)
  
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


