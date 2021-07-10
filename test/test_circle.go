package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type (
	circle struct {
		center fpoint
		radius int32
	}
	point struct {
		x, y int32
	}
	fpoint struct {
		x, y float32
	}
)

func (c circle) drawCircleInt(renderer *sdl.Renderer) {
	p := make([]point, 0)
	x, y := float64(c.radius), 0.0
	fi := 1.0 / float64(c.radius)
	cos, sin := math.Cos(fi), math.Sin(fi)
	for x > y {
		p = append(p, point{int32(math.Round(x)), int32(math.Round(y))})
		x, y = x*cos-y*sin, x*sin+y*cos
	}

	pp := make([]sdl.Point, len(p))
	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) + v.x, int32(c.center.y) + v.y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) - v.x, int32(c.center.y) + v.y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) - v.y, int32(c.center.y) + v.x}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) - v.y, int32(c.center.y) - v.x}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) - v.x, int32(c.center.y) - v.y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) + v.x, int32(c.center.y) - v.y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) + v.y, int32(c.center.y) - v.x}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{int32(c.center.x) + v.y, int32(c.center.y) + v.x}
	}
	renderer.DrawLines(pp)
}

func (c circle) drawCircleFloat(renderer *sdl.Renderer) {
	p := make([]fpoint, 0)
	x, y := float64(c.radius), 0.0
	fi := 1.0 / float64(c.radius)
	cos, sin := math.Cos(fi), math.Sin(fi)
	for x > y {
		p = append(p, fpoint{float32(math.Round(x)), float32(math.Round(y))})
		x, y = x*cos-y*sin, x*sin+y*cos
	}

	pp := make([]sdl.FPoint, len(p))
	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x + v.x), float32(c.center.y + v.y)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x - v.x), float32(c.center.y + v.y)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x - v.y), float32(c.center.y + v.x)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x - v.y), float32(c.center.y - v.x)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x - v.x), float32(c.center.y - v.y)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x + v.x), float32(c.center.y - v.y)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x + v.y), float32(c.center.y - v.x)}
	}
	renderer.DrawLinesF(pp)

	for i, v := range p {
		pp[i] = sdl.FPoint{float32(c.center.x + v.y), float32(c.center.y + v.x)}
	}
	renderer.DrawLinesF(pp)
}

func main() {
	window, _ := sdl.CreateWindow("Test", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1280, 720, sdl.WINDOW_SHOWN)
	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()
	renderer.SetDrawColor(0, 0, 0, 255)
	var circ circle

	circ.center.x = 100
	circ.center.y = 100
	for rad := int32(30); rad > 0; rad-- {
		circ.radius = rad
		circ.drawCircleInt(renderer)
	}

	circ.center.x = 200
	circ.center.y = 100
	for rad := int32(30); rad > 0; rad-- {
		circ.radius = rad
		circ.drawCircleFloat(renderer)
	}

	renderer.Present()
	sdl.Delay(10000)
}
