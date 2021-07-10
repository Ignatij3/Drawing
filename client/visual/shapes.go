package visual

import (
	"fmt"
	"image"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type (
	Figure interface {
		DrawToWindow(renderer *sdl.Renderer)
		DrawToImage(img *image.RGBA)
		SetShape()
		SetCoordinates(click Point)
		GetVerticalSize() int32
	}
	Rectangle struct {
		Center        Point
		Width, Height int32
	}
	Circle struct {
		Center Point
		Radius int32
	}
	Triangle struct {
		P1, P2, P3 Point
		Size       int32
	}
	Hexagon struct {
		Center Point
		Size   int32
	}
	Point struct {
		X, Y int32
	}
)

func (rect *Rectangle) DrawToWindow(renderer *sdl.Renderer) {
	renderer.FillRect(&sdl.Rect{X: (*rect).Center.X - ((*rect).Width / 2), Y: (*rect).Center.Y - ((*rect).Height / 2), W: (*rect).Width, H: (*rect).Height})
}

func (circ *Circle) DrawToWindow(renderer *sdl.Renderer) {
	newCircle := (*circ)
	for rad := (*circ).Radius; rad > 0; rad-- {
		newCircle.Radius = rad
		newCircle.drawCircle(renderer)
	}
}

func (circ Circle) drawCircle(renderer *sdl.Renderer) {
	p := make([]Point, 0)
	x, y := float64(circ.Radius), 0.0
	fi := 1.0 / float64(circ.Radius)
	cos, sin := math.Cos(fi), math.Sin(fi)
	for x > y {
		p = append(p, Point{int32(math.Round(x)), int32(math.Round(y))})
		x, y = x*cos-y*sin, x*sin+y*cos
	}

	pp := make([]sdl.Point, len(p))
	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X + v.X, Y: circ.Center.Y + v.Y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X - v.X, Y: circ.Center.Y + v.Y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X - v.Y, Y: circ.Center.Y + v.X}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X - v.Y, Y: circ.Center.Y - v.X}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X - v.X, Y: circ.Center.Y - v.Y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X + v.X, Y: circ.Center.Y - v.Y}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X + v.Y, Y: circ.Center.Y - v.X}
	}
	renderer.DrawLines(pp)

	for i, v := range p {
		pp[i] = sdl.Point{X: circ.Center.X + v.Y, Y: circ.Center.Y + v.X}
	}
	renderer.DrawLines(pp)
}

func (tri *Triangle) DrawToWindow(renderer *sdl.Renderer) {
	tri.sortPoints()

	if (*tri).P2.Y == (*tri).P3.Y {
		(*tri).fillBottomFlatTriangle(renderer)
	} else if (*tri).P1.Y == (*tri).P2.Y {
		(*tri).fillTopFlatTriangle(renderer)
	} else {
		p4 := Point{
			X: (*tri).P1.X + int32(
				(float64((*tri).P2.Y-(*tri).P1.Y)/float64((*tri).P3.Y-(*tri).P1.Y)))*
				((*tri).P3.X-(*tri).P1.X),
			Y: (*tri).P2.Y}

		(*tri).P3, p4 = p4, (*tri).P3
		(*tri).fillBottomFlatTriangle(renderer)
		(*tri).P1, (*tri).P2, (*tri).P3 = (*tri).P2, (*tri).P3, p4
		(*tri).fillTopFlatTriangle(renderer)
	}
}

func (tri *Triangle) sortPoints() {
	if (*tri).P1.Y > (*tri).P2.Y {
		(*tri).P1, (*tri).P2 = (*tri).P2, (*tri).P1
	}
	if (*tri).P2.Y > (*tri).P3.Y {
		(*tri).P2, (*tri).P3 = (*tri).P3, (*tri).P2
	}
	if (*tri).P1.Y > (*tri).P2.Y {
		(*tri).P1, (*tri).P2 = (*tri).P2, (*tri).P1
	}
}

