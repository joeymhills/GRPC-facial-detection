package server

import (
  "context"
  "database/sql"
  "fmt"
  "log"
  "os"

  _ "github.com/go-sql-driver/mysql"
)

//InitDb connects us to the sql server and returns the connection object
func InitDb() (*sql.DB, error){
  //grabs credentials from .env file and puts them into the connection string
  conn := fmt.Sprintf("root:%s@tcp(%s)/%s", os.Getenv("SQL_PASS"), os.Getenv("SQL_ADDR"), os.Getenv("DATABASE_NAME"))
  
  log.Println("attempting to connect to db...")
  db, err := sql.Open("mysql", conn)
  if err != nil {
    return nil, err
  }
  log.Println("Database connection successful")
  
  return db, nil
}

func Send_query(ctx context.Context, db *sql.DB) (*string, error){
  select {
  case<-ctx.Done():
    return nil, ctx.Err()
  default:
    q := "CREATE DATABASE faces"
    _, err := db.QueryContext(ctx, q)
    if err != nil {
      return nil, err
    }
  }
  return nil, nil 
}
