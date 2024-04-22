//go:build server
// +build server

package main

import (
	"log"

	"github.com/joeymhills/rpi-facial-detection/server"
	"github.com/joho/godotenv"
)

func main(){
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    db, err := server.InitDb()
    if err != nil{
        log.Fatalln(err)
    }
    
    if err = server.InitGrpcServer(db); err != nil {
        log.Fatalln(err)
    }

}
