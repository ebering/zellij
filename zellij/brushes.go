package zellij

import "../quadratic/quadratic"

func PlainBrush(F *quadratic.Face) (float64, float64, float64, float64) {
	for i,t := range(Tiles) {
		if F.Type == t {
			return float64(i)/float64(len(Tiles)),0.,0.,1.
		}
	}
	return 0.,0.,0.,0.
}
