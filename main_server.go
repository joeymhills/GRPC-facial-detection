//go:build server
// +build server

package main

import (
	"fmt"
	"log"

	"github.com/joeymhills/rpi-facial-detection/server"
	"github.com/joho/godotenv"
)

func main(){
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    err = server.InitOpencv()
    if err != nil {
        fmt.Println("error when starting opencv", err)
    }

    db, err := server.InitDb()
    if err != nil{
        log.Fatalln(err)
    }
    
    if err = server.InitGrpcServer(db); err != nil {
        log.Fatalln(err)
    }

}
