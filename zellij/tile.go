package zellij

import "../quadratic/quadratic"

func TileMap(s string, Generation int) *quadratic.Map {
	tilePoints := make([]*quadratic.Point, len(s))
	for i, c := range s {
		tilePoints[i] = Points[string(c)].Copy()
	}
	tile := quadratic.PolygonMap(tilePoints)
	tile.Edges.Do(func(f interface{}) {
		f.(*quadratic.Edge).Generation = Generation
	})
	tile.Faces.Do(func(f interface{}) {
		if f.(*quadratic.Face).Value.(string) == "inner" {
			f.(*quadratic.Face).Type = s
		}
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
