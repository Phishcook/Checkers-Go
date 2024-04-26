package checkers

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	input      *Input
	board      *Board
	boardImage *ebiten.Image
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{
		input: NewInput(),
	}
	var err error
	g.board, err = NewBoard()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) Update() error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(screen.Bounds().Dx()-15, screen.Bounds().Dy()-15)
	}
	g.board.draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 320
}
