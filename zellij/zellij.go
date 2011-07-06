package zellij

import "../quadratic/quadratic"
import "container/list"
import "os"
//import "fmt"

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

func TileSkeleton(skeleton string,showIntermediate bool) (<-chan *quadratic.Map, chan<- int) {
	intermediateTilings := make(chan *list.List,1)
	finalTilings := make(chan *quadratic.Map,1000)
	halt := make(chan int,Workers)
	L := list.New()
	skel,ok := SkeletonMap(skeleton)
	if ok != nil {
		panic("Bad skeleton: "+ok.String()+"\n")
	}
	L.PushBack(skel)
	intermediateTilings <- L
	for i:= 0; i < Workers; i++ {
		localMaps := make([]*quadratic.Map,len(TileMaps))
		for j,r := range(TileMaps) {
			localMaps[j] = r.Copy()
		}
		go tileWorker(intermediateTilings,finalTilings,halt,localMaps,chooseNextEdgeByLocation,0,showIntermediate)
	}
	return finalTilings,halt
}

func TilePlane(maxtiles int,showIntermediate bool) (<-chan *quadratic.Map, chan<- int) {
	//center := quadratic.NewPoint(xmax.Sub(xmin),ymax.Sub(ymin))
	intermediateTilings := make(chan *list.List,1)
	finalTilings := make(chan *quadratic.Map,1000)
	halt := make(chan int,Workers)
	L := list.New()
	L.PushBack(TileMap(Tiles[0],0))
	L.Front().Value.(*quadratic.Map).Faces.Do(func (f interface{}) {
		F := f.(*quadratic.Face)
		if F.Value.(string) == "outer" {
			F.Value = "active"
		}
	})
	intermediateTilings <- L
	for i:= 0; i < Workers; i++ {
		localMaps := make([]*quadratic.Map,len(TileMaps))
		for j,r := range(TileMaps) {
			localMaps[j] = r.Copy()
		}
		go tileWorker(intermediateTilings,finalTilings,halt,localMaps,chooseNextEdgeByGeneration,maxtiles,showIntermediate)
	}
	return finalTilings,halt
}
	
func tileWorker (source chan *list.List, sink chan<- *quadratic.Map, halt chan int, tileMaps []*quadratic.Map,chooseNextEdge func (*quadratic.Map) (*quadratic.Edge), maxtiles int,showIntermediate bool) {
	for {
		select {
			case L := <-source:
				if L.Len() == 0 { source <- L; continue }
				T := L.Remove(L.Front()).(*quadratic.Map)
				source <- L
				if T.Faces.Len() > maxtiles && maxtiles > 0 {
					sink <- T
					continue
				} else if noActiveFaces(T) {
					sink <- T
					continue
				} else if showIntermediate {
					sink <- T
				}
				//fmt.Fprintf(os.Stderr,"currently have %v faces\n",T.Faces.Len())
				localSink := addTilesByEdge(T,tileMaps,chooseNextEdge)
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
	} else if f.(string) == "active" || g.(string) == "active" {
		return "active",nil
	}
	return "outer",nil
}

func addTilesByEdge(T *quadratic.Map,tileMaps []*quadratic.Map, chooseNextEdge func (*quadratic.Map) (*quadratic.Edge)) (* list.List) {
	e := chooseNextEdge(T)
	onGeneration := e.Generation
	sink := new(list.List)
	for _,t := range(tileMaps) {
		q := t.Copy()
		q.SetGeneration(onGeneration+1)

		q.Edges.Do(func (l interface{}) {
			f := l.(*quadratic.Edge)
			if e.IntHeading() == f.IntHeading()  && 
				e.LengthSquared().Equal(f.LengthSquared())  &&
				legalVertexFigure(vertexFigure(e.Start()) | vertexFigure(f.Start())) &&
				legalVertexFigure(vertexFigure(e.End()) | vertexFigure(f.End())) { 
				Q,ok := T.Overlay(q.Copy().Translate(f.Start(),e.Start()),Overlay)
				if ok == nil && Q != nil &&  LegalVertexFigures(Q)  {
					//fmt.Fprintf(os.Stderr,"adding %v to %v dup checks %v %v\n",f,e,duplicateTiling(sink,Q),duplicateTiling(oldSink,Q))
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
	Q.Faces.Do(func (f interface{}) {
		ret = ret && f.(*quadratic.Face).Value.(string) != "active"
	})
	return ret
}
 

func duplicateTiling(tilings *list.List,T *quadratic.Map) bool {
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