func (tri Triangle) fillBottomFlatTriangle(renderer *sdl.Renderer) {
	var (
		slope1 float64 = float64(tri.P2.X-tri.P1.X) / float64(tri.P2.Y-tri.P1.Y)
		slope2 float64 = float64(tri.P3.X-tri.P1.X) / float64(tri.P3.Y-tri.P1.Y)
		curx1  float64 = float64(tri.P1.X)
		curx2  float64 = float64(tri.P1.X)
	)

	for PointY := tri.P1.Y; PointY <= tri.P2.Y; PointY++ {
		renderer.DrawLine(int32(curx1), PointY, int32(curx2), PointY)
		curx1 += slope1
		curx2 += slope2
	}
}

func (tri Triangle) fillTopFlatTriangle(renderer *sdl.Renderer) {
	var (
		slope1 float64 = float64(tri.P3.X-tri.P1.X) / float64(tri.P3.Y-tri.P1.Y)
		slope2 float64 = float64(tri.P3.X-tri.P2.X) / float64(tri.P3.Y-tri.P2.Y)
		curx1  float64 = float64(tri.P3.X)
		curx2  float64 = float64(tri.P3.X)
	)

	for PointY := tri.P3.Y; PointY > tri.P1.Y; PointY-- {
		renderer.DrawLine(int32(curx1), PointY, int32(curx2), PointY)
		curx1 -= slope1
		curx2 -= slope2
	}
}

func (hex *Hexagon) DrawToWindow(renderer *sdl.Renderer) {
	var (
		b, c Point
		tri  Triangle
	)

	for angle := 60; angle <= 360; angle += 60 {
		b.X = int32(math.Cos(float64(angle-60)*(math.Pi/180.0))*float64((*hex).Size)) + (*hex).Size + (*hex).Center.X
		b.Y = int32(math.Sin(float64(angle-60)*(math.Pi/180.0))*float64((*hex).Size)) + (*hex).Size + (*hex).Center.Y
		c.X = int32(math.Cos(float64(angle)*(math.Pi/180.0))*float64((*hex).Size)) + (*hex).Size + (*hex).Center.X
		c.Y = int32(math.Sin(float64(angle)*(math.Pi/180.0))*float64((*hex).Size)) + (*hex).Size + (*hex).Center.Y
		tri = Triangle{P1: Point{(*hex).Center.X + (*hex).Size, (*hex).Center.Y + (*hex).Size}, P2: b, P3: c, Size: 0}
		tri.DrawToWindow(renderer)
	}
}

func (rect *Rectangle) SetShape() {
	var size int32

	fmt.Print("Enter width: ")
	fmt.Scan(&size)
	(*rect).Width = size

	fmt.Print("Enter height: ")
	fmt.Scan(&size)
	(*rect).Height = size
}

func (circ *Circle) SetShape() {
	var rad int32

	fmt.Print("Enter radius: ")
	fmt.Scan(&rad)
	(*circ).Radius = rad
}

func (tri *Triangle) SetShape() {
	var size int32

	fmt.Print("Enter size: ")
	fmt.Scan(&size)
	(*tri).Size = size
}

func (hex *Hexagon) SetShape() {
	var size int32

	fmt.Print("Enter size: ")
	fmt.Scan(&size)
	(*hex).Size = size
}

func (rect *Rectangle) SetCoordinates(click Point) {
	(*rect).Center.X = click.X
	(*rect).Center.Y = click.Y
}

func (circ *Circle) SetCoordinates(click Point) {
	(*circ).Center.X = click.X
	(*circ).Center.Y = click.Y
}

func (tri *Triangle) SetCoordinates(click Point) {
	(*tri).P1.X, (*tri).P1.Y = click.X, click.Y-(*tri).Size
	(*tri).P2.X, (*tri).P2.Y = click.X-(*tri).Size, click.Y+(*tri).Size
	(*tri).P3.X, (*tri).P3.Y = click.X+(*tri).Size, click.Y+(*tri).Size
}

func (hex *Hexagon) SetCoordinates(click Point) {
	(*hex).Center.X = click.X - (*hex).Size
	(*hex).Center.Y = click.Y - (*hex).Size
}

func (rect Rectangle) GetVerticalSize() int32 {
	return rect.Center.Y - (rect.Height / 2)
}

func (circ Circle) GetVerticalSize() int32 {
	return circ.Center.Y - circ.Radius
}

func (tri Triangle) GetVerticalSize() int32 {
	return tri.P1.Y - (tri.Size / 2)
}

func (hex Hexagon) GetVerticalSize() int32 {
	return hex.Center.Y - (hex.Size / 2)
}
