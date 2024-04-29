//go:build server
// +build server

package main

import (
    "fmt"
    "log"
    "os"
    "bytes"
    "image/jpeg"

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
    
    file, err := os.Open("img/joey1.jpg") // Replace with the actual path to your image
    if err != nil {
        fmt.Println("Error opening image file:", err)
        return
    }
    defer file.Close()

    // Read the image data into a byte slice
    imageData, err := jpeg.Decode(file)
    if err != nil {
        fmt.Println("Error reading image file:", err)
        return
    }
    
    buf := new(bytes.Buffer)

    // Encode the image into JPEG format and write it to the buffer
    err = jpeg.Encode(buf, imageData, nil)
    if err != nil {
        fmt.Println("Error encoding image:", err)
        return
    }
    imageBytes := buf.Bytes()

    err = server.HandleImage(db, &imageBytes)
    if err != nil {
        fmt.Println("error when starting opencv", err)
    }
    
    if err = server.InitGrpcServer(db); err != nil {
        log.Fatalln(err)
    }

}
