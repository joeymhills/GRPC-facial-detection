package server

import (
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

func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error) {
   
  imagePath := "img/newfile.jpg"

  err := os.WriteFile(imagePath, req.ImageData, 0777)  
  if err != nil {
    log.Fatalln("error in saving jpeg:", err)
  }

  log.Println("Received image data")
  return &pb.ImageResponse{Message: "Image received successfully"}, nil
}

func StartServer(wg *sync.WaitGroup) {
  lis, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  srv := grpc.NewServer()
  pb.RegisterImageServiceServer(srv, &imageServer{})
  log.Println("Server started")

  wg.Done()

  if err := srv.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

