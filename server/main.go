package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", connect)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, X-Auth-Token")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  921600,
		WriteBufferSize: 0,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "Upgrading error: %v", err)
		return
	}

	waitForImage(conn)
}

func waitForImage(conn *websocket.Conn) {
	for {
		img := getAndDecodeImage(conn)
		saveImage(img)
	}
}

func getAndDecodeImage(conn *websocket.Conn) image.Image {
	_, mssg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(mssg)
	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func saveImage(img image.Image) {
	rand.Seed(time.Now().UnixNano())
	file := createFile()
	png.Encode(file, img)
	file.Close()
}

func createFile() *os.File {
	var directory strings.Builder
	directory.Grow(30)
	directory.WriteString("userImages/")
	directory.WriteString("userImage")
	directory.WriteString(strconv.Itoa(rand.Intn(1000000) + 100000))
	directory.WriteString(".png")

	file, err := os.Create(directory.String())
	if err != nil {
		makeDirectory()
		file, _ = os.Create(directory.String())
	}

	return file
}

func makeDirectory() {
	err := os.Mkdir("UserImages", 0755)
	if err != nil {
		log.Fatal(err)
	}
}
