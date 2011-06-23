package zellij

import "container/vector"
import "../quadratic/quadratic"

type GenerationalEdges struct { *vector.Vector }

func (ge GenerationalEdges) Less(i,j int) bool {
	return ge.Vector.At(i).(*quadratic.Edge).Generation < ge.Vector.At(j).(*quadratic.Edge).Generation 
}

func (ge GenerationalEdges) Swap(i,j int) {
	ge.Vector.Swap(i,j)
}

func (ge GenerationalEdges) Len() int {
	return ge.Vector.Len()
}
