package server

import (
  "context"
  "database/sql"
  "fmt"
  "os"
  "log"

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
  
  log.Println("Sql server connected")
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
