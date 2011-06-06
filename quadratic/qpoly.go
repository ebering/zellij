package quadratic

import "os"
import "regexp"

func PolygonMap(verts [] *Point) (* Map) {
	m := NewMap()
	for _,pt := range(verts) {
		v := NewVertex(pt)
		m.verticies.PushFront(v)
	}
	L := m.verticies
	innerFace := new(Face)
	innerFace.Value = "inner"
	outerFace := new(Face)
	outerFace.Value = "outer"
	for l := L.Front(); l != nil  ; l = l.Next() {
		if(l.Next() == nil) {
			f := m.JoinVerticies(l.Value.(*Vertex),L.Front().Value.(*Vertex))
			f.face = innerFace
			innerFace.boundary = f
			f.twin.face = outerFace
			outerFace.boundary = f
			break
		} 
		e := m.JoinVerticies(l.Value.(*Vertex),l.Next().Value.(*Vertex))
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
