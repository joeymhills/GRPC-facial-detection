package server

import (
    "fmt"
    "gocv.io/x/gocv"
    "image/color"
)

func InitOpencv() error {

    //var landmarkData [][]float32
    //var points []image.Point

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

    img := gocv.IMRead("img/lebron1.jpg", gocv.IMReadColor)
    if img.Empty() {
	return fmt.Errorf("Error loading image")
    }


    imgRects := faces.DetectMultiScale(img)

    for _, rect := range imgRects {
	//draws rectangle around face
	gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)
	face := img.Region(rect)
	
	//draws rectangle around nose
	noseBox := nose.DetectMultiScale(face)
	gocv.Rectangle(&face, noseBox[0], color.RGBA{255, 0, 0, 0}, 2)

	//draws rectangle around eyes
	eyes := eyes.DetectMultiScale(face)
	gocv.Rectangle(&face, eyes[0], color.RGBA{255, 0, 0, 0}, 2)

	//draws rectangle around mouth
	mouthBox := mouth.DetectMultiScale(face)
	gocv.Rectangle(&face, mouthBox[0], color.RGBA{255, 0, 0, 0}, 2)
    }

    // Save the image with landmarks
    gocv.IMWrite("output_image2.jpg", img)

    return nil
}
