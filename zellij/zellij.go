package zellij

import "../quadratic/quadratic"
import "container/vector"
import "container/list"
import "os"
import "fmt"

func TileMap(s string,Generation int) *quadratic.Map {
	tilePoints := make([]*quadratic.Point, len(s))
	for i, c := range s {
		tilePoints[i] = Points[string(c)].Copy()
	}
	tile := quadratic.PolygonMap(tilePoints)
	tile.Edges.Do(func (f interface{}) {
		f.(*quadratic.Edge).Generation = Generation
	})
	return tile
}

func PathMap(s string) *quadratic.Map {
	tilePoints := make([]*quadratic.Point, len(s))
	for i, c := range s {
		tilePoints[i] = Points[string(c)].Copy()
	}
	return quadratic.PathMap(tilePoints)
}

func TilePlane() (<-chan *quadratic.Map, chan<- int) {
	//center := quadratic.NewPoint(xmax.Sub(xmin),ymax.Sub(ymin))
	intermediateTilings := make(chan *list.List,1)
	finalTilings := make(chan *quadratic.Map,1000)
	halt := make(chan int,Workers)
	L := list.New()
	L.PushBack(TileMap(Tiles[0],0))
	intermediateTilings <- L
	for i:= 0; i < Workers; i++ {
		go tileWorker(intermediateTilings,finalTilings,halt)
	}
	return finalTilings,halt
}
	
func tileWorker (source chan *list.List, sink chan<- *quadratic.Map, halt chan int) {
	localSink := list.New()
	for {
		select {
			case L := <-source:
				if L.Len() == 0 { source <- L; continue }
				T := L.Remove(L.Front()).(*quadratic.Map)
				source <- L
				if T.Faces.Len() > 10 {
					sink <- T
					continue
				}
				fmt.Fprintf(os.Stderr,"currently have %v faces\n",T.Faces.Len())
				localSink.Init()
				addTilesByEdge(localSink,T)
				L = <-source
				L.PushFrontList(localSink)
				source <- L
			case <-halt:
				halt <- 1
				return
		}
	}
}

func Overlay(f interface{}, g interface{}) (interface{},os.Error) {
	if f.(string) == "inner" && g.(string) == "inner" {
		return nil,os.NewError("cannot overlap zellij tiles")
	} else if f.(string) == "inner" || g.(string) == "inner" {
		return "inner",nil
	}
	return "outer",nil
}

func addTilesByEdge(sink *list.List, T *quadratic.Map ) {
	ActiveFaces := new(vector.Vector)
	onGeneration := -1
	T.Faces.Do(func (F interface{}) {
		Fac := F.(*quadratic.Face)
		if(Fac.Value.(string) != "outer") {
			return
		}
		ActiveFaces.Push(Fac)
		Fac.DoEdges(func (e *quadratic.Edge) {
			//fmt.Fprintf(os.Stderr,"edge generation: %v onGen: %v\n",e.Generation,onGeneration)
			if onGeneration < 0 || e.Generation < onGeneration {
				onGeneration = e.Generation
			}
		})
	})
	//fmt.Fprintf(os.Stderr,"onGen: %v\n",onGeneration)
	for _,t := range(TileMaps) {
		q := t.Copy()
		q.SetGeneration(onGeneration+1)
		ActiveFaces.Do(func (F interface{}) {
			Fac := F.(*quadratic.Face)
			Fac.DoEdges(func (e (*quadratic.Edge)) {
				if (e.Generation != onGeneration) {
					return
				}
				q.Edges.Do(func (l interface{}) {
					f := l.(*quadratic.Edge)
					if e.IntHeading() == f.IntHeading()  {
						Q,ok := T.Overlay(q.Copy().Translate(f.Start(),e.Start()),Overlay)
						if ok == nil && !Q.Isomorphic(T) && legalVertexFigures(Q) && !duplicateTiling(sink,Q) {
							sink.PushBack(Q)
						}
					}
				})
			})
		})
	}
}

func duplicateTiling(tilings *list.List,T *quadratic.Map) bool {
	for l := tilings.Front(); l != nil; l = l.Next() {
		if T.Isomorphic(l.Value.(*quadratic.Map)) {
			return true
		}
	}
	return false
}
