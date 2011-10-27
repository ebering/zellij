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

type Colour struct {
	R,G,B,A float64
}

func CreateZellijBrush( skeleton, primary, secondary Colour) func (*quadratic.Face) (float64, float64,float64, float64) {
	secondaryFaces := make(map[*quadratic.Face]bool)
	return func(F *quadratic.Face) (float64, float64, float64, float64) {
		if F.Value.(string) == "outer" {
			return 0.,0.,0.,0.
		}
		if F.Value.(string) == "skeleton" {
			return skeleton.R,skeleton.G,skeleton.B,skeleton.A
		}
		if WhiteTiles[F.Type] {
			return 1.,1.,1.,1.
		}

		neighbors := F.Neighbors()
		nW := 0

		for _,n := range(neighbors) {
			if secondaryFaces[n] {
				return primary.R,primary.G,primary.B,primary.A
			}
			if WhiteTiles[n.Type] {
				nW++
			}
		}

		if nW < len(neighbors)/2 {
			secondaryFaces[F] = true
			return secondary.R,secondary.G,secondary.B,secondary.A
		}

		return primary.R,primary.G,primary.B,primary.A
	}
}
