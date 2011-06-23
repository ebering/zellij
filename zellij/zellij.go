package zellij

import "./quadratic/quadratic"
import "os"

var Points map[string]*quadratic.Point
var Tiles []string
var gRs chan int

func init() {
	Points = make(map[string]*quadratic.Point)
	Points["a"], _ = quadratic.PointFromString("0,0,0,2")
	Points["b"], _ = quadratic.PointFromString("-2,0,2,0")
	Points["c"], _ = quadratic.PointFromString("2,-2,2,0")
	Points["d"], _ = quadratic.PointFromString("-2,2,2,0")
	Points["e"], _ = quadratic.PointFromString("2,0,2,0")

	Points["g"], _ = quadratic.PointFromString("-2,0,-2,2")
	Points["h"], _ = quadratic.PointFromString("2,0,-2,2")
	Points["i"], _ = quadratic.PointFromString("-4,0,0,0")
	Points["j"], _ = quadratic.PointFromString("0,-2,0,0")

	Points["l"], _ = quadratic.PointFromString("0,0,0,0")

	Points["n"], _ = quadratic.PointFromString("0,2,0,0")
	Points["o"], _ = quadratic.PointFromString("4,0,0,0")
	Points["p"], _ = quadratic.PointFromString("-2,0,2,-2")

	Points["r"], _ = quadratic.PointFromString("2,0,2,-2")
	Points["s"], _ = quadratic.PointFromString("-2,0,-2,0")
	Points["t"], _ = quadratic.PointFromString("2,-2,-2,0")
	Points["u"], _ = quadratic.PointFromString("-2,2,-2,0")
	Points["v"], _ = quadratic.PointFromString("2,0,-2,0")
	Points["w"], _ = quadratic.PointFromString("0,0,0,-2")



	Points["A"], _ = quadratic.PointFromString("-2,-2,2,0")
	Points["B"], _ = quadratic.PointFromString("2,2,2,0")
	Points["C"], _ = quadratic.PointFromString("-2,-2,-2,0")
	Points["D"], _ = quadratic.PointFromString("2,2,-2,0")

	Tiles = make([]string,6 )
	Tiles[0] = "adehnrvuwtspjgbc"
	Tiles[1] = "beovsi"
	Tiles[2] = "Abgj"
	Tiles[3] = "Cjps"
	Tiles[4] = "Dvrn"
	Tiles[5] = "Bnhe"
}

func TileMap(s string) *quadratic.Map {
	tilePoints := make([]*quadratic.Point, len(s))
	for i, c := range s {
		tilePoints[i] = Points[string(c)].Copy()
	}
	return quadratic.PolygonMap(tilePoints)
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
	gRs = make(chan int,1)
	bc := BoundsChecker(xmin,xmax,ymin,ymax)
	intermediateTilings <- TileMap(Tiles[0])
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
	T.Edges.Do(func (l interface {}) {
		e := l.(*quadratic.Edge)
		q := TileMap(t)
		q.Edges.Do(func (l interface{}) {
			f := l.(*quadratic.Edge)
			if e.IntHeading() == f.IntHeading() && !e.Start().Point.Equal(f.Start().Point) && e.Start().Less(e.End()) {
				//os.Stderr.WriteString("boundary\n")
				Q,ok := T.Overlay(q.Copy().Translate(f.Start(),e.Start()),Overlay)
				if ok == nil && !Q.Isomorphic(T) && boundsCheck(Q) {
					//os.Stderr.WriteString("accept\n")
					sink <- Q
				} 
			}
		})
	})
}

func addTileByVertex(sink chan<- *quadratic.Map, T *quadratic.Map, t string) {
	T.Verticies.Do(func (l interface{}) {
		v := l.(*quadratic.Vertex)
		q := TileMap(t)
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
