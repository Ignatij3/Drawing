package visual

import (
	"image"
	"image/color"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH, WINDOW_HEIGHT int   = 1280, 820
	TOOLBAR_Y_POS               int32 = 100
)

//Used to not pass sdl.Color variable through almost all functions
var CurrentColor sdl.Color

func CreateImage() *image.RGBA {
	return image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{WINDOW_WIDTH, WINDOW_HEIGHT - int(TOOLBAR_Y_POS)}})
}

func (rect Rectangle) DrawToImage(img *image.RGBA) {
	var a Point = Point{X: rect.Center.X - (rect.Width / 2), Y: rect.Center.Y - (rect.Height / 2)}
	var b Point = Point{X: rect.Center.X + (rect.Width / 2), Y: rect.Center.Y + (rect.Height / 2)}
	for y := a.Y; y <= b.Y; y++ {
		drawImageLine(img, Point{a.X, y}, Point{b.X, y})
	}
}

func (circ Circle) DrawToImage(img *image.RGBA) {
	for rad := circ.Radius; rad > 0; rad-- {
		circ.Radius = rad
		circ.drawImageCircle(img)
	}
}

func (circ Circle) drawImageCircle(img *image.RGBA) {
	p := make([]Point, 0)
	x, y := float64(circ.Radius), 0.0
	fi := 1.0 / float64(circ.Radius)
	cos, sin := math.Cos(fi), math.Sin(fi)
	for x > y {
		p = append(p, Point{int32(math.Round(x)), int32(math.Round(y))})
		x, y = x*cos-y*sin, x*sin+y*cos
	}

	pp := make([]Point, len(p))
	for i, v := range p {
		pp[i] = Point{circ.Center.X + v.X, circ.Center.Y + v.Y}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X - v.X, circ.Center.Y + v.Y}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X - v.Y, circ.Center.Y + v.X}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X - v.Y, circ.Center.Y - v.X}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X - v.X, circ.Center.Y - v.Y}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X + v.X, circ.Center.Y - v.Y}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X + v.Y, circ.Center.Y - v.X}
	}
	drawImageLines(img, pp)

	for i, v := range p {
		pp[i] = Point{circ.Center.X + v.Y, circ.Center.Y + v.X}
	}
	drawImageLines(img, pp)
}

func drawImageLines(img *image.RGBA, P1 []Point) {
	for n := 1; n < len(P1); n++ {
		drawImageLine(img, Point{P1[n].X, P1[n].Y}, Point{P1[n-1].X, P1[n-1].Y})
	}
}

func (tri Triangle) DrawToImage(img *image.RGBA) {
	tri.sortPoints()

	if tri.P2.Y == tri.P3.Y {
		tri.fillImageBottomFlatTriangle(img)
	} else if tri.P1.Y == tri.P2.Y {
		tri.fillImageTopFlatTriangle(img)
	} else {
		p4 := Point{X: tri.P1.X + int32((float64(tri.P2.Y-tri.P1.Y)/float64(tri.P3.Y-tri.P1.Y)))*(tri.P3.X-tri.P1.X), Y: tri.P2.Y}

		tri.P3, p4 = p4, tri.P3
		tri.fillImageBottomFlatTriangle(img)
		tri.P1, tri.P2, tri.P3 = tri.P2, tri.P3, p4
		tri.fillImageTopFlatTriangle(img)
	}
}

func (tri Triangle) fillImageBottomFlatTriangle(img *image.RGBA) {
	var (
		slope1 float64 = float64(tri.P2.X-tri.P1.X) / float64(tri.P2.Y-tri.P1.Y)
		slope2 float64 = float64(tri.P3.X-tri.P1.X) / float64(tri.P3.Y-tri.P1.Y)
		curx1  float64 = float64(tri.P1.X)
		curx2  float64 = float64(tri.P1.X)
	)

	for PointY := tri.P1.Y; PointY <= tri.P2.Y; PointY++ {
		drawImageLine(img, Point{int32(curx1), PointY}, Point{int32(curx2), PointY})
		curx1 += slope1
		curx2 += slope2
	}
}

func (tri Triangle) fillImageTopFlatTriangle(img *image.RGBA) {
	var (
		slope1 float64 = float64(tri.P3.X-tri.P1.X) / float64(tri.P3.Y-tri.P1.Y)
		slope2 float64 = float64(tri.P3.X-tri.P2.X) / float64(tri.P3.Y-tri.P2.Y)
		curx1  float64 = float64(tri.P3.X)
		curx2  float64 = float64(tri.P3.X)
	)

	for PointY := tri.P3.Y; PointY > tri.P1.Y; PointY-- {
		drawImageLine(img, Point{int32(curx1), PointY}, Point{int32(curx2), PointY})
		curx1 -= slope1
		curx2 -= slope2
	}
}

func (hex Hexagon) DrawToImage(img *image.RGBA) {
	var (
		b, c Point
		tri  Triangle
	)

	//Получаю координаты углов и по ним рисую 6 треугольников, первая точка всегда центр
	for angle := 60; angle <= 360; angle += 60 {
		b.X = int32(math.Cos(float64(angle-60)*(math.Pi/180.0))*float64(hex.Size)) + hex.Size + hex.Center.X
		b.Y = int32(math.Sin(float64(angle-60)*(math.Pi/180.0))*float64(hex.Size)) + hex.Size + hex.Center.Y
		c.X = int32(math.Cos(float64(angle)*(math.Pi/180.0))*float64(hex.Size)) + hex.Size + hex.Center.X
		c.Y = int32(math.Sin(float64(angle)*(math.Pi/180.0))*float64(hex.Size)) + hex.Size + hex.Center.Y
		tri = Triangle{P1: Point{hex.Center.X + hex.Size, hex.Center.Y + hex.Size}, P2: b, P3: c, Size: 0}
		tri.DrawToImage(img)
	}
}

func drawImageLine(img *image.RGBA, a, b Point) {
	a.Y -= TOOLBAR_Y_POS
	b.Y -= TOOLBAR_Y_POS
	if a.Y < 0 && b.Y < 0 {
		return
	}

	if a.Y < 0 {
		a.Y = 0
	}
	if b.Y < 0 {
		b.Y = 0
	}

	if math.Abs(float64(b.Y-a.Y)) < math.Abs(float64(b.X-a.X)) {
		if a.X > b.X {
			drawImageLineLow(img, b, a)
		} else {
			drawImageLineLow(img, a, b)
		}
	} else {
		if a.Y > b.Y {
			drawImageLineHigh(img, b, a)
		} else {
			drawImageLineHigh(img, a, b)
		}
	}
}

func drawImageLineLow(img *image.RGBA, a, b Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	yi := int32(1)
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := (2 * dy) - dx

	col := color.RGBA{CurrentColor.R, CurrentColor.G, CurrentColor.B, CurrentColor.A}
	for ; a.X < b.X; a.X++ {
		img.Set(int(a.X), int(a.Y), col)
		if D > 0 {
			a.Y += yi
			D += 2 * -dx
		}
		D += 2 * dy
	}
}

func drawImageLineHigh(img *image.RGBA, a, b Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	xi := int32(1)
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := (2 * dx) - dy

	col := color.RGBA{CurrentColor.R, CurrentColor.G, CurrentColor.B, CurrentColor.A}
	for ; a.Y < b.Y; a.Y++ {
		img.Set(int(a.X), int(a.Y), col)
		if D > 0 {
			a.X += xi
			D += 2 * -dy
		}
		D += 2 * dx
	}
}
