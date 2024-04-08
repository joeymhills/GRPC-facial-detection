package main

import(
    "github.com/joeymhills/rpi-facial-detection/server"
    "github.com/joeymhills/rpi-facial-detection/client"

    "flag"
)

func main() {

    mode := flag.String("m", "default", "enter mode")

    flag.Parse()

    switch *mode{
    
    case "default":
        server.StartServer()
        client.InitializeClient()
    case "client":
        client.InitializeClient()
    case "server":
        server.StartServer()
    }
}
