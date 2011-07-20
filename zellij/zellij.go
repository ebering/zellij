package zellij

import "../quadratic/quadratic"
import "container/list"
import "os"
import "fmt"

var initializationTime int64

func TileSkeleton(skeleton string, showIntermediate bool) (<-chan *quadratic.Map, chan int) {
	finalTilings := make(chan *quadratic.Map, 10000)
	halt := make(chan int, 1)
	skel, ok := SkeletonMap(skeleton)
	if ok != nil {
		panic("Bad skeleton: " + ok.String() + "\n")
	}
	go tileDriver(skel, finalTilings, halt, chooseNextEdgeByLocation, 0, showIntermediate)
	initializationTime,_,_ = os.Time()
	fmt.Fprintf(os.Stderr,"intitialized at %v\n",initializationTime)
	return finalTilings, halt
}

func TilePlane(maxtiles int, showIntermediate bool) (<-chan *quadratic.Map, chan int) {
	finalTilings := make(chan *quadratic.Map, 10000)
	halt := make(chan int, 1)
	go tileDriver(TileMap("adehnrvuwtspjgbc",0), finalTilings, halt, chooseNextEdgeByLocation, maxtiles, showIntermediate)
	initializationTime,_,_ = os.Time()
	fmt.Fprintf(os.Stderr,"intitialized at %v\n",initializationTime)
	return finalTilings, halt
}

func tileDriver(startTiling * quadratic.Map,sink chan<- *quadratic.Map,halt chan int, chooseNextEdge func(*quadratic.Map) *quadratic.Edge, maxtiles int, showIntermediate bool) {
	alternativeStack := make(chan *list.List,1)
	alternatives := new(list.List)
	alternatives.PushBack(startTiling)
	alternativeStack <- alternatives
	workerCount := make(chan int,1)
	workerCount <- 0
	for {
		select {
			case workers := <-workerCount:
				alternatives = <-alternativeStack
				if alternatives.Len() == 0 {
					if (workers == 0) {
						workerCount <- 0
						halt <- 1
						finishTime, _, _ := os.Time()
						fmt.Fprintf(os.Stderr,"we're done, took %v seconds\n",finishTime-initializationTime)
						return
					}
					alternativeStack <- alternatives
					workerCount <- workers
					continue
				} else if workers < Workers {
					T := alternatives.Remove(alternatives.Front()).(*quadratic.Map)
					alternativeStack <- alternatives
					localMaps := make([]*quadratic.Map, len(TileMaps))
					for j, r := range TileMaps {
						localMaps[j] = r.Copy()
					}
					workerCount <- workers + 1
					go tileWorker(T, alternativeStack,sink,workerCount, halt, localMaps,chooseNextEdge,maxtiles,showIntermediate)
				} else {
					alternativeStack <- alternatives
					workerCount <- workers
				}
			case <-halt:
				halt <- 1
				fmt.Fprintf(os.Stderr,"premature halt\n")
				return
		}
	}
}

func tileWorker(T *quadratic.Map, alternativeStack chan *list.List, sink chan<- *quadratic.Map, workerCount chan int, halt chan int, tileMaps []*quadratic.Map, chooseNextEdge func(*quadratic.Map) *quadratic.Edge, maxtiles int, showIntermediate bool) {
	localAlternatives := new(list.List)
	Work: for {
		select {
			case <-halt:
				halt <- 1
				fmt.Fprintf(os.Stderr,"premature halt\n")
				return
			case L := <-alternativeStack:
				L.PushFrontList(localAlternatives)
				localAlternatives.Init()
				alternativeStack <- L
			default:
				if T.Faces.Len() > maxtiles && maxtiles > 0 {
					sink <- T
					break Work
				} else if noActiveFaces(T) {
					finishTime, _, _ := os.Time()
					fmt.Fprintf(os.Stderr, "new tiling complete, took %v seconds\n",finishTime-initializationTime)
					sink <- T
					break Work
				} else if showIntermediate {
					sink <- T
				}
				alternatives := addTilesByEdge(T, tileMaps, chooseNextEdge)
				if alternatives.Len() == 0 {
					break Work
				}
				T = alternatives.Remove(alternatives.Front()).(*quadratic.Map)
				localAlternatives.PushFrontList(alternatives)
				//fmt.Fprintf(os.Stderr,"currently have %v faces\n",T.Faces.Len())
		}
	}
	L := <-alternativeStack
	L.PushFrontList(localAlternatives)
	localAlternatives.Init()
	alternativeStack <- L

	workers := <-workerCount
	workerCount <- workers - 1
}

func Overlay(f interface{}, g interface{}) (interface{}, os.Error) {
	if f.(string) == "inner" && g.(string) == "inner" {
		return nil, os.NewError("cannot overlap zellij tiles")
	} else if f.(string) == "inner" || g.(string) == "inner" {
		return "inner", nil
	} else if f.(string) == "active" || g.(string) == "active" {
		return "active", nil
	}
	return "outer", nil
}

func addTilesByEdge(T *quadratic.Map, tileMaps []*quadratic.Map, chooseNextEdge func(*quadratic.Map) *quadratic.Edge) *list.List {
	e := chooseNextEdge(T)
	onGeneration := e.Generation
	sink := new(list.List)
	for _, t := range tileMaps {
		q := t.Copy()
		q.SetGeneration(onGeneration + 1)

		q.Edges.Do(func(l interface{}) {
			f := l.(*quadratic.Edge)
			if e.IntHeading() == f.IntHeading() &&
				e.LengthSquared().Equal(f.LengthSquared()) &&
				legalVertexFigure(vertexFigure(e.Start())|vertexFigure(f.Start())) &&
				legalVertexFigure(vertexFigure(e.End())|vertexFigure(f.End())) {
				toLay := GenerateOrbit(q.Copy().Translate(f.Start(), e.Start()),"d4")
				Q := T.Copy()
				goodTiling := true
				for _,o := range(toLay) {
					var ok os.Error
					Q, ok = Q.Overlay(o, Overlay)
					if ok == nil && Q != nil && LegalVertexFigures(Q) {
						continue
					}
					goodTiling = false
					break
				}
				if goodTiling {
					sink.PushBack(Q)
				}
			}
		})
	}
	//fmt.Fprintf(os.Stderr,"on generation %v found %v children\n",onGeneration,sink.Len())

	return sink
}

func noActiveFaces(Q *quadratic.Map) bool {
	ret := true
	Q.Faces.Do(func(f interface{}) {
		ret = ret && f.(*quadratic.Face).Value.(string) != "active"
	})
	return ret
}


func duplicateTiling(tilings *list.List, T *quadratic.Map) bool {
	for l := tilings.Front(); l != nil; l = l.Next() {
		if T.Isomorphic(l.Value.(*quadratic.Map)) {
			return true
		}
	}
	return false
}

func removeDuplicates(bigList *list.List, littleList *list.List) {
	for l := bigList.Front(); l != nil; l = l.Next() {
		for m := littleList.Front(); m != nil; m = m.Next() {
			if l.Value.(*quadratic.Map).Isomorphic(m.Value.(*quadratic.Map)) {
				littleList.Remove(m)
			}
		}
	}
}
