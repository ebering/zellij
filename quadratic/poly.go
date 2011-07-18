package quadratic

import "os"
import "regexp"

// Verticies should be the boundary word of the Outer face.
func PolygonMap(verts []*Point) *Map {
	m := NewMap()
	for _, pt := range verts {
		v := NewVertex(pt)
		m.Verticies.Push(v)
	}
	innerFace := new(Face)
	innerFace.fromMap = m
	innerFace.Value = "inner"
	outerFace := new(Face)
	outerFace.fromMap = m
	outerFace.Value = "outer"
	for i := 0; i < m.Verticies.Len(); i++ {
		if i == m.Verticies.Len()-1 {
			f := m.JoinVerticies(m.Verticies.Last().(*Vertex), m.Verticies.At(0).(*Vertex))
			f.face = outerFace
			outerFace.boundary = f
			f.twin.face = innerFace
			innerFace.boundary = f.twin
			break
		}
		e := m.JoinVerticies(m.Verticies.At(i).(*Vertex), m.Verticies.At(i+1).(*Vertex))
		e.face = outerFace
		outerFace.boundary = e
		e.twin.face = innerFace
		innerFace.boundary = e.twin
	}
	m.Faces.Push(innerFace)
	m.Faces.Push(outerFace)
	m.Init()
	return m
}

func PathMap(verts []*Point) *Map {
	m := NewMap()
	for _, pt := range verts {
		v := NewVertex(pt)
		m.Verticies.Push(v)
	}
	innerFace := new(Face)
	innerFace.fromMap = m
	innerFace.Value = "inner"
	outerFace := new(Face)
	outerFace.fromMap = m
	outerFace.Value = "outer"
	for i := 0; i < m.Verticies.Len()-1; i++ {
		e := m.JoinVerticies(m.Verticies.At(i).(*Vertex), m.Verticies.At(i+1).(*Vertex))
		e.face = innerFace
		innerFace.boundary = e
		e.twin.face = outerFace
		outerFace.boundary = e
	}
	m.Init()
	return m
}

func PolygonMapFromString(str string) (*Map, os.Error) {
	re := regexp.MustCompile("([0-9]+,[0-9]+,[0-9]+,[0-9]+)")
	matches := re.FindAllString(str, -1)
	points := make([]*Point, len(matches))
	for i, m := range matches {
		var ok os.Error
		points[i], ok = PointFromString(m)
		if ok != nil {
			return nil, ok
		}
	}
	return PolygonMap(points), nil
}
