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
    db, err := server.InitDb()
    if err != nil{
        log.Fatalln(err)
    }
    err = server.ProcessTestImages() 
    if err != nil {
        fmt.Println(err) 
    }
    server.TestModelFromImage(db, "joey.jpg")
    //server.TrainModelFromMp4("vid/justin.mp4", "justin") 

    if err = server.InitGrpcServer(db); err != nil {
        log.Fatalln(err)
    }

}
