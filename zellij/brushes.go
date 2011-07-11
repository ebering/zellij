package zellij

import "../quadratic/quadratic"

func PlainBrush(F *quadratic.Face) (float64, float64, float64, float64) {
	return 0., 0., 0., 0.
	/*if F.Type == Tiles[6] || F.Type == Tiles[5] || F.Type == Tiles[4] || F.Type == Tiles[6] {
		return 1.,1.,1.,1.
	} else if F.Type == Tiles[0] || F.Type == Tiles[1] {
		return 0.,28./255.,95./255.,1.
	} else if F.Type == Tiles[2] || F.Type == Tiles[3] || F.Type == Tiles[3] {
		return 31./255.,202./255.,1.,1.
	} */
}
