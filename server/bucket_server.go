package server

import (
  "bytes"
  "log"
  "net/http"
  "os"
  "fmt"
)
//Receives imageData and a fileName and stores it in the gcp bucket
func StoreImage(imageData *[]byte, fileName string) error {
  
  client := http.Client{}
  url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/%s/o?uploadType=media&name=%s", os.Getenv("GCP_BUCKET_NAME"), fileName)
  
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
