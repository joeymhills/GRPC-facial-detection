package server

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"io"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/joeymhills/rpi-facial-detection/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func generateUniqueName() string {
	// Get the current date and time as a string
	currentTime := time.Now().UTC().Format(time.RFC3339Nano)

	// Create a SHA-256 hasher
	hasher := sha256.New()

	// Write the current time to the hasher; since Write never returns an error, we ignore its error return
	hasher.Write([]byte(currentTime))

	// Compute the SHA-256 hash and return it as a hexadecimal string
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}

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

// Handles the image data when its uploaded from raspberry pi
func HandleImage(db *sql.DB, imgBytes *[]byte) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
  defer cancel()
  
  // Check for facial features 
  numFaces, err := ScanImage(imgBytes)
  if err != nil {
    return err 
  }
  log.Println("Faces detected: ", numFaces)
  
  // Loads image
  img, err := os.Open("output_image.jpg")
  if err != nil {
    log.Println("error opening image: ", err)
  }

  // Reads image into []byte
  imgData, err := io.ReadAll(img)
  if err != nil {
    log.Println("error reading image: ", err)
  }
  imageName := generateUniqueName()

  
  err = StoreImage(ctx, db, &imgData, imageName+".jpg")
  if err != nil {
    log.Println("err storing image in GCP Bucket: ", err)
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
