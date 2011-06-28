package zellij

import "../quadratic/quadratic"

func LegalVertexFigures(Q *quadratic.Map) bool {
	ret := true
	Q.Verticies.Do(func (l interface{}) {
		v := l.(*quadratic.Vertex)
		ret = ret && legalVertexFigure(vertexFigure(v))
	})
	return ret
}

func legalVertexFigure(figure byte) bool {
	vtx := false
	for _,f := range(VertexFigures) {
		 vtx = vtx || ( figure | f == f)
	}
	return vtx
}

func vertexFigure(v *quadratic.Vertex) byte {
	var figure byte
	v.OutgoingEdges.Do(func (l interface{}) {
		e := l.(*quadratic.Edge)
		figure |= byte(1 << uint(e.IntHeading()))
	})
	return figure
}

func leftRotate(b byte, i int) byte {
	j := i % 8
	k := 8-j
	return b << uint(j) | b >> uint(k)
}

var VertexFigures []byte
