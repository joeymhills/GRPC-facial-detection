//go:build server
// +build server 

package main

import(
    "github.com/joeymhills/rpi-facial-detection/server"
    "github.com/joho/godotenv"
    "log"
)

func main(){
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    server.InitGrpcServer()
    _, err = server.InitDb()
}
