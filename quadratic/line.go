package quadratic

import "math"
//import "fmt"
//import "os"

type Line struct {
	start,end *Point
	deltax,deltay *Integer
}

func NewLine(p1,p2 *Point) (*Line) {
	l := new(Line)

	if p1.Less(p2) {
		l.start,l.end = p1,p2
	} else {
		l.start,l.end = p2,p1
	} 

	l.deltax = l.end.x.Sub(l.start.x)
	l.deltay = l.end.y.Sub(l.start.y)
	
	return l
}

func (l *Line) Vertical() bool {
	return l.deltax.Equal(&Integer{0,0})
}

func (l *Line) Horizontal() bool {
	return l.deltay.Equal(&Integer{0,0})
}

func (l *Line) On(p *Point) (bool) {
	lhs := l.deltax.Mul(p.y.Sub(l.start.y))
	rhs := l.deltay.Mul(p.x.Sub(l.start.x))
	
	return  lhs.Equal(rhs) && ( l.start.x.Less(p.x) && p.x.Less(l.end.x) || 
					l.deltax.Equal(&Integer{0,0}) && l.start.y.Less(p.y) && p.y.Less(l.end.y) ||
					l.start.Equal(p) || 
					l.end.Equal(p))
}

func (l *Line) FloatOn(x,y float64) (bool) {
	lhs := l.deltax.Float64()*(y - l.start.y.Float64())
	rhs := l.deltay.Float64()*(x - l.start.x.Float64())
	
	return math.Fabs(lhs-rhs) < FLOAT64_EPSILON  && ( l.start.x.Float64() < x && x < l.end.x.Float64() || l.deltax.Equal(&Integer{0,0}) && l.start.y.Float64() < y && y < l.end.y.Float64())
}

// determines if p is below l
func (l *Line) Below(p *Point) (bool) {
	lhs := l.deltax.Mul(p.y.Sub(l.start.y))
	rhs := l.deltay.Mul(p.x.Sub(l.start.x))

	return lhs.Less(rhs)
}

// Determines if l < m in the plane sweep line ordering, assuming the sweep is currently at p
func (l *Line) LessAt(m *Line,p *Point) bool {
	// First special case, l,m isect at p, in which case go by slope
	//fmt.Fprintf(os.Stderr,"Testing: %v and %v",l,m)
	if l.On(p) && m.On(p) {
		//fmt.Fprintf(os.Stderr," by slope: %v\n",l.end.y.Less(m.end.y))
		return l.end.y.Less(m.end.y)
	}

	// Cases for vertical lines:
	// l,m vertical, compare start points:
	if l.Vertical() && m.Vertical() {
		//fmt.Fprintf(os.Stderr," by height (both vert): %v\n",l.end.y.Less(m.end.y))
		return l.start.y.Less(m.end.y)

	// l vertical m not, do case analysis on new yorker subscription card
	} else if l.Vertical() {
		//fmt.Fprintf(os.Stderr," by l vertical: %v (cmp %v) || %v || %v (cmp %v) && %v\n", m.Below(l.end),l.end.y.Less(m.start.y),m.On(l.end),m.Below(l.start),l.start.y.Less(m.end.y) ,m.Below(p) )
		return m.Below(l.end) || m.On(l.end) || (m.Below(l.start) && m.Below(p) )

	// m vertical, l not, do above but take negation
	} else if m.Vertical() {
		//fmt.Fprintf(os.Stderr," by m vertical: %v\n",l.Below(p))
		return !(l.Below(m.end) || l.On(m.end) || (l.Below(m.start) && l.Below(p) )) 

	// l, m both horizontal
	} else if l.Horizontal() && m.Horizontal() {
		//fmt.Fprintf(os.Stderr," by both horizontal: %v\n", l.start.y.Less(m.start.y))
		return l.start.y.Less(m.start.y)
	}

	// General case:
	// derived from l(p.x) <? m(p.x) but compensates for possible non-invertibility
	// of deltax
	lhs := m.deltax.Mul(l.deltay.Mul(p.x.Sub(l.start.x)).Add(l.start.y.Mul(l.deltax)))
	rhs := l.deltax.Mul(m.deltay.Mul(p.x.Sub(m.start.x)).Add(m.start.y.Mul(m.deltax)))
	//fmt.Fprintf(os.Stderr," by general case: %v\n",lhs.Less(rhs))
	//fmt.Fprintf(os.Stderr,"lhs: %v rhs: %v\n",lhs,rhs)
	return lhs.Less(rhs)
}

func (l *Line) IntersectAsIntegers(m *Line) (bool, *Point) {
	if l.start.Equal(m.start) {
		return true,l.start
	} else if l.start.Equal(m.end) {
		return true,l.start
	} else if l.end.Equal(m.start) {
		return true,l.end
	} else if l.end.Equal(m.end) {
		return true,l.end
	}

	det := l.start.x.Sub(l.end.x).Mul(m.start.y.Sub(m.end.y)).Sub(
		m.start.x.Sub(m.end.x).Mul(l.start.y.Sub(l.end.y)))
	if det.Equal(&Integer{0,0}) {
		return false,nil
	}	

	xnumerator := l.start.x.Mul(l.end.y).Sub(l.start.y.Mul(l.end.x)).Mul(m.start.x.Sub(m.end.x)).Sub(
			m.start.x.Mul(m.end.y).Sub(m.start.y.Mul(m.end.x)).Mul(l.start.x.Sub(l.end.x)))
	ynumerator := l.start.x.Mul(l.end.y).Sub(l.start.y.Mul(l.end.x)).Mul(m.start.y.Sub(m.end.y)).Sub(
			m.start.x.Mul(m.end.y).Sub(m.start.y.Mul(m.end.x)).Mul(l.start.y.Sub(l.end.y)))

	x,okx := xnumerator.Div(det)
	y,oky := ynumerator.Div(det)

	if okx && oky {
		isect := &Point{x,y}
		return l.On(isect) && m.On(isect), isect
	}
	return false,nil
}

func (l *Line) IntersectAsFloats(m *Line) (bool, float64,float64) {
	if l.start.Equal(m.start) {
		return true,l.start.x.Float64(),l.start.y.Float64()
	} else if l.start.Equal(m.end) {
		return true,l.start.x.Float64(),l.start.y.Float64()
	} else if l.end.Equal(m.start) {
		return true,l.end.x.Float64(),l.end.y.Float64()
	} else if l.end.Equal(m.end) {
		return true,l.end.x.Float64(),l.end.y.Float64()
	}

	det := l.start.x.Sub(l.end.x).Mul(m.start.y.Sub(m.end.y)).Sub(
		m.start.x.Sub(m.end.x).Mul(l.start.y.Sub(l.end.y)))
	if det.Equal(&Integer{0,0})  {
		return false,0,0
	}	

	xnumerator := l.start.x.Mul(l.end.y).Sub(l.start.y.Mul(l.end.x)).Mul(m.start.x.Sub(m.end.x)).Sub(
			m.start.x.Mul(m.end.y).Sub(m.start.y.Mul(m.end.x)).Mul(l.start.x.Sub(l.end.x)))
	ynumerator := l.start.x.Mul(l.end.y).Sub(l.start.y.Mul(l.end.x)).Mul(m.start.y.Sub(m.end.y)).Sub(
			m.start.x.Mul(m.end.y).Sub(m.start.y.Mul(m.end.x)).Mul(l.start.y.Sub(l.end.y)))
	x := xnumerator.Float64()/det.Float64()
	y := ynumerator.Float64()/det.Float64()	
		
	return l.FloatOn(x,y) && m.FloatOn(x,y),x,y
}
