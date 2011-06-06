package quadratic

import "math"

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

func (l *Line) Below(p *Point) (bool) {
	lhs := l.deltax.Mul(p.y.Sub(l.start.y))
	rhs := l.deltay.Mul(p.x.Sub(l.start.x))

	return rhs.Less(lhs)
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

