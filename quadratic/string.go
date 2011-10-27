package quadratic

import "fmt"

func (i *Integer) String() string {
	return fmt.Sprintf("%v+%v√2", i.a, i.b)
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.x, p.y)
}

func (e *Edge) String() string {
	return fmt.Sprintf("Edge: ( %v , %v ) %p", e.start, e.end,e.face)
}
