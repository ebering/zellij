package zellij

import "./quadratic/quadratic"

var Points map[string]*quadratic.Point
var Tiles []string

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

	Tiles = make([]string, 2)
	Tiles[0] = "adehnrvuwtspjgbc"
	Tiles[1] = "beovsi"
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
	intermediateTilings := make(chan *quadratic.Map,1000)
	finalTilings := make(chan *quadratic.Map,100)
	halt := make(chan int)
	intermediateTilings <- TileMap(Tiles[0])
	go func () {
		for {
			T := <-intermediateTilings
			for _,t := range(Tiles) {
				go addTile(intermediateTilings,T,t)
			}
			finalTilings <- T
		}
	}()
	return finalTilings,halt
}

func addTile(sink chan<- *quadratic.Map, T *quadratic.Map, t string) {
	T.Verticies.Do(func (l interface{}) {
		v := l.(*quadratic.Vertex)
		q := TileMap(t)
		q.Verticies.Do(func (l interface{}) {
			u := l.(*quadratic.Vertex)
			if !v.Point.Equal(u.Point) {
				Q,ok := T.Overlay(q.Copy().Translate(u,v))
				if ok == nil && !Q.Isomorphic(T) {
					sink <- Q
				} 
			}
		})
	})
}
