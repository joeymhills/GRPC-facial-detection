package main

import (
  "context"
  "log"
  "net"

  "google.golang.org/grpc"
  pb "github.com/joeymhills/rpi-facial-detection/proto"
  )
  // Implement the ImageServiceServer interface
  type imageServer struct{}

  func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error) {
    // Handle the image upload here (e.g., save the image data)
    log.Println("Received image data")
    return &pb.ImageResponse{Message: "Image received successfully"}, nil
  }

  func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
      log.Fatalf("failed to listen: %v", err)
    }
    srv := grpc.NewServer()
    pb.RegisterImageServiceServer(srv, &imageServer{})
    log.Println("Server started")
    if err := srv.Serve(lis); err != nil {
      log.Fatalf("failed to serve: %v", err)
    }
  }

