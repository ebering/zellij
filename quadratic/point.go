package quadratic

import (
	"os"
	"strconv"
	"regexp"
)

type Point struct {
	x,y *Integer
}

func NewPoint( x,y *Integer) (*Point) {
	return &Point{x,y}
}

func (p *Point) Copy() (*Point) {
	return &Point{p.x.Copy(),p.y.Copy()}
}

func (p * Point) X() *Integer {
	return p.x
}

func (p* Point) Y() *Integer {
	return p.y
}

func (p *Point) Float64() (float64,float64) {
	return p.x.Float64(),p.y.Float64()
}

func (p *Point) SetX(a *Integer){
	p.x = a
}

func (p *Point) SetY(a *Integer){
	p.y = a
}

func (p *Point) Less(q *Point) (bool){
	if p.x.Equal(q.x) {
		return p.y.Less(q.y)
	}
	return p.x.Less(q.x) 
}

func (p *Point) Equal(q *Point) (bool){
	return p.x.Equal(q.x) && p.y.Equal(q.y)
}

func DistanceSquared( p,q *Point) (* Integer) {
	dy := q.y.Sub(p.y)
	dx := q.x.Sub(p.x)
	return dy.Mul(dy).Add(dx.Mul(dx))
}

func MakeTranslation(from,to *Point) (func (*Point) (*Point)) {
	deltax := to.x.Sub(from.x)
	deltay := to.y.Sub(from.y)
	
	return func (p *Point) (*Point) {
		return &Point{p.x.Add(deltax),p.y.Add(deltay)}
	}
}

func (p *Point) RotatePi4(n int) {
	n = n % 8 + 8;

	for i := 0; i < n; i++ {
		p.x,p.y = p.x.MultR2On2().Sub(p.y.MultR2On2()),
			p.y.MultR2On2().Add(p.x.MultR2On2())
	}
}

func PointFromString(ptstr string)  (* Point,os.Error) {
	re := regexp.MustCompile("(-?[0-9]+),(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)")
	matches := re.FindStringSubmatch(ptstr)
	if len(matches) < 5 {
		return nil,os.NewError("invalid pointstring")
	}
	x,ex :=strconv.Atoi64(matches[1])
	x2,ex2 :=strconv.Atoi64(matches[2])
	y,ey :=strconv.Atoi64(matches[3])
	y2,ey2 :=strconv.Atoi64(matches[4])
	if ex != nil || ex2 != nil || ey != nil || ey2 != nil {
		return &Point{new(Integer),new(Integer)},os.NewError("point.FromString: failed to parse coordinates.")
	}
	return &Point{&Integer{x,x2},&Integer{y,y2}},nil
}

func PointMustFromString(ptstr string) (* Point) {
	pt,ok := PointFromString(ptstr)
	if ok != nil {
		panic("PointMustFromString: error in point string: "+ok.String())
	}
	return pt
}	

