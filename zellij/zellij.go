package zellij

import "../quadratic/quadratic"
import "sort"
import "os"

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

func TileRegion(xmin,xmax,ymin,ymax *quadratic.Integer) (<-chan *quadratic.Map, chan<- int) {
	//center := quadratic.NewPoint(xmax.Sub(xmin),ymax.Sub(ymin))
	intermediateTilings := make(chan *quadratic.Map,100)
	finalTilings := make(chan *quadratic.Map)
	halt := make(chan int)
	bc := BoundsChecker(xmin,xmax,ymin,ymax)
	intermediateTilings <- TileMap(Tiles[0],0)
	go func () {
		for {
			T := <-intermediateTilings
			for _,t := range(Tiles) {
				go addTileByEdge(intermediateTilings,finalTilings,bc,T,t)
			} 
			finalTilings <- T
		}
	}()
	return finalTilings,halt
}

func BoundsChecker(xmin,xmax,ymin,ymax *quadratic.Integer) (func (*quadratic.Map) bool) {
	return func (m *quadratic.Map) bool {
		ret := true
		m.Verticies.Do(func (l interface{}) {
			v := l.(*quadratic.Vertex)
			ret = ret && xmin.Less(v.X()) && v.X().Less(xmax) && ymin.Less(v.Y()) && v.Y().Less(ymax)
		})
		return ret
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

func addTileByEdge(sink chan<- *quadratic.Map, finalSink chan<- *quadratic.Map, boundsCheck func(*quadratic.Map) bool, T *quadratic.Map, t string) {
	sort.Sort(GenerationalEdges{T.Edges})
	onGeneration := T.Edges.At(0).(*quadratic.Edge).Generation
	q := TileMap(t,onGeneration+1)
	T.Faces.Do(func (F interface{}) {
		Fac := F.(*quadratic.Face)
		if(Fac.Value.(string) != "outer") {
			return
		}
		Fac.DoEdges(func (e (*quadratic.Edge)) {
			if (e.Generation != onGeneration) {
				return
			}
			q.Edges.Do(func (l interface{}) {
				f := l.(*quadratic.Edge)
				if e.IntHeading() == f.IntHeading() && !e.Start().Point.Equal(f.Start().Point) && e.Start().Less(e.End()) {
					Q,ok := T.Overlay(q.Copy().Translate(f.Start(),e.Start()),Overlay)
					if ok == nil && !Q.Isomorphic(T) && LegalVertexFigures(Q) {
						sink <- Q
					} 
				}
			})
		})
	})
}

func addTileByVertex(sink chan<- *quadratic.Map, T *quadratic.Map, t string) {
	T.Verticies.Do(func (l interface{}) {
		v := l.(*quadratic.Vertex)
		q := TileMap(t,0)
		q.Verticies.Do(func (l interface{}) {
			u := l.(*quadratic.Vertex)
			if !v.Point.Equal(u.Point) {
				Q,ok := T.Overlay(q.Copy().Translate(u,v),Overlay)
				if ok == nil && !Q.Isomorphic(T) {
					sink <- Q
				} 
			}
		})
	})
}
