package server

import (
    "bytes"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "io"
    "net"
    "os"
    "os/exec"
    "time"
    "path/filepath"
    "strings"

    _ "github.com/go-sql-driver/mysql"

    "gocv.io/x/gocv"
)
func imageToBytes(img gocv.Mat) (*[]byte, error) {
    //Convert mat object to []byte
    newImg, err := img.ToImage()
    if err != nil {
	return  nil, err
    }
    var buf bytes.Buffer
    _ = jpeg.Encode(&buf, newImg, nil)
    returnImg := buf.Bytes()

    return &returnImg, err
}

//Calls python script that trains tensorflow model from the images in "trainimg"
func TrainModel(name string) error {
    cmd := exec.Command("python3", "python/train.py", name)
    output, err := cmd.CombinedOutput()
    if err != nil {
	fmt.Println(string(output))
	return err    
    }
    fmt.Println(string(output))

    return nil
}

//CheckFace is a function that will take in an img(as type *[]byte)
//and return whether or not it was succesfully recognized by
//any of the CNN's
func CheckFace(imgData *[]byte) ([]string, error) {

    pythonHost := "127.0.0.1"
    pythonPort := "49522"
    
    var conn net.Conn

    for attempt := 1; attempt <= 5; attempt++ {
	fmt.Printf("Attempting to connect (attempt %d/%d)...\n", attempt, 5)
	
	var err error
	
	// Dials python server
	conn, err = net.Dial("tcp", pythonHost+":"+pythonPort)
	if err != nil {
	    // Wait before the next attempt
	    fmt.Printf("Connection failed: %v, retrying in 2 seconds\n", err)
	    time.Sleep(time.Second * 2)	
	    continue
	}
    
	fmt.Println("Connection successful!")
	
	// Writes to TCP socket
	_, err = conn.Write(*imgData)
	if err != nil {
	    fmt.Println("Error sending image data:", err)
	    return nil, err
	}
	defer conn.Close()

	// Receive the string response from Python
	response, err := io.ReadAll(conn)
	if err != nil {
	    fmt.Println("Error receiving model names response:", err)
	    return nil, err
	}

	if len(response) == 0 {
	    fmt.Println("Received empty response from TCP stream")
	    return nil, nil
	}

	// Convert the comma-separated string to a slice of strings
	models := strings.Split(string(response), ",")
	return models, nil
    }

    return nil, fmt.Errorf("failed to connect after maximum attempts")
}

// spawnPythonProcess spins a up a python script on a new goroutine
// and opens a tcp socket that we will connect to with our go program
// and send images to to be processed with our tensorflow models
func spawnPythonProcess(){

    // Get a list of files from the TensorFlow models directory
    files, err := os.ReadDir("python/savedModels")
    if err != nil {
	fmt.Println("Error reading directory:", err)
	return
    }

    // Build a comma-separated list of filenames
    var models []string
    for _, file := range files {
	if !file.IsDir() {  // Ensure we're only dealing with files
	    models = append(models, file.Name())
	}
    }

    // Join all model filenames with a comma
    modelsList := strings.Join(models, ",")

    // Print the models list
    fmt.Println("Models to process:", modelsList)
    //get list of models and give to the python command
    cmd := exec.Command("python3", "python/test.py", modelsList)
    _, err = cmd.StdoutPipe()
    if err != nil {
	fmt.Println("Error creating stdout pipe:", err)
    }

    cmd.Start()
}

