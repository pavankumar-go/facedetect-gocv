package main

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/machinebox/sdk-go/facebox"
	"gocv.io/x/gocv"
)

const haarClassifier = "haar_classifier_face.xml"
const fontFace = gocv.FontHersheyPlain

var (
	fontSize      = 2.0
	fontThickness = 3
	fontColor     = color.RGBA{255, 0, 0, 0}
	rectColor     = color.RGBA{0, 0, 255, 0}
	rectThickness = 1
)

func main() {
	machineboxClient := facebox.New("http://localhost:8080")

	// Initialize Webcam
	// if using built-in webcam -> ID = 0
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalf("Failed to initialize webcam : %v", err)
	}
	defer webcam.Close()

	// Open a window to start webcam
	window := gocv.NewWindow("Face Detection")
	defer window.Close()

	// returns image matrix
	imgMat := gocv.NewMat()
	defer imgMat.Close()

	// initializes classifier
	classifer := gocv.NewCascadeClassifier()
	classifer.Load(haarClassifier) // using HAAR_face_default algorithm
	defer classifer.Close()

	for {
		// reads image from webcam and read it to image matrix
		if ok := webcam.Read(&imgMat); !ok || imgMat.Empty() {
			log.Print("Faild to read image from webcam")
			continue
		}
		// initially we dont the person in the image
		// Before training
		faceName := "Unknown Person"

		// image detection using haar classifier
		imgRects := classifer.DetectMultiScale(imgMat)
		for _, imgRect := range imgRects {

			imgFace := imgMat.Region(imgRect)

			/* IMWrite is used to capture current image of a person being read from webcam
			commented after getting the images
			uncomment and run it to get your images captured as its being displayed in webcam
			*/

			// if ok := gocv.IMWrite(fmt.Sprintf("pavan-%d.jpg", time.Now().Unix()), imgMat); !ok {
			// 	log.Print("Failed to write image")
			// }

			// IMEncode stores in mem buffer
			// can be read from buffer later
			imgBuf, err := gocv.IMEncode(".jpg", imgMat)
			if err != nil {
				log.Print("failed to encode image", err)
			}
			imgFace.Close()

			// uses the image buffer to check with the trained image using machine-box client
			faces, err := machineboxClient.Check(bytes.NewBuffer(imgBuf))
			if err != nil {
				log.Println("failed to check image: ", err)
			}
			if len(faces) > 0 {
				faceName = faces[0].Name // if there is face i.e., matched we show that name
			}

			textSize := gocv.GetTextSize(faceName, fontFace, fontSize, fontThickness) // configure text-properties to display
			imgPt := image.Pt(imgRect.Min.X+(imgRect.Min.X/2)-(textSize.X/2), imgRect.Min.Y-2)

			gocv.PutText(&imgMat, faceName, imgPt, fontFace, fontSize, fontColor, fontThickness) // displays text for a image
			gocv.Rectangle(&imgMat, imgRect, rectColor, rectThickness)                           // draws rectangle around the face

		}

		window.IMShow(imgMat) // displays the face on screen
		window.WaitKey(1)     // must be called for image processing
	}
}
