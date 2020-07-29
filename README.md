Face Detection and Recognition Using GoCV
=========
## A simple face detection and recognition program written in golang by making use of Haar feature based cascade classifier
### Requirements for Running it on Mac OS
* Go lang - `version go1.13.6 darwin/amd64` (https://golang.org/dl/)
* Docker for Mac or Docker Desktop (https://docs.docker.com/docker-for-mac/install)
* OpenCV 3 - `brew install opencv3`
* Go Package Dependencies 
   * `go get gocv.io/x/gocv`
   * `go get github.com/machinebox/go-sdk/facebox`

### Steps to Run 
1. Clone the repo 
2. Register in machine-box (https://machinebox.io)
3. Get Machine-Box Access Key, looks something like this `MB_KEY=xxxxxxxxxxx....`
4. Run MachineBox locally 
    *  export your machine-box key `export MB_KEY="xxxxx...."`
    *  run machinebox using docker `docker run -d -p 8080:8080 -e "MB_KEY=$MB_KEY" machinebox/facebox` remove `-d` if you dont want machine-box container to run in detached mode
5. Now check machinebox by hitting `http://localhost:8080`
6. Upload your image for training under `Post a file` -> `Try it now`
7. cd to `/your/directory/main.go` and `go run main.go` in your terminal

### FAQ 
How to improve recognition accuracy
There are few suggestions:
* Purely uses Haar Cascade clssifier default face xml, hence no control over the algorithm
* Upload more samples of your faces in different postures to machinebox's facebox for better recognition
* Try to test it under sufficient lighting environment where your face is clearly visible, that said avoid low-light and less grain or noise when using webcam

### Useful links for references
1. https://github.com/opencv/opencv/tree/master/data/haarcascades  - Haar cascade classifier xml files
2. https://gocv.io
3. https://machinebox.io