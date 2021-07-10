package object

import (
	"fmt"
	"image"
	"math"
	"sync"

	"../visual"

	"github.com/veandco/go-sdl2/sdl"
)

type (
	SelectedObject struct {
		Shape            visual.Figure
		LastCheckedColor int
		LastCheckedShape int
	}
	toolHitbox struct {
		upperLeftCorner visual.Point
		toolColor       sdl.Color
	}
)

const (
	WINDOW_WIDTH, WINDOW_HEIGHT int32 = 1280, 820
	TOOLBAR_Y_POS               int32 = 100
	AMOUNT_TOOLS                      = 14
	AMOUNT_SHAPES                     = 4
	TOOL_FRAME_SIDE                   = 80
)

var (
	hitbox [14]toolHitbox

	renderer     *sdl.Renderer
	rendererSync sync.WaitGroup

	red    = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	orange = sdl.Color{R: 255, G: 145, B: 0, A: 255}
	yellow = sdl.Color{R: 255, G: 255, B: 0, A: 255}
	green  = sdl.Color{R: 0, G: 255, B: 0, A: 255}
	cyan   = sdl.Color{R: 0, G: 255, B: 255, A: 255}
	blue   = sdl.Color{R: 0, G: 0, B: 255, A: 255}
	purple = sdl.Color{R: 150, G: 0, B: 255, A: 255}
	white  = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	brown  = sdl.Color{R: 60, G: 30, B: 0, A: 255}
	black  = sdl.Color{R: 0, G: 0, B: 0, A: 255}
)

func drawToolbar(obj SelectedObject) {
	cleartoolBar()

	for toolPosition := 0; toolPosition < AMOUNT_TOOLS; toolPosition++ {
		obj.redrawFrame(&toolPosition)
		drawOptions(toolPosition)
	}
}

func cleartoolBar() {
	renderer.SetDrawColor(white.R, white.G, white.B, white.A)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: WINDOW_WIDTH, H: TOOLBAR_Y_POS})
	renderer.SetDrawColor(black.R, black.G, black.B, black.A)
	renderer.DrawLine(0, TOOLBAR_Y_POS, WINDOW_WIDTH, TOOLBAR_Y_POS)
}

func (object *SelectedObject) redrawFrame(toolPosition *int) {
	//Comparing if the number sent to function occupies the same memory location as in the object
	//If it does: paint the tool's frame black as it is not used anymore
	//If it does not: paint the tool's frame yellow
	//Look how it used in redrawFrameSetPosition()

	if &(*object).LastCheckedShape == toolPosition || &(*object).LastCheckedColor == toolPosition {
		renderer.SetDrawColor(black.R, black.G, black.B, black.A)

	} else if (*toolPosition < AMOUNT_SHAPES && *toolPosition == object.LastCheckedShape) ||
		(*toolPosition >= AMOUNT_SHAPES && *toolPosition == object.LastCheckedColor) {
		renderer.SetDrawColor(yellow.R, yellow.G, yellow.B, yellow.A)

	} else {
		renderer.SetDrawColor(black.R, black.G, black.B, black.A)
	}

	renderer.DrawRect(&sdl.Rect{X: hitbox[*toolPosition].upperLeftCorner.X - 1, Y: hitbox[*toolPosition].upperLeftCorner.Y - 1,
		W: TOOL_FRAME_SIDE + 2, H: TOOL_FRAME_SIDE + 2})
	renderer.SetDrawColor(visual.CurrentColor.R, visual.CurrentColor.G, visual.CurrentColor.B, visual.CurrentColor.A)
}

