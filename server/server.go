package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"io"
	"encoding/json"

	pb "github.com/joeymhills/rpi-facial-detection/proto"
	vision "google.golang.org/genproto/googleapis/cloud/vision/v1p4beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
)

// Implement the ImageServiceServer interface
type imageServer struct{
  pb.UnimplementedImageServiceServer
}

// Handles the image data when it's uploaded from Raspberry Pi
func handleImage() func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    // Reads image from request body
    imageData, err := io.ReadAll(r.Body)
    if err != nil {
      log.Fatalln("Error in reading request body:", err)
      http.Error(w, "Error reading request body", http.StatusBadRequest)
      return
    }

    // Sends image to be annotated
    resp, err := detectFaces(&imageData)
    if err != nil {
      log.Fatalln("Error in detectFaces:", err)
      http.Error(w, "Error processing image", http.StatusInternalServerError)
      return
    }

    log.Println("Response from GCP:", resp)
    // Return response to client
    json.NewEncoder(w).Encode(resp)
  }
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
  resp, err := visionClient.BatchAnnotateImages(ctx, &vision.BatchAnnotateImagesRequest{
    Requests: []*vision.AnnotateImageRequest{&req}})
  if err != nil{
    log.Fatalln("err in image annotation request: ", err)
  }
  return resp, nil
}

func StartServer() {
  //Creates a TCP listener
  lis, err := net.Listen("tcp", ":80")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  //Initializes HTTP image server
  srv := http.NewServeMux()
  srv.HandleFunc("POST /detect/", handleImage())
  http.Serve(lis, srv)
}
