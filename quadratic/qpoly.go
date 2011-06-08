package quadratic

import "os"
import "regexp"

func PolygonMap(verts [] *Point) (* Map) {
	m := NewMap()
	for _,pt := range(verts) {
		v := NewVertex(pt)
		m.Verticies.Push(v)
	}
	innerFace := new(Face)
	innerFace.Value = "inner"
	outerFace := new(Face)
	outerFace.Value = "outer"
	for i := 0; i < m.Verticies.Len(); i++ {
		if(i == m.Verticies.Len()-1) {
			f := m.JoinVerticies(m.Verticies.Last().(*Vertex),m.Verticies.At(0).(*Vertex))
			f.face = innerFace
			innerFace.boundary = f
			f.twin.face = outerFace
			outerFace.boundary = f
			break
		} 
		e := m.JoinVerticies(m.Verticies.At(i).(*Vertex),m.Verticies.At(i+1).(*Vertex))
		e.face = innerFace
		innerFace.boundary = e
		e.twin.face = outerFace
		outerFace.boundary = e
	}
	return m
}

func PathMap(verts [] *Point) (* Map) {
	m := NewMap()
	for _,pt := range(verts) {
		v := NewVertex(pt)
		m.Verticies.Push(v)
	}
	innerFace := new(Face)
	innerFace.Value = "inner"
	outerFace := new(Face)
	outerFace.Value = "outer"
	for i := 0; i < m.Verticies.Len()-1; i++ {
		e := m.JoinVerticies(m.Verticies.At(i).(*Vertex),m.Verticies.At(i+1).(*Vertex))
		e.face = innerFace
		innerFace.boundary = e
		e.twin.face = outerFace
		outerFace.boundary = e
	}
	return m
}

func PolygonMapFromString(str string) (* Map,os.Error) {
	re := regexp.MustCompile("([0-9]+,[0-9]+,[0-9]+,[0-9]+)")
	matches := re.FindAllString(str,-1)
	points := make([]*Point,len(matches))
	for i,m := range(matches) {
		var ok os.Error
		points[i],ok = PointFromString(m)
		if ok != nil {
			return nil,ok
		}
	}
	return PolygonMap(points),nil
}