func drawOptions(toolPosition int) {
	renderer.SetDrawColor(hitbox[toolPosition].toolColor.R, hitbox[toolPosition].toolColor.G, hitbox[toolPosition].toolColor.B, hitbox[toolPosition].toolColor.A)

	if toolPosition >= AMOUNT_SHAPES {
		renderer.FillRect(&sdl.Rect{X: hitbox[toolPosition].upperLeftCorner.X, Y: hitbox[toolPosition].upperLeftCorner.Y,
			W: TOOL_FRAME_SIDE, H: TOOL_FRAME_SIDE})
	} else {
		switch toolPosition {
		case 0:
			rectangle := visual.Rectangle{
				Center: visual.Point{X: hitbox[0].upperLeftCorner.X + TOOL_FRAME_SIDE/2, Y: hitbox[0].upperLeftCorner.Y + TOOL_FRAME_SIDE/2},
				Width:  TOOL_FRAME_SIDE, Height: TOOL_FRAME_SIDE}
			rectangle.DrawToWindow(renderer)
		case 1:
			circle := visual.Circle{
				Center: visual.Point{X: hitbox[1].upperLeftCorner.X + TOOL_FRAME_SIDE/2, Y: hitbox[1].upperLeftCorner.Y + TOOL_FRAME_SIDE/2},
				Radius: TOOL_FRAME_SIDE / 2}
			circle.DrawToWindow(renderer)
		case 2:
			triangle := visual.Triangle{
				P1: visual.Point{X: hitbox[2].upperLeftCorner.X + TOOL_FRAME_SIDE/2, Y: hitbox[2].upperLeftCorner.Y},
				P2: visual.Point{X: hitbox[2].upperLeftCorner.X, Y: hitbox[2].upperLeftCorner.Y + TOOL_FRAME_SIDE},
				P3: visual.Point{X: hitbox[2].upperLeftCorner.X + TOOL_FRAME_SIDE, Y: hitbox[2].upperLeftCorner.Y + TOOL_FRAME_SIDE}, Size: 0}
			triangle.DrawToWindow(renderer)
		case 3:
			hexagon := visual.Hexagon{
				Center: visual.Point{X: hitbox[3].upperLeftCorner.X, Y: hitbox[3].upperLeftCorner.Y},
				Size:   TOOL_FRAME_SIDE / 2}
			hexagon.DrawToWindow(renderer)
		}
	}
}

func (object *SelectedObject) ChangeObject(click visual.Point) {
	if ok, toolPosition := checkToolOnClickReturnPosition(click); ok {
		object.redrawFrameSetPosition(toolPosition)
		object.changeTool(toolPosition, click)
	}
	renderer.Present()
}

func checkToolOnClickReturnPosition(a visual.Point) (bool, int) {
	for n, c := range hitbox {
		if a.X >= c.upperLeftCorner.X && a.X <= c.upperLeftCorner.X+TOOL_FRAME_SIDE &&
			a.Y >= c.upperLeftCorner.Y && a.Y <= c.upperLeftCorner.Y+TOOL_FRAME_SIDE {
			return true, n
		}
	}
	return false, 0
}

func (object *SelectedObject) redrawFrameSetPosition(toolPosition int) {
	if toolPosition < AMOUNT_SHAPES {
		object.redrawFrame(&(*object).LastCheckedShape)
		(*object).LastCheckedShape = toolPosition
		object.redrawFrame(&toolPosition)
	} else {
		object.redrawFrame(&(*object).LastCheckedColor)
		(*object).LastCheckedColor = toolPosition
		object.redrawFrame(&toolPosition)
	}
}

func (object *SelectedObject) changeTool(toolPosition int, click visual.Point) {
	switch toolPosition {
	case 0, 1, 2, 3:
		(*object).setNewShape(toolPosition, click)
	case 4:
		renderer.SetDrawColor(red.R, red.G, red.B, red.A)
		visual.CurrentColor = red
	case 5:
		renderer.SetDrawColor(orange.R, orange.G, orange.B, orange.A)
		visual.CurrentColor = orange
	case 6:
		renderer.SetDrawColor(yellow.R, yellow.G, yellow.B, yellow.A)
		visual.CurrentColor = yellow
	case 7:
		renderer.SetDrawColor(green.R, green.G, green.B, green.A)
		visual.CurrentColor = green
	case 8:
		renderer.SetDrawColor(cyan.R, cyan.G, cyan.B, cyan.A)
		visual.CurrentColor = cyan
	case 9:
		renderer.SetDrawColor(blue.R, blue.G, blue.B, blue.A)
		visual.CurrentColor = blue
	case 10:
		renderer.SetDrawColor(purple.R, purple.G, purple.B, purple.A)
		visual.CurrentColor = purple
	case 11:
		renderer.SetDrawColor(white.R, white.G, white.B, white.A)
		visual.CurrentColor = white
	case 12:
		renderer.SetDrawColor(brown.R, brown.G, brown.B, brown.A)
		visual.CurrentColor = brown
	case 13:
		renderer.SetDrawColor(black.R, black.G, black.B, black.A)
		visual.CurrentColor = black
	}
	object.changeLastCheckedColorOrShape(toolPosition)
}

func (object *SelectedObject) setNewShape(toolPosition int, click visual.Point) {
	switch toolPosition {
	case 0:
		(*object).Shape = new(visual.Rectangle)
	case 1:
		(*object).Shape = new(visual.Circle)
	case 2:
		(*object).Shape = new(visual.Triangle)
	case 3:
		(*object).Shape = new(visual.Hexagon)
	}

	(*object).Shape.SetShape()
}

func (object *SelectedObject) changeLastCheckedColorOrShape(toolPosition int) {
	if toolPosition < AMOUNT_SHAPES {
		(*object).LastCheckedShape = toolPosition
	} else {
		(*object).LastCheckedColor = toolPosition
	}
}

