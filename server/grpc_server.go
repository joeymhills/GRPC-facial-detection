package server

import (
  "context"
  "log"
  "net"
  "crypto/tls"
  "os"

  pb "github.com/joeymhills/rpi-facial-detection/proto"
  //vision "google.golang.org/genproto/googleapis/cloud/vision/v1p4beta1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
)
// Implement the ImageServiceServer interface
type imageServer struct{
  pb.UnimplementedImageServiceServer
}
//Handles the image data when its uploaded from raspberry pi
func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error){
  log.Println("gRPC endpoint hit!")
  //DetectFaces(&req.ImageData)
  return &pb.ImageResponse{Message: "Image received successfully"}, nil
}
func InitGrpcServer(){

  port := os.Getenv("GCP_PORT")

  //Load the server's certificate and private key
  cert, err := tls.LoadX509KeyPair("server/server.crt", "server/server.key")
  if err != nil {
    log.Fatalf("Failed to load certificate: %v", err)
  }

  creds := credentials.NewTLS(&tls.Config{
    Certificates: []tls.Certificate{cert},
  })
  // Create a new gRPC server with TLS credentials
  grpcServer := grpc.NewServer(grpc.Creds(creds))

  // Register your gRPC service implementation
  pb.RegisterImageServiceServer(grpcServer, &imageServer{})

  // Create a TCP listener on port 8080 with TLS configuration
  listener, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("Failed to create listener: %v", err)
  }

  // Start the gRPC server with TLS-enabled listener                  
  log.Println("gRPC server running on port", port)                     
  if err := grpcServer.Serve(listener); err != nil {                  
    log.Fatalf("Failed to serve: %v", err)                      
  }

}
