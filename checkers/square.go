package checkers

type Square struct {
	X, Y      int
	moveable  bool
	checker   *Checker
	selected  bool
	canMoveTo bool
}

//	func (s *Square) getBounts() (x1, y1, x2, y2 int) {
//		// https://github.com/hajimehoshi/ebiten/blob/main/examples/2048/2048/game.go
//		// Input, Board need their own state
//		// Input will be stored to g Game and in Game.Update we call each in order
//	}
func (s *Square) Update() error {

	return nil
}

func (s *Square) checkerCanMoveTo(ts *Square) bool {
	xDiff := ts.X - s.X
	if xDiff > 0 {
		xDiff = -xDiff
	}
	yDiff := ts.Y - s.Y
	if s.checker.king {

		if yDiff > 0 {
			yDiff = -yDiff
		}
	}
	movDir := -1
	if !s.checker.player {
		yDiff = -yDiff
	}

	if xDiff == movDir && yDiff == movDir {
		if ts.checker == nil {
			return true
		}
	}
	return false

}
