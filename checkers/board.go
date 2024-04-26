package checkers

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Board struct {
	squares        [8][8]*Square
	selectedSquare *Square
	clickedX       *int
	clickedY       *int
	opoonentTurn   bool
}

func NewBoard() (*Board, error) {
	boardGrid := [8][8]*Square{
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil}}

	for j := 0; j < 8; j++ {
		for i := 0; i < 8; i++ {
			var checker *Checker
			if j%2 != i%2 && j != 3 && j != 4 {
				checkerG := Checker{j > 4, false}
				checker = &checkerG
			}

			boardGrid[i][j] = &Square{i, j, j%2 != i%2, checker, false, false}
		}
	}
	return &Board{boardGrid, nil, nil, nil, false}, nil
}

func (b *Board) Update(input *Input) error {
	for _, r := range b.squares {
		for _, s := range r {
			if err := s.Update(); err != nil {
				return err
			}
		}
	}

	if input.mouseState == mouseStatePressing {
		x, y := ebiten.CursorPosition()
		b.clickedX, b.clickedY = &x, &y

	}

	return nil
}

func (b *Board) doOpponentMove() {
	moveablePieces := b.getMoveablePieces(false)
	if len(moveablePieces) == 0 {
		fmt.Println("No moves remain")
		b.opoonentTurn = false
		return
	}
	squareToMove := moveablePieces[rand.Intn(len(moveablePieces))]

	squareOptions := b.getMoveableSquaresFromSquare(squareToMove)
	selectedDestination := squareOptions[rand.Intn(len(squareOptions))]

	if jumpedSquare := b.getSquareBetween(squareToMove, selectedDestination); jumpedSquare != nil {
		jumpedSquare.checker = nil
	}
	selectedDestination.checker = squareToMove.checker
	if selectedDestination.Y == 7 {
		selectedDestination.checker.king = true
	}
	squareToMove.checker = nil
	b.opoonentTurn = false
}

func (b *Board) getMoveableSquaresFromSquare(s *Square) []*Square {
	squaresToMoveTo := []*Square{}
	adjacentSquares := b.getAdjacentSquares(s)
	for k := range adjacentSquares {
		if s.checkerCanMoveTo(adjacentSquares[k]) {
			squaresToMoveTo = append(squaresToMoveTo, adjacentSquares[k])
		} else if jumpSquare := b.getCheckerJumpSquare(s, adjacentSquares[k]); jumpSquare != nil {
			squaresToMoveTo = append(squaresToMoveTo, jumpSquare)
		}
	}

	return squaresToMoveTo
}

func (b *Board) getMoveablePieces(player bool) []*Square {
	moveableSquares := []*Square{}
	for j := range b.squares {
		for i := range b.squares[j] {
			s := b.squares[j][i]
			if s.checker == nil || s.checker.player != player {
				continue
			}
			if len(b.getMoveableSquaresFromSquare(s)) > 0 {
				moveableSquares = append(moveableSquares, s)
				break
			}
		}
	}

	return moveableSquares
}

func (b *Board) getAdjacentSquares(s *Square) []*Square {
	adjacentSquares := []*Square{}
	yDir := -1
	if !s.checker.player {
		yDir = 1
	}
	yTarget := s.Y + yDir
	if 0 <= yTarget && yTarget <= 7 {
		if s.X > 0 {
			adjacentSquares = append(adjacentSquares, b.squares[s.X-1][yTarget])
		}
		if s.X < 7 {
			adjacentSquares = append(adjacentSquares, b.squares[s.X+1][yTarget])
		}
	}

	if s.checker.king {
		yTarget = s.Y - yDir
		if 0 <= yTarget && yTarget <= 7 {
			if s.X > 0 {
				adjacentSquares = append(adjacentSquares, b.squares[s.X-1][yTarget])
			}
			if s.X < 7 {
				adjacentSquares = append(adjacentSquares, b.squares[s.X+1][yTarget])
			}
		}
	}
	return adjacentSquares
}

func (b *Board) highlightAdjacentSquares(s *Square) {
	b.deHighlightSquares()
	adjacentSquares := b.getAdjacentSquares(s)

	for i := range adjacentSquares {
		bs := adjacentSquares[i]
		if s.checkerCanMoveTo(bs) {
			bs.canMoveTo = true
		} else if jumpSquare := b.getCheckerJumpSquare(s, bs); jumpSquare != nil {
			jumpSquare.canMoveTo = true
		}
	}
}

