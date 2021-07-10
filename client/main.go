package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"./object"
	"./visual"

	"github.com/gorilla/websocket"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	conn := establishConnectionPrintResult()
	object.InitWindowAndToolbar()
	defer conn.Close()

	img := visual.CreateImage()
	waitForInput(conn, img)
}

func establishConnectionPrintResult() *websocket.Conn {
	conn, resp, err := dialServer()
	if err != nil {
		fmt.Printf("Error occured when connecting to server: %v\n", err)
		return nil
	}

	readBody(resp)
	return conn
}

func dialServer() (*websocket.Conn, *http.Response, error) {
	addr := flag.String("addr", "localhost:8080", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	return websocket.DefaultDialer.Dial(u.String(), nil)
}

func readBody(res *http.Response) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err, string(body))
	} else {
		fmt.Println("Connection established successfully")
	}
}

func waitForInput(conn *websocket.Conn, img *image.RGBA) {
	var (
		event       sdl.Event
		obj         object.SelectedObject
		ctrlPressed bool
	)

	for {
		event = sdl.WaitEvent()
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			ctrlPressed = processKeyboardEvent(t, conn, img, ctrlPressed)
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				processMouseButtonEvent(&obj, t, img)
			}
		}
	}
}

func processKeyboardEvent(keyboardEvent *sdl.KeyboardEvent, conn *websocket.Conn, img *image.RGBA, ctrlPressed bool) bool {
	if keyboardEvent.State == sdl.PRESSED {
		switch keyboardEvent.Keysym.Sym {
		case sdl.K_ESCAPE:
			os.Exit(0)
		case sdl.K_LCTRL, sdl.K_RCTRL:
			ctrlPressed = true
		case sdl.K_s:
			if ctrlPressed {
				sendImageToServer(conn, img)
				saveImageToUser(img)
			}
		}
	} else {
		return false
	}

	return ctrlPressed
}

func sendImageToServer(conn *websocket.Conn, img *image.RGBA) {
	if conn != nil {
		buf := new(bytes.Buffer)
		err := png.Encode(buf, img)
		if err != nil {
			log.Fatal(err)
		}
		byteImage := buf.Bytes()

		conn.WriteMessage(websocket.TextMessage, byteImage)
	}
}

func saveImageToUser(img *image.RGBA) {
	var fName string
	fmt.Print("Enter file name: ")
	fmt.Scan(&fName)

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	file := createFile(dirname, fName)
	png.Encode(file, img)
	file.Close()
}

func processMouseButtonEvent(obj *object.SelectedObject, mouseEvent *sdl.MouseButtonEvent, img *image.RGBA) {
	switch mouseEvent.Button {
	case sdl.BUTTON_LEFT:
		if mouseEvent.Y <= object.TOOLBAR_Y_POS {
			click := visual.Point{X: mouseEvent.X, Y: mouseEvent.Y}
			obj.ChangeObject(click)
		} else {
			drawOnMouseCoordinates(obj, img)
		}
	}
}

func drawOnMouseCoordinates(obj *object.SelectedObject, img *image.RGBA) {
	var lastClickCoordinates, currentClickCoordinates visual.Point

	for {
		event := sdl.WaitEvent()

		switch tp := event.(type) {
		case *sdl.MouseButtonEvent:
			currentClickCoordinates = visual.Point{X: tp.X, Y: tp.Y}
			obj.Draw(img, currentClickCoordinates)
			if tp.Type == sdl.MOUSEBUTTONUP {
				return
			}

		case *sdl.MouseMotionEvent:
			currentClickCoordinates = visual.Point{X: tp.X, Y: tp.Y}
			obj.DrawBetweenClicks(img, lastClickCoordinates, currentClickCoordinates)
			obj.Draw(img, currentClickCoordinates)
			lastClickCoordinates = currentClickCoordinates
		}
	}
}

func createFile(dirname, fName string) *os.File {
	var fullDir strings.Builder

	fullDir.Grow(50)
	fullDir.WriteString(dirname)
	fullDir.WriteString("/Downloads/")
	fullDir.WriteString(fName)
	fullDir.WriteString(".png")

	file, err := os.Create(fullDir.String())
	if err != nil {
		log.Fatal(err)
	}

	return file
}
