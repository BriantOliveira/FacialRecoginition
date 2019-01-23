package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tfacedetect [camera ID] [classifier XML file]")
		return
	}
	// Parse args
	deviceID, _ := strconv.Atoi(os.Args[1])
	xmlFile := os.Args[2]

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// Open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// Prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// Color for the rect when faces detected
	blue := color.RGBA{ 0, 0, 255, 0}

	// Load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("Start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}
		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("Found %d faces\n", len(rects))

		// Draw a retangle around each found face
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		// Show the image in the window and wait 1 mins
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}