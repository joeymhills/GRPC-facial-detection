package server

import (
    "bytes"
    "database/sql"
    "fmt"
    "image"
    "image/jpeg"
    "os"

    _ "github.com/go-sql-driver/mysql"

    "gocv.io/x/gocv"
)

//Test function
func TestModelFromImage(db *sql.DB, filePath string) {
    
    file, err := os.Open(filePath)
    if err != nil {
	fmt.Println("Error opening image file:", err)
	return
    }
    defer file.Close()

    // Read the image data into a byte slice
    imageData, err := jpeg.Decode(file)
    if err != nil {
	fmt.Println("Error reading image file:", err)
	return
    }

    // Encode the image into JPEG format and write it to the buffer
    buf := new(bytes.Buffer)
    err = jpeg.Encode(buf, imageData, nil)
    if err != nil {
	fmt.Println("Error encoding image:", err)
	return
    }
    imageBytes := buf.Bytes()

    err = HandleImage(db, &imageBytes)
    if err != nil {
	fmt.Println("error when starting opencv", err)
    }
}

// For extracting face images from pictures for test data
func ProcessTestImages() error {

    haarCascade := gocv.NewCascadeClassifier()
    haarCascade.Load("aimodels/haarfacetree.xml")
    defer haarCascade.Close()

    // Get a list of files from the tensorflow models directory
    files, err := os.ReadDir("img")
    if err != nil {
	fmt.Print(err)
    }

    // Slice to store filenames
    var imgs []string

    // Iterate over the files and add their names to the slice
    for _, file := range files {
	if file.IsDir(){
	    continue
	}

	imgs = append(imgs, file.Name())
	imagePath := fmt.Sprintf("img/%s", file.Name())

	imgData, err := os.ReadFile(imagePath)
	if err != nil {
	    return  fmt.Errorf("err decoding image:", err)
	}

	img, err := gocv.IMDecode(imgData, gocv.IMReadColor)
	if err != nil {
	    return  fmt.Errorf("err decoding image:", err)
	}
	if img.Empty() {
	    return  fmt.Errorf("Error loading image")
	}
	fmt.Println("image loaded:", file.Name())
	
	imgRects := haarCascade.DetectMultiScaleWithParams(img, 3.5, 1, 0, image.Point{60, 60}, image.Point{1500,1500})

	var faceIdx int
	//Iterates through faces and adds them to the returned array
	for _, rect := range imgRects {
	    //Increases number of faces
	    faceIdx++

	    face := img.Region(rect)

	    imageTitle := fmt.Sprintf("python/trainimg/class_0/%d.jpg", faceIdx)
	    gocv.IMWrite(imageTitle, face)

	}
    }
    return nil
}