func ScanImage(imgData *[]byte) (faceNum int, err error)  {
    haarCascade := gocv.NewCascadeClassifier()
    haarCascade.Load("aimodels/haarfrontalface.xml")
    defer haarCascade.Close()

    img, err := gocv.IMDecode(*imgData, gocv.IMReadColor)
    if err != nil {
	return 0, err
    }
    if img.Empty() {
	return 0, fmt.Errorf("Error loading image")
    }
    //Starts python tcp listener to listen for image data to run on CNNs
    //go spawnPythonProcess()
    
    //Detects region of picture with face in it
    imgRects := haarCascade.DetectMultiScaleWithParams(img, 1.1, 3, 0, image.Point{200, 200}, image.Point{1500,1500})
    faceNum = 0

    //Iterates through faces and adds them to the returned array
    for _, rect := range imgRects {
	//Increases number of faces
	faceNum++

	//Iterates through all of the CNNs and marks image succesful scans are
	//Converts image into bytes
	face, err := imageToBytes(img.Region(rect))
	if err != nil {
	    return 0, err
	}

	// This is the code that actually checks for facial recognition
	names, err := CheckFace(face)
	if err != nil{
	    fmt.Println(err)
	}

	if len(names) > 0 {
	    fmt.Printf("Face recognized, welcome %s!\n", names[0])
	    //Draws rectangle around face
	    gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

	    // Calculate text size and position
	    size := gocv.GetTextSize(names[0], gocv.FontHersheyPlain, 3.0, 5)
	    pt := image.Pt(rect.Min.X+(rect.Max.X-rect.Min.X)/2-(size.X/2), rect.Min.Y-size.Y-10)

	    // Put text on the image
	    gocv.PutText(&img, names[0], pt, gocv.FontHersheyPlain, 3.0, color.RGBA{0, 255, 0, 1}, 5)

	} else {
	    fmt.Printf("Face not recognized, possible intruder.\n")

	    //Draws rectangle around face
	    gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)

	    // Draws rectangle around face
	    gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)

	    // Calculate text size and position
	    size := gocv.GetTextSize("Not Recognized", gocv.FontHersheyPlain, 3.0, 5)
	    pt := image.Pt(rect.Min.X+(rect.Max.X-rect.Min.X)/2-(size.X/2), rect.Min.Y-size.Y-10)
	    
	    // Put text on the image
	    gocv.PutText(&img, "Not Recognized", pt, gocv.FontHersheyPlain, 3.0, color.RGBA{0, 255, 0, 1}, 5)
	}
    }
    gocv.IMWrite("output_image.jpg", img)
    if err != nil {
	return 0, err
    }

    return faceNum, nil
}

func emptyDir(dirPath string) error {
    files, err := os.ReadDir(dirPath)
    if err != nil {
	return err
    }

    for _, file := range files {
	err := os.RemoveAll(filepath.Join(dirPath, file.Name()))
	if err != nil {
	    return err
	}
    }
    return nil
}

func TrainModelFromMp4(filepath string, name string) {
    
    // Empties training image directory
    emptyDir("python/trainimg/class_1")

    haarCascade := gocv.NewCascadeClassifier()
    haarCascade.Load("aimodels/haarfrontalface.xml")
    defer haarCascade.Close()
    
    webcam, err := gocv.VideoCaptureFile(filepath)
    if err != nil {
	//handle err
    }
    defer webcam.Close()

    // Loop to continuously read frames from the video file

    i := 0
    for {
	// Read a frame from the video file
	frame := gocv.NewMat()
	if ok := webcam.Read(&frame); !ok {
	    break
	}
	if frame.Empty() {
	    continue
	}
	//Detects region of picture with face in it
	imgRects := haarCascade.DetectMultiScaleWithParams(frame, 1.1, 3, 0, image.Point{200, 200}, image.Point{1500,1500})
	n := 0;
	//Iterates through faces and adds them to the returned array
	for _, rect := range imgRects {
	     n++
	    //gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)
	    face := frame.Region(rect)
	    //returnFaces = append(returnFaces, face)

	    //Save the image with landmarks
	    outputName := fmt.Sprintf("python/trainimg/class_1/%d.jpg", i); i++
	    gocv.IMWrite(outputName, face)
	}

	// Release resources (important to avoid memory leaks)
	frame.Close()
    }
    err = TrainModel(name)
    if err != nil{
	fmt.Println(err)
    }
}
