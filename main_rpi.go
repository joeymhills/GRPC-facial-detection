//go:build rpi
// +build rpi

package main

import(
   "log"
   "github.com/joeymhills/rpi-facial-detection/client"
   "github.com/joho/godotenv"
)

func main(){
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
   client.WaitForMotion()
}
