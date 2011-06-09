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
		if e.start.Less(e.end) {
			e.Line().Draw(ctx)
		}
	})
	m.Verticies.Do(func (v interface{}) {
		v.(*Vertex).Draw(ctx)
	})
}

func (m *Map) ColourFaces(ctx *cairo.Surface)  {
	m.Faces.Do(func (f interface{}) {
		F,_ := f.(*Face)
		e := F.boundary;
		ctx.MoveTo(e.start.Float64())
		for f:= e.Next(); f != e; f = f.Next() {
			ctx.LineTo(f.start.Float64())
		}
		ctx.ClosePath()
		ctx.SetSourceRGBA(1.,0.,0.,.1)
		ctx.Fill()
	})
}	
