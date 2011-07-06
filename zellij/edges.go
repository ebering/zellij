package zellij

import "container/vector"
import "../quadratic/quadratic"

type generationOrderedEdges struct {
	vector.Vector
}

func (G *generationOrderedEdges) Less(i,j int) bool {
	eI := G.At(i).(*quadratic.Edge)
	eJ := G.At(j).(*quadratic.Edge)
	if eI.Start().Equal(eJ.Start().Point) {
		return eI.IntHeading() < eJ.IntHeading()
	}
	return eI.Start().Less(eJ.Start())
}

func EdgeStillActive(T *quadratic.Map, e *quadratic.Edge) bool {
	for i := 0; i < T.Edges.Len(); i ++ {
		f := T.Edges.At(i).(*quadratic.Edge)
		if f.Start().Equal(e.Start().Point) && f.End().Equal(e.End().Point) {
			return f.Face().Value.(string) == "active"
		}
	}
	return false
}