func (object *SelectedObject) DrawBetweenClicks(img *image.RGBA, a, b visual.Point) {
	if (*object).Shape == nil {
		fmt.Println("Shape is not chosen, unable to draw")
	} else if a.X != 0 && a.Y != 0 {
		if math.Abs(float64(b.Y-a.Y)) < math.Abs(float64(b.X-a.X)) {
			if a.X > b.X {
				object.drawShapeLow(img, b, a)
			} else {
				object.drawShapeLow(img, a, b)
			}
		} else {
			if a.Y > b.Y {
				object.drawShapeHigh(img, b, a)
			} else {
				object.drawShapeHigh(img, a, b)
			}
		}
	}
}

func (object *SelectedObject) drawShapeLow(img *image.RGBA, a, b visual.Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	yi := int32(1)
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := (2 * dy) - dx

	var (
		counter    int
		screenSync *sync.Mutex = new(sync.Mutex)
	)

	for ; a.X < b.X; a.X++ {
		rendererSync.Add(1)
		go object.drawShapeToWindowAndImage(img, a, screenSync)

		if D > 0 {
			a.Y += yi
			D += 2 * -dx
		}
		D += 2 * dy

		if counter++; counter >= 10 {
			rendererSync.Wait()
			counter = 0
		}
	}

	rendererSync.Wait()
	renderer.Present()
}

func (object *SelectedObject) drawShapeHigh(img *image.RGBA, a, b visual.Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	xi := int32(1)
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := (2 * dx) - dy

	var (
		counter    int
		screenSync *sync.Mutex = new(sync.Mutex)
	)

	for ; a.Y < b.Y; a.Y++ {
		rendererSync.Add(1)
		go object.drawShapeToWindowAndImage(img, a, screenSync)

		if D > 0 {
			a.X += xi
			D += 2 * -dy
		}
		D += 2 * dx

		if counter++; counter >= 10 {
			rendererSync.Wait()
			counter = 0
		}
	}

	rendererSync.Wait()
	renderer.Present()
}

func (object *SelectedObject) Draw(img *image.RGBA, click visual.Point) {
	if (*object).Shape == nil {
		fmt.Println("Shape is not chosen, unable to draw")
	} else {
		rendererSync.Add(1)
		var screenSync *sync.Mutex = new(sync.Mutex)

		object.drawShapeToWindowAndImage(img, click, screenSync)
		renderer.Present()
	}
}

func (object *SelectedObject) drawShapeToWindowAndImage(img *image.RGBA, click visual.Point, screenSync *sync.Mutex) {
	screenSync.Lock()

	object.prepareAndDraw(click)
	(*object).Shape.DrawToImage(img)

	screenSync.Unlock()

	if (*object).Shape.GetVerticalSize() <= TOOLBAR_Y_POS {
		screenSync.Lock()

		drawToolbar(*object)
		renderer.SetDrawColor(visual.CurrentColor.R, visual.CurrentColor.G, visual.CurrentColor.B, visual.CurrentColor.A)

		screenSync.Unlock()
	}

	rendererSync.Done()
}

func (object *SelectedObject) prepareAndDraw(click visual.Point) {
	(*object).Shape.SetCoordinates(click)
	(*object).Shape.DrawToWindow(renderer)
}

func InitWindowAndToolbar() {
	window, _ := sdl.CreateWindow("Paint", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	renderer, _ = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	initToolHitbox()
	clearWindow()
	drawToolbar(SelectedObject{Shape: nil, LastCheckedColor: -1, LastCheckedShape: -1})
	renderer.Present()
}

func initToolHitbox() {
	for n := 0; n < AMOUNT_TOOLS; n++ {
		hitbox[n] = toolHitbox{upperLeftCorner: visual.Point{X: int32(32/3 + n*32/3 + n*TOOL_FRAME_SIDE), Y: 10}}
		if n < AMOUNT_SHAPES || n == AMOUNT_TOOLS-1 {
			hitbox[n].toolColor = black
		} else {
			switch n {
			case 4:
				hitbox[n].toolColor = red
			case 5:
				hitbox[n].toolColor = orange
			case 6:
				hitbox[n].toolColor = yellow
			case 7:
				hitbox[n].toolColor = green
			case 8:
				hitbox[n].toolColor = cyan
			case 9:
				hitbox[n].toolColor = blue
			case 10:
				hitbox[n].toolColor = purple
			case 11:
				hitbox[n].toolColor = white
			case 12:
				hitbox[n].toolColor = brown
			}
		}
	}
}

func clearWindow() {
	renderer.SetDrawColor(white.R, white.G, white.B, white.A)
	renderer.Clear()
}
