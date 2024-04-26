package checkers

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type mouseState int

const (
	mouseStateNone mouseState = iota
	mouseStatePressing
	mouseStateSettled
)

// Input represents the current key states.
type Input struct {
	mouseState    mouseState
	mouseInitPosX int
	mouseInitPosY int
}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}

// Update updates the current input states.
func (i *Input) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		i.mouseState = mouseStatePressing
		i.mouseInitPosX, i.mouseInitPosY = ebiten.CursorPosition()
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		i.mouseState = mouseStateSettled
	} else {
		i.mouseState = mouseStateNone
	}
}
