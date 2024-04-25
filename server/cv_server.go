package server

import (
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "bytes"

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

func HaarCascade(imgData *[]byte) (faceNum int, markedImage *[]byte, err error)  {

    faces := gocv.NewCascadeClassifier()
    faces.Load("aimodels/haarfrontalface.xml")
    defer faces.Close()

    nose := gocv.NewCascadeClassifier()
    nose.Load("aimodels/haarnose.xml")
    defer nose.Close()

    leftEye := gocv.NewCascadeClassifier()
    leftEye.Load("aimodels/haarlefteye.xml")
    defer leftEye.Close()

    rightEye := gocv.NewCascadeClassifier()
    rightEye.Load("aimodels/haarrighteye.xml")
    defer rightEye.Close()

    eyes := gocv.NewCascadeClassifier()
    eyes.Load("aimodels/haareyes.xml")
    defer eyes.Close()

    mouth := gocv.NewCascadeClassifier()
    mouth.Load("aimodels/haarmouth.xml")
    defer mouth.Close()

    img, err := gocv.IMDecode(*imgData, gocv.IMReadColor)
    if err != nil {
	return 0, nil, err
    }
    if img.Empty() {
	return 0, nil, fmt.Errorf("Error loading image")
    }
    
    imgRects := faces.DetectMultiScale(img)

    for _, rect := range imgRects {
	//draws rectangle around face
	gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)
	
	face := img.Region(rect)
	
	//draws rectangle around nose
	noseBox := nose.DetectMultiScale(face)
	gocv.Rectangle(&face, noseBox[0], color.RGBA{255, 0, 0, 0}, 2)
	points := RectToFeature(&noseBox[0], "nose")
	points.PrintPoints()

	//draws rectangle around eyes
	eyes := eyes.DetectMultiScale(face)
	gocv.Rectangle(&face, eyes[0], color.RGBA{255, 0, 0, 0}, 2)
	
	// Region of interest for the lower half of the face
	//halfFace := img.Region(image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, avg))

	//draws rectangle around mouth
	mouthBox := mouth.DetectMultiScale(face)
	gocv.Rectangle(&face, mouthBox[0], color.RGBA{0, 0, 255, 0}, 2)
    }

    // Save the image with landmarks
    gocv.IMWrite("output_image2.jpg", img)
    if err != nil {
	return 0, nil, err
    }
    
    //Convert mat object to []byte
    img.ToImage()
    var buf bytes.Buffer
    _ = jpeg.Encode(&buf, img, nil)
    returnImg := buf.Bytes()

    return len(imgRects), &returnImg, err
}
