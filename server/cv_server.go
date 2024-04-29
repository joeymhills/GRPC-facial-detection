package server

import (
    "bytes"
    "context"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "net"
    "os/exec"
    "time"
    "io"

    "gocv.io/x/gocv"
)

type Feature struct {
    BottomLeft image.Point
    TopLeft image.Point
    TopRight image.Point
    BottomRight image.Point
    FeatureType string
}

func(r *Feature) PrintPoints() {
    fmt.Printf("%s can be found at points: (%d,%d),(%d,%d),(%d,%d),(%d,%d)\n", r.FeatureType, r.BottomLeft.X, r.BottomLeft.Y,
	r.TopLeft.X, r.TopLeft.Y, r.TopRight.X, r.TopRight.Y, r.BottomRight.X, r.BottomRight.Y)
}

func(r *Feature) mean() *image.Point{
    
    avgx := r.BottomLeft.X + r.TopLeft.X + r.TopRight.X + r.BottomRight.X / 4
    avgy := r.BottomLeft.Y + r.TopLeft.Y + r.TopRight.Y + r.BottomRight.Y / 4
    
    fmt.Printf("%s can be found at points: (%d,%d),(%d,%d),(%d,%d),(%d,%d)\n", r.FeatureType, r.BottomLeft.X, r.BottomLeft.Y,
	r.TopLeft.X, r.TopLeft.Y, r.TopRight.X, r.TopRight.Y, r.BottomRight.X, r.BottomRight.Y)

    return &image.Point{avgx, avgy}
}

type Face struct {
    eyes Feature
    nose Feature
}


func FacialLayout(f *Face) error {
    
    return nil
}

func RectToFeature(r *image.Rectangle, f string) *Feature {    
    
    point0 := image.Point{r.Min.X, r.Min.Y}
    point1 := image.Point{r.Min.X, r.Max.Y}
    point2 := image.Point{r.Max.X, r.Max.Y}
    point3 := image.Point{r.Max.X, r.Min.Y}
    
    rect := &Feature{point0, point1, point2, point3, f}

    return rect
}
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
func TrainModel(name string) error {
    cmd := exec.Command("python3", "../python/train.py", name)
    output, err := cmd.CombinedOutput()
    if err != nil {
	return err
    }
    fmt.Println(output)

    return nil
}

func CheckFace(imgData *[]byte, model string) (bool, error) {
    cmd := exec.Command("python3", "python/test.py", model)
    _, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println("Error creating stdout pipe:", err)
    }
    ctx := context.Background()
    
    go func(ctx context.Context) {
	//err = cmd.Start()
	if err != nil {
	    fmt.Println("error spawning python socket script", err)
	}
    }(ctx)

    pythonHost := "127.0.0.1"
    pythonPort := "49522"
    
    var conn net.Conn

    for attempt := 1; attempt <= 4; attempt++ {
	fmt.Printf("Attempting to connect (attempt %d/%d)...\n", attempt, 10)
	
	conn, err = net.Dial("tcp", pythonHost+":"+pythonPort)
	if err == nil {
	    fmt.Println("Connection successful!")

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
	    attempt = 10
	    resultBool := boolBuf[0] != 0
	    return resultBool, nil
	}

	fmt.Printf("Connection failed: %v\n", err)
	time.Sleep(time.Second * 1) // Wait before the next attempt
	//retryDelay *= 2        // Exponential backoff (optional)

    }
	return false, nil
}


func GetFaceImages(imgData *[]byte) (faceNum int, faceImages *[]*[]byte, err error)  {

    haarCascade := gocv.NewCascadeClassifier()
    haarCascade.Load("aimodels/haarfrontalface.xml")
    defer haarCascade.Close()
    img, err := gocv.IMDecode(*imgData, gocv.IMReadColor)
    if err != nil {
	return 0, nil, err
    }
    if img.Empty() {
	return 0, nil, fmt.Errorf("Error loading image")
    }
    
    imgRects := haarCascade.DetectMultiScale(img)
    
    faceNum = 0
    var returnFaces []*[]byte
    //Iterates through faces and adds them to the returned array
    for _, rect := range imgRects {
	//Increases number of faces
	faceNum++
	//Draws rectangle around face
	gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)
	
	face, err := imageToBytes(img.Region(rect))
	if err != nil {
	    return 0, nil, err
	}
	_ = face
	returnFaces = append(returnFaces, face)
    }


    /* Save the image with landmarks
    gocv.IMWrite("output_image2.jpg", img)
    if err != nil {
	return 0, nil, err
    }
    */

    return faceNum, &returnFaces, nil
}
