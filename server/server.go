package server

import (
	"context"
	"log"
	"net"

	pb "github.com/joeymhills/rpi-facial-detection/proto"
	vision "google.golang.org/genproto/googleapis/cloud/vision/v1p4beta1"
	"google.golang.org/grpc"
)

// Implement the ImageServiceServer interface
type imageServer struct{
  pb.UnimplementedImageServiceServer
}

//Handles the image data when its uploaded from raspberry pi
func (s *imageServer) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error) {
  
  resp, err := detectFaces(&req.ImageData)
  if err != nil {
    log.Fatalln("Error in detectFaces:", err)
  }
  log.Println("Response from GCP:", &resp)

  return &pb.ImageResponse{Message: "Image received successfully"}, nil
}

//Initializes and returns a GCP Vision client
func setupVisionClient() (vision.ImageAnnotatorClient, *grpc.ClientConn, error) {
  ctx := context.Background()
  conn, err := grpc.DialContext(ctx, "vision.googleapis.com:443", grpc.WithInsecure())
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
  visionClient, conn, err := setupVisionClient()
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
  resp, err := visionClient.BatchAnnotateImages(ctx, &vision.BatchAnnotateImagesRequest{
    Requests: []*vision.AnnotateImageRequest{&req}})
  if err != nil{
    log.Fatalln("err in image annotation request: ", err)
  }
  return resp, nil
}

//TODO: Remove waitgroup
func StartServer() {
  //Creates a tcp listener on port 50051
  lis, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  //Initializes grpc image server
  srv := grpc.NewServer()
  pb.RegisterImageServiceServer(srv, &imageServer{})
  log.Println("Server started")

  if err := srv.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

