package zellij

import "../quadratic/quadratic"

func LegalVertexFigures(Q *quadratic.Map) bool {
	ret := false
	Q.Verticies.Do(func (l interface{}) {
		v := l.(*quadratic.Vertex)
		var figure int
		v.OutgoingEdges.Do(func (l interface{}) {
			e := l.(*quadratic.Vertex)
			figure |= 1 << e.IntHeading()
		}
		for _,f := range(VertexFigures) {
			ret = ret || ( figure | f == f)
		}
	})
	return ret
}

func LeftRotate(b byte, i int) byte {
	j := i % 8
	k := 8-j
	return b << uint(j) | b >> uint(k)
}

VertexFigures :=  []byte{
LeftRotate(5,0),
LeftRotate(5,1),
LeftRotate(5,2),
LeftRotate(5,3),
LeftRotate(5,4),
LeftRotate(5,5),
LeftRotate(5,6),
LeftRotate(5,7),
LeftRotate(9,0),
LeftRotate(9,1),
LeftRotate(9,2),
LeftRotate(9,3),
LeftRotate(9,4),
LeftRotate(9,5),
LeftRotate(9,6),
LeftRotate(9,7),
LeftRotate(85,0),
LeftRotate(85,1),
LeftRotate(51,0),
LeftRotate(51,1),
LeftRotate(51,2),
LeftRotate(51,3)}
