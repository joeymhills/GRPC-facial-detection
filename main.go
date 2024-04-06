package main

import(
    "github.com/joeymhills/rpi-facial-detection/server"
    "github.com/joeymhills/rpi-facial-detection/client"

    "sync"
)

func main() {
    var wg sync.WaitGroup
    
    wg.Add(1)
    go server.StartServer(&wg)

    wg.Wait()

    client.SendImage()

}
