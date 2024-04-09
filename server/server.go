package server

import (
  "context"
  "log"
  "net"
  "crypto/tls"

  pb "github.com/joeymhills/rpi-facial-detection/proto"
  vision "google.golang.org/genproto/googleapis/cloud/vision/v1p4beta1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/credentials/alts"
)

// Implement the ImageServiceServer interface
type imageServer struct{
  pb.UnimplementedImageServiceServer
}
//Handles the image data when its uploaded from raspberry pi
func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error){
  log.Println("gRPC endpoint hit!")
  /*
resp, err := detectFaces(&req.ImageData)
if err != nil {
log.Println("Error in detectFaces:", err)
return &pb.ImageResponse{Message: "Error in detectFaces()"}, nil
}
log.Println("Response from GCP:", &resp)
*/

  return &pb.ImageResponse{Message: "Image received successfully"}, nil
}
//Initializes and returns a GCP Vision client
func SetupVisionClient() (vision.ImageAnnotatorClient, *grpc.ClientConn, error) {
  addr := "vision.googleapis.com:443"
  conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(alts.NewClientCreds(alts.DefaultClientOptions())))
  if err != nil {
    return nil, nil, err
  }
  visionClient := vision.NewImageAnnotatorClient(conn)
  return visionClient, conn, nil
}

//Takes in an image and detects face landmarks
func detectFaces(imageData *[]byte) (*vision.BatchAnnotateImagesResponse, error){
  ctx := context.Background()

  //Gets our vision client and gRPC connection
  visionClient, conn, err := SetupVisionClient()
  if err != nil {
    log.Fatalln("error in creating connection: ", err)
  }
  defer conn.Close()
  //Populates struct with information about the 
  req := vision.AnnotateImageRequest{
    Image: &vision.Image{
      Content: *imageData,
    },
    //Fills the feature type enum(currently set to landmark detection)
    Features: []*vision.Feature{
      &vision.Feature{Type: 2},
    },
    ImageContext: &vision.ImageContext{
      // Optionally, include additional context
    },
  }
  //Sends image annotation request to GCP services
  resp, err := visionClient.BatchAnnotateImages(ctx, &vision.BatchAnnotateImagesRequest{Requests: []*vision.AnnotateImageRequest{&req}})
  if err != nil{
    log.Fatalln("err in image annotation request: ", err)
  }
  return resp, nil
}

func StartServer() {


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
  listener, err := net.Listen("tcp", ":8080")
  if err != nil {
    log.Fatalf("Failed to create listener: %v", err)
  }

  // Start the gRPC server with TLS-enabled listener                  
  log.Println("gRPC server running on port 8080")                     
  if err := grpcServer.Serve(listener); err != nil {                  
    log.Fatalf("Failed to serve: %v", err)                      
  }                                                                   
}    
