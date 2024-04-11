//go:build server
// +build server 

package main

import(
    "github.com/joeymhills/rpi-facial-detection/server"
)

func main() {
    server.StartServer()
}
