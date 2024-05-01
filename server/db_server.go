package server

import (
  "context"
  "database/sql"
  "fmt"
  "os"
  "log"
  "bytes"
  "net/http"

  _ "github.com/go-sql-driver/mysql"
)

type FaceData struct {
  ids []int
  faces []string
}

//InitDb connects us to the sql server and returns the connection object
func InitDb() (*sql.DB, error){
  //Grabs credentials from .env file and puts them into the connection string
  conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",os.Getenv("SQL_USER"), os.Getenv("SQL_PASS"), os.Getenv("SQL_ADDR"), os.Getenv("SQL_NAME"))
  
  //starts db connection
  db, err := sql.Open("mysql", conn)
  if err != nil {
    return nil, err
  }
  
  return db, nil
}

func Get_Faces(ctx context.Context, db *sql.DB) (*FaceData, error) {

  q := "SELECT * FROM faces;"
  rows, err := db.QueryContext(ctx, q)
  if err != nil {
    return nil, err
  }
  var ids []int
  var faces []string

  for rows.Next() {
    var id int
    var face string
    if err := rows.Scan(&id, &face); err != nil{
      return nil, err 
    }
    ids = append(ids, id)
    faces = append(faces, face)
  }
  faceData := &FaceData{ids, faces}

  return faceData, nil
}

//Receives imageData and a fileName and stores it in the gcp bucket
func StoreImage(ctx context.Context, db *sql.DB, imageData *[]byte, fileName string) error {
  
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

  
  /* Insert image into sql table
  q := "INSERT INTO faces VALUES id, ;"
  if _, err := db.QueryContext(ctx, q); err != nil {
    return err
  }
  */
  return nil
}
