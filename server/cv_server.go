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
func CheckFace(imgData *[]byte) (bool, error) {

    pythonHost := "127.0.0.1"
    pythonPort := "49522"
    
    var conn net.Conn

    for attempt := 1; attempt <= 4; attempt++ {
	fmt.Printf("Attempting to connect (attempt %d/%d)...\n", attempt, 10)
	
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
	i, err := conn.Write(*imgData)
	if err != nil {
	    fmt.Println("Error sending image data:", err)
	    return false, err
	}
	fmt.Println(i)
	defer conn.Close()

	// Receive the boolean response from Python
	boolBuf, err := io.ReadAll(conn)
	if err != nil {
	    fmt.Println("Error receiving boolean response:", err)
	    return false, err
	}

	if len(boolBuf) == 0 {
	    fmt.Println("Received empty response TCP stream")
	    return false, nil
	}
	var resultBool bool
	resultBool = boolBuf[0] != 0

	return resultBool, nil

    }

    return false, nil
}

func spawnPythonScript(){
    //get list of models and give to the python command
    cmd := exec.Command("python3", "python/test.py", "joey.keras,missy.keras")
    _, err := cmd.StdoutPipe()
    if err != nil {
	fmt.Println("Error creating stdout pipe:", err)
    }

    cmd.Start()
}

func ScanAndAnalyzeImage(imgData *[]byte) (faceNum int, err error)  {

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

    // Get a list of files from the tensorflow models directory
    files, err := os.ReadDir("python/savedModels")
    if err != nil {
	fmt.Print(err)
    }

    // Slice to store filenames
    var models []string

    // Iterate over the files and add their names to the slice
    for _, file := range files {
	if !file.IsDir() {
	    models = append(models, file.Name())
	}
    }

    // Print the filenames
    fmt.Println("Files in the directory:")
    for _, model := range models {
	    fmt.Println(model)
    }
    
    //Detects region of picture with face in it
    imgRects := haarCascade.DetectMultiScaleWithParams(img, 1.1, 3, 0, image.Point{200, 200}, image.Point{1500,1500})
    faceNum = 0
    
    //Starts python tcp listener to listen for image data to run on CNNs
    go spawnPythonScript()

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
	success, err := CheckFace(face)
	if err != nil{
	    fmt.Println(err)
	}

	if success {
	    fmt.Printf("Face recognized, welcome!\n")
	    //Draws rectangle around face
	    gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

	    size := gocv.GetTextSize("Recognized", gocv.FontHersheyPlain, 1.2, 2)
	    pt := image.Pt(rect.Min.X+(rect.Min.X/2)-(size.X/2), rect.Min.Y-2)
	    gocv.PutText(&img, "Recognized", pt, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 1}, 2)

	} else {
	    fmt.Printf("Face not recognized, possible intruder.\n")

	    //Draws rectangle around face
	    gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)

	    size := gocv.GetTextSize("Not Recognized", gocv.FontHersheyPlain, 1.2, 2)
	    pt := image.Pt(rect.Min.X+(rect.Min.X/2)-(size.X/2), rect.Min.Y-2)
	    gocv.PutText(&img, "Not Recognized", pt, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 1}, 2)
	}
    }
    gocv.IMWrite("output_image.jpg", img)
    if err != nil {
	return 0, err
    }

    return faceNum, nil
}

func TrainModelFromMp4(filepath string, name string) {

    haarCascade := gocv.NewCascadeClassifier()
    haarCascade.Load("aimodels/haarfrontalface.xml")
    defer haarCascade.Close()
    
    webcam, err := gocv.VideoCaptureFile(filepath)
    if err != nil {
	    fmt.Printf("Error opening video file: %v", err)
    }
    defer webcam.Close()

    // Loop to continuously read frames from the video file

    i := 0
    for {
	// Read a frame from the video file
	frame := gocv.NewMat()
	if ok := webcam.Read(&frame); !ok {
	    fmt.Println("Error reading frame from video file")
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
