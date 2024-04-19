package server

import (
  "bytes"
  "log"
  "net/http"
  "os"
  "fmt"

  "cloud.google.com/go/storage"
  "github.com/joho/godotenv"
)

//Uploader struct for gcp bucket
type ClientUploader struct {
  cl *storage.Client
  projectID string
  bucketName string
  fileName string
}

func (c *ClientUploader) StoreImage(imageData *[]byte) error {
  err := godotenv.Load()
  if err != nil {
      log.Fatal("Error loading .env file")
  }
  url := os.Getenv("GCP_BUCKET_ADDR")
  url = fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/%s/o?uploadType=media&name=%s", c.bucketName, c.fileName)
  //timeout context
  //ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
  //defer cancel()

  client := http.Client{}
  
  r := bytes.NewReader(*imageData)
  req, err := http.NewRequest(http.MethodPost, url, r)
  if err != nil {
    return err
  }
  // Set the Content-Type header
  req.Header.Set("Content-Type", "application/octet-stream")

  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  
  if resp.StatusCode != http.StatusOK {
      log.Println("Upload failed. Status code", resp.StatusCode)
  }
  
  return nil
}
