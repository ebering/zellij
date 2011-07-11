package zellij

import "container/vector"
import "../quadratic/quadratic"
import "sort"

type generationOrderedEdges struct {
	vector.Vector
}

func (G *generationOrderedEdges) Less(i, j int) bool {
	eI := G.At(i).(*quadratic.Edge)
	eJ := G.At(j).(*quadratic.Edge)
	if eI.Start().Equal(eJ.Start().Point) {
		return eI.IntHeading() < eJ.IntHeading()
	}
	return eI.Start().Less(eJ.Start())
}

func EdgeStillActive(T *quadratic.Map, e *quadratic.Edge) bool {
	for i := 0; i < T.Edges.Len(); i++ {
		f := T.Edges.At(i).(*quadratic.Edge)
		if f.Start().Equal(e.Start().Point) && f.End().Equal(e.End().Point) {
			return f.Face().Value.(string) == "active"
		}
	}
	return false
}

func chooseNextEdgeByGeneration(T *quadratic.Map) *quadratic.Edge {
	ActiveFaces := new(vector.Vector)
	onGeneration := -1
	T.Faces.Do(func(F interface{}) {
		Fac := F.(*quadratic.Face)
		if Fac.Value.(string) != "active" {
			return
		}
		ActiveFaces.Push(Fac)
		Fac.DoEdges(func(e *quadratic.Edge) {
			//fmt.Fprintf(os.Stderr,"edge generation: %v onGen: %v\n",e.Generation,onGeneration)
			if onGeneration < 0 || e.Generation < onGeneration {
				onGeneration = e.Generation
			}
		})
	})
	activeEdges := new(generationOrderedEdges)

	ActiveFaces.Do(func(F interface{}) {
		Fac := F.(*quadratic.Face)
		Fac.DoEdges(func(e (*quadratic.Edge)) {
			if e.Generation != onGeneration {
				return
			}
			activeEdges.Push(e)
		})
	})
	//fmt.Fprintf(os.Stderr,"onGen: %v have %v active edges\n",onGeneration,activeEdges.Len())
	sort.Sort(activeEdges)
	return activeEdges.At(0).(*quadratic.Edge)
}

func chooseNextEdgeByLocation(T *quadratic.Map) *quadratic.Edge {
	activeEdges := new(generationOrderedEdges)

	T.Faces.Do(func(F interface{}) {
		Fac := F.(*quadratic.Face)
		if Fac.Value.(string) != "active" {
			return
		}
		Fac.DoEdges(func(e *quadratic.Edge) {
			activeEdges.Push(e)
		})
	})

	sort.Sort(activeEdges)
	return activeEdges.At(0).(*quadratic.Edge)
}
