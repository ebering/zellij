package zellij

import "../quadratic/quadratic"

func legalVertexFigures(Q *quadratic.Map) bool {
	ret := true
	Q.Verticies.Do(func (l interface{}) {
		v := l.(*quadratic.Vertex)
		var figure byte
		vtx := false
		v.OutgoingEdges.Do(func (l interface{}) {
			e := l.(*quadratic.Edge)
			figure |= byte(1 << uint(e.IntHeading()))
		})
		for _,f := range(VertexFigures) {
			 vtx = vtx || ( figure | f == f)
		}
		ret = ret && vtx
	})
	return ret
}

func leftRotate(b byte, i int) byte {
	j := i % 8
	k := 8-j
	return b << uint(j) | b >> uint(k)
}

var VertexFigures []byte
