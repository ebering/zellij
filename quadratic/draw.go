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
	m.DoEdges(func (e *Edge) {
		if e.start.Less(e.end.Point) {
			e.Line().Draw(ctx)
		}
	})
	m.DoVerticies(func (v *Vertex) {
		v.Draw(ctx)
	})
}
