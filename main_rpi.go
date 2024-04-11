//go:build rpi
// +build rpi

package main

import(
    "github.com/joeymhills/rpi-facial-detection/client"
)

func main(){
   client.WaitForMotion()
}
