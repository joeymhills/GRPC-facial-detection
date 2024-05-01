//go:build server
// +build server

package main

import (
    "fmt"
    "log"
    "context"
    "github.com/charmbracelet/huh"
    "github.com/charmbracelet/huh/spinner"
    "github.com/joeymhills/rpi-facial-detection/server"
    "github.com/joho/godotenv"
    "errors"
)
var (
    action string
    toppings []string
    sauceLevel int
    name string
    filepath string
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

            huh.NewInput().
                Title("Please enter file path to face scan scan video").
                Value(&filepath).
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
	ctx, cancel := context.WithCancel(context.Background())
    
    title := fmt.Sprintf("Training a CNN model for %s's face scan...", name)
	s := spinner.New().Title(title).
		Context(ctx)
    
    go s.Run()

	if err != nil {
		log.Fatalln(err)
        cancel()
	}
        server.TrainModelFromMp4(filepath, name) 
        cancel()
        fmt.Printf("Training success!\n")
        return
    }

    db, err := server.InitDb()
    if err != nil{
        log.Fatalln(err)
    }
    if err != nil {
        fmt.Println(err) 
    }
    
    server.TestModelFromImage(db, "justin.jpg")

    if err = server.InitGrpcServer(db); err != nil {
        log.Fatalln(err)
    }

}
