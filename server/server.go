package server

import (
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/joeymhills/rpi-facial-detection/proto"
	"google.golang.org/grpc"
)

// Implement the ImageServiceServer interface
type imageServer struct{
  pb.UnimplementedImageServiceServer
}

//Implementationn of grpc UploadImage function
func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error) {
  //Path to saved image
  imagePath := "img/newfile.jpg"
  
  //saves uploaded image.
  err := os.WriteFile(imagePath, req.ImageData, 0777)  
  if err != nil {
    log.Fatalln("error in saving jpeg:", err)
  }

  return &pb.ImageResponse{Message: "Image received successfully"}, nil
}

func detectFaces(imageData []byte){

}

//TODO: Remove waitgroup
func StartServer(wg *sync.WaitGroup) {
  //Creates a tcp listener on port 50051
  lis, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  //Initializes grpc image server
  srv := grpc.NewServer()
  pb.RegisterImageServiceServer(srv, &imageServer{})
  log.Println("Server started")

  wg.Done()

  if err := srv.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

