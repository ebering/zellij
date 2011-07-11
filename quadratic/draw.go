package quadratic

import "cairo"
import "math"

func (p *Point) Draw(ctx *cairo.Surface) {
	ctx.Arc(p.x.Float64(),p.y.Float64(),.1,0.,2.*math.Pi)
	ctx.Fill()
}

func (l *Line) DrawEndpoints(ctx *cairo.Surface) {
	ctx.MoveTo(l.start.Float64())
	ctx.LineTo(l.end.Float64())
	ctx.Stroke()
	l.start.Draw(ctx)
	l.end.Draw(ctx)
}

func (l *Line) Draw(ctx *cairo.Surface) {
	ctx.MoveTo(l.start.Float64())
	ctx.LineTo(l.end.Float64())
	ctx.Stroke()
}

func (m *Map) DrawEdges(ctx *cairo.Surface) {
	m.Edges.Do(func (f interface{} ) {
		e,_ := f.(*Edge)
		e.Line().Draw(ctx)
	})
	m.Verticies.Do(func (v interface{}) {
		ctx.SetSourceRGBA(0.,0.,0.,1.)
		v.(*Vertex).Draw(ctx)
	})
}

func (m *Map) DrawDebugEdges(ctx *cairo.Surface) {
	m.Edges.Do(func (f interface{} ) {
		e,_ := f.(*Edge)
		ctx.SetSourceRGBA(0.,float64((e.Generation*20)%255)/255.,0.,1.)
		e.Line().Draw(ctx)
	})
	m.Verticies.Do(func (v interface{}) {
		ctx.SetSourceRGBA(0.,0.,0.,1.)
		v.(*Vertex).Draw(ctx)
	})
}

func (m *Map) ColourDebugFaces(ctx *cairo.Surface)  {
	n := float64(m.Faces.Len())
	i := 0.
	m.Faces.Do(func (f interface{}) {
		F,_ := f.(*Face)
		if F.Value.(string) == "outer"  { return }
		e := F.boundary;
		ctx.MoveTo(e.start.Float64())
		for f:= e.Next(); f != e; f = f.Next() {
			ctx.LineTo(f.start.Float64())
		}
		ctx.ClosePath()
		if F.Value.(string) == "active" {
			ctx.SetSourceRGBA(0.,0.,1.,1.)
		} else {
			ctx.SetSourceRGBA(i/n,0.,0.,1.)
		}
		ctx.Fill()
		i = i+1.
	})
}	

func (m *Map) ColourFaces(ctx *cairo.Surface, brush func (*Face) (float64,float64,float64,float64) ) {
	m.Faces.Do(func (f interface {}) {
		F := f.(*Face)
		ctx.SetSourceRGBA(brush(F))
		e := F.boundary;
		ctx.MoveTo(e.start.Float64())
		for f:= e.Next(); f != e; f = f.Next() {
			ctx.LineTo(f.start.Float64())
		}
		ctx.ClosePath()
		ctx.Fill()
	})
}
