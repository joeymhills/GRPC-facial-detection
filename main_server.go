//go:build server
// +build server

package main

import (
    "fmt"
    "log"
    "github.com/charmbracelet/huh"
    "github.com/joeymhills/rpi-facial-detection/server"
    "github.com/joho/godotenv"
    "errors"
)
var (
    action string
    toppings []string
    sauceLevel int
    name string
    instructions string
    discount bool
)

func main(){
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    initialForm := huh.NewForm(
        huh.NewGroup(
            // Ask the user for a base burger and toppings.
            huh.NewSelect[string]().
                Title("Choose an action").
                Options(
                huh.NewOption("Start server", "server"),
                huh.NewOption("Add a new face scan", "train"),
                ).
                Value(&action), // store the chosen option in the "burger" variable
            ),

        )
    trainForm := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Who's face are we scanning?").
                Value(&name).
                // Validating fields is easy. The form will mark erroneous fields
                // and display error messages accordingly.
                Validate(func(str string) error {
                    if str == "missy" {
                        return errors.New("Sorry, that name has already been chosen.")
                    }
                    return nil
                }),
            ),
        )

    err = initialForm.Run()
    if err != nil {
        log.Fatal(err)
    }

    if action == "train" {
        err = trainForm.Run()
        if err != nil {
            log.Fatal(err)
        }

        server.TrainModelFromMp4("vid/justin.mp4", name) 
    }

    db, err := server.InitDb()
    if err != nil{
        log.Fatalln(err)
    }
    //err = server.ProcessTestImages() 
    if err != nil {
        fmt.Println(err) 
    }
    server.TestModelFromImage(db, "joey.jpg")

    if err = server.InitGrpcServer(db); err != nil {
        log.Fatalln(err)
    }

}
