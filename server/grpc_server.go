package server

import (
  "context"
  "crypto/tls"
  "database/sql"
  "log"
  "net"
  "os"
  "fmt"

  pb "github.com/joeymhills/rpi-facial-detection/proto"
  //vision "google.golang.org/genproto/googleapis/cloud/vision/v1p4beta1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
)

// Implement the ImageServiceServer interface
type imageServer struct{
  pb.UnimplementedImageServiceServer
  db *sql.DB
}

//gRPC endpoint handler
func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error){
  log.Println("gRPC endpoint hit!")

  HandleImage(s.db, &req.ImageData)

  return &pb.ImageResponse{Message: "Image received successfully"}, nil
}

//Handles the image data when its uploaded from raspberry pi
func HandleImage(db *sql.DB, imgBytes *[]byte) error {
  //ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
  //defer cancel()

  
  //Check for facial features 
  numFaces, faces, err := GetFaceImages(imgBytes)
  if err != nil {
    return  err
  }

  //StoreImage(img, "test123.jpg")
  log.Println("Faces detected: ", numFaces)
  
  if faces != nil {
    for _, face := range *faces {
      success, err := CheckFace(face, "lebron")
      if err != nil{
        log.Println(err)
      }
      if success {
        fmt.Printf("Face recognized, welcome Lebron!")

      } else {
        fmt.Printf("Face not recognized, possible intruder.")
      }
    }
  }
  return nil
}

func InitGrpcServer(db *sql.DB) error{

  port := os.Getenv("GCP_PORT")

  //Load the server's certificate and private key
  cert, err := tls.LoadX509KeyPair("server/cert/server.crt", "server/cert/server.key")
  if err != nil {
    return err
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
    return err
  }

  // Start the gRPC server with TLS-enabled listener                  
  log.Println("gRPC server running on port", port)                     
  if err := grpcServer.Serve(listener); err != nil {                  
    return err
  }
return err
}
