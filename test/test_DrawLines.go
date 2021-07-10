package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func testLines(renderer *sdl.Renderer) { //Пытался понять как работает функция renderer.DrawLines(). Понял.
	renderer.SetDrawColor(0, 0, 0, 255)
	var lines1 [1]sdl.Point
	var lines2 [2]sdl.Point
	var lines3 [3]sdl.Point
	var lines4 [4]sdl.Point

	lines1 = [1]sdl.Point{sdl.Point{10, 10}}
	lines2 = [2]sdl.Point{sdl.Point{10, 20}, sdl.Point{100, 20}}
	lines3 = [3]sdl.Point{sdl.Point{10, 30}, sdl.Point{100, 30}, sdl.Point{200, 40}}
	lines4 = [4]sdl.Point{sdl.Point{10, 50}, sdl.Point{100, 50}, sdl.Point{200, 60}, sdl.Point{300, 60}}
	renderer.DrawLines(lines1[:])
	renderer.DrawLines(lines2[:])
	renderer.DrawLines(lines3[:])
	renderer.DrawLines(lines4[:])
	renderer.Present()
}

func main() {
	window, err := sdl.CreateWindow("Test", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1280, 720, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()
	defer window.Destroy()
	defer renderer.Destroy()
	testLines(renderer)
	sdl.Delay(10000)
}
