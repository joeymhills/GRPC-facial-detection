package server

import (
  "context"
  "database/sql"
  "log"
  "time"

  _ "github.com/go-sql-driver/mysql"
)

func Send_query() error{
  
  ctx, cancel := context.WithTimeout(context.Background(), 8 * time.Second)
  defer cancel()
  
  db, err := sql.Open("mysql", "root@tcp(35.193.123.84:3306)/dsc-333-412716:us-central1:rpi")
  if err != nil {
    log.Fatalln("error connecting to db: ", err)
  }
  defer db.Close()
  log.Println("Database connection successful")
 
  select {
  case<-ctx.Done():
    return ctx.Err()
  default:
    q := "CREATE DATABASE faces"
    _, err := db.QueryContext(ctx, q)
    if err != nil {
      return err
    }
  }
  return nil 
}
