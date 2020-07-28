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

// properties for font and rectangle
var (
	fontSize      = 2.0
	fontThickness = 3
	fontColor     = color.RGBA{255, 0, 0, 0}
	rectColor     = color.RGBA{0, 0, 255, 0}
	rectThickness = 1
)

func main() {
	// address of the machinebox running on your local
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

	// returns image matrix to which your images will be written
	imgMat := gocv.NewMat()
	defer imgMat.Close()

	// initializes cascade classifier
	classifer := gocv.NewCascadeClassifier()
	classifer.Load(haarClassifier) // loads Haar cascade classifier
	defer classifer.Close()

	for {
		// reads image from the video frame and writes it to image matrix
		if ok := webcam.Read(&imgMat); !ok || imgMat.Empty() {
			log.Print("Faild to read image from webcam")
			continue
		}

		// initially we dont the person in the image
		faceName := "Unknown Person"

		// face detection using haar classifier from the image matrix
		imgRects := classifer.DetectMultiScale(imgMat)
		for _, imgRect := range imgRects {
			imgFace := imgMat.Region(imgRect)

			// IMWrite is used to capture current image of a person being read from webcam
			// commented after getting the images
			// uncomment and rerun it to get your images captured and
			// saved in the current directory as its being displayed in webcam
			// can be used for image training in your machinebox
			// if ok := gocv.IMWrite(fmt.Sprintf("imagename-%d.jpg", time.Now().Unix()), imgMat); !ok {
			// 	log.Print("Failed to write image")
			// }

			// IMEncode stores image matrix in mem buffer
			// will be read from buffer later
			imgBuf, err := gocv.IMEncode(".jpg", imgMat)
			if err != nil {
				log.Print("failed to encode image", err)
			}
			imgFace.Close()

			// reads from the buffer to match with the image trained in machinebox's facebox
			faces, err := machineboxClient.Check(bytes.NewBuffer(imgBuf))
			if err != nil {
				log.Println("failed to check image: ", err)
			}
			if len(faces) > 0 {
				faceName = faces[0].Name // name of the matched face/image
			}
			// configure image text properties to display
			textSize := gocv.GetTextSize(faceName, fontFace, fontSize, fontThickness)
			imgPt := image.Pt(imgRect.Min.X+(imgRect.Min.X/2)-(textSize.X/2), imgRect.Min.Y-2)
			// displays text for an face/image
			gocv.PutText(&imgMat, faceName, imgPt, fontFace, fontSize, fontColor, fontThickness)
			// draws rectangle around the face/image
			gocv.Rectangle(&imgMat, imgRect, rectColor, rectThickness)
		}

		window.IMShow(imgMat) // displays the face/image on screen
		window.WaitKey(1)     // must be called for image processing
	}
}