func (b *Board) deHighlightSquares() {
	for j := range b.squares {
		for i := range b.squares[j] {
			b.squares[i][j].canMoveTo = false
		}
	}
}

func (b *Board) getSquareBetween(s *Square, ts *Square) *Square {
	x := (s.X + ts.X)
	y := (s.Y + ts.Y)
	if x%2 != 0 || y%2 != 0 {
		return nil
	}
	return b.squares[x/2][y/2]
}

func (b *Board) getCheckerJumpSquare(s *Square, ts *Square) *Square {
	if ts.checker == nil || ts.checker.player == s.checker.player {
		return nil
	}
	lX := ts.X - (s.X - ts.X)
	lY := ts.Y - (s.Y - ts.Y)
	if lX > 7 || lX < 0 || lY > 7 || lY < 0 {
		return nil
	}

	if b.squares[lX][lY].checker == nil {
		return b.squares[lX][lY]
	}
	return nil
}

func pointIsInSquare(x1 int, y1 int, x2 int, y2 int, x int, y int) bool {
	return (x > x1 && x < x2 && y > y1 && y < y2)
}

func (b *Board) draw(screen *ebiten.Image) {

	// board
	vector.DrawFilledRect(screen, float32(screen.Bounds().Min.X)+10.0, float32(screen.Bounds().Min.Y)+10.0, float32(screen.Bounds().Max.X)-20.0, float32(screen.Bounds().Max.Y)-20.0, color.RGBA{160, 82, 45, 255}, true)

	// startX/Y represent where the Square grid starts
	startX := float32(screen.Bounds().Min.X) + 15.0
	startY := float32(screen.Bounds().Min.Y) + 15.0
	squareWidth := (float32(screen.Bounds().Max.X) - 30.0) / 8

	for j := range b.squares {
		for i := range b.squares[j] {
			s := b.squares[j][i]
			var c color.Color
			if s.moveable {
				c = color.RGBA{255, 0, 0, 255}
			} else {
				c = color.Black
			}

			if s.selected {
				c = color.RGBA{0, 230, 0, 255}
			}

			if s.canMoveTo {
				c = color.RGBA{250, 250, 205, 255}
			}

			lowerX := startX + (float32(s.X) * squareWidth)
			lowerY := startY + (float32(s.Y) * squareWidth)
			vector.DrawFilledRect(screen, lowerX, lowerY, squareWidth, squareWidth, c, true)

			if s.checker != nil {
				var checkerColor color.Color
				if s.checker.player {
					checkerColor = color.White
				} else {
					checkerColor = color.RGBA{255, 255, 0, 255}
				}
				vector.DrawFilledCircle(screen, lowerX+(squareWidth/2), lowerY+(squareWidth/2), squareWidth/3, checkerColor, true)

				if s.checker.king {
					vector.DrawFilledCircle(screen, lowerX+(squareWidth/2), lowerY+(squareWidth/2), squareWidth/8, color.RGBA{30, 215, 40, 255}, true)
				}
			}
		}
	}

	if b.opoonentTurn {
		b.doOpponentMove()
	}

	for j := range b.squares {
		for i := range b.squares[j] {
			s := b.squares[j][i]
			lowerX := startX + (float32(s.X) * squareWidth)
			lowerY := startY + (float32(s.Y) * squareWidth)
			if b.clickedX != nil && b.clickedY != nil {
				if pointIsInSquare(int(lowerX), int(lowerY), int(lowerX+squareWidth), int(lowerY+squareWidth), *b.clickedX-8, *b.clickedY-8) {
					if s.checker != nil && s.checker.player {
						s.selected = true
						b.selectedSquare = s
						b.highlightAdjacentSquares(s)
					} else if s.canMoveTo {
						if jumpedSquare := b.getSquareBetween(s, b.selectedSquare); jumpedSquare != nil {
							jumpedSquare.checker = nil
						}
						s.checker = b.selectedSquare.checker
						b.selectedSquare.checker = nil
						b.selectedSquare = nil
						b.deHighlightSquares()
						b.opoonentTurn = true
						if s.Y == 0 {
							s.checker.king = true
						}
					} else {
						b.deHighlightSquares()
					}
				} else {
					s.selected = false
				}
			}
		}
	}
	b.clickedX, b.clickedY = nil, nil
}
