package quadratic

import "container/vector"
import "sort"
import "math"
import "os"

// Stores a vertex, which is a point that has a list of incident edges
// IMPORTANT INVARIANT: The vector of outgoing edges is a vector of things of type *Edge and it is KEPT IN SORTED ORDER BY HEADING
// DON'T COCK THAT UP.
type Vertex struct {
	*Point
	OutgoingEdges *vector.Vector
	copy          *Vertex //used only in the copy routine
	inFace        *Face
}

func NewVertex(p *Point) *Vertex {
	nv := new(Vertex)
	nv.OutgoingEdges = new(vector.Vector)
	nv.Point = p.Copy()
	return nv
}

func (v *Vertex) EdgeIndexCounterClockwiseOf(e *Edge) int {
	return sort.Search(v.OutgoingEdges.Len(), func(i int) bool {
		return e.Less(v.OutgoingEdges.At(i))
	})
}

func (v *Vertex) Less(u interface{}) bool {
	return v.Point.Less(u.(*Vertex).Point)
}

// Stores a directed half edge in a planar map
type Edge struct {
	start, end       *Vertex
	next, prev, twin *Edge
	face             *Face
	newFace          *Face
	visited          bool
	Generation       int
	copy             *Edge //used only in the copy routine
	fromMap *Map // pointer back to the map we're in
}

func (e *Edge) Start() *Vertex {
	return e.start
}

func (e *Edge) End() *Vertex {
	return e.end
}

func (e *Edge) Face() *Face {
	return e.face
}

func (e *Edge) Coterminal(f *Edge) bool {
	return e.start.Point.Equal(f.start.Point) ||
		e.start.Point.Equal(f.end.Point) ||
		e.end.Point.Equal(f.start.Point) ||
		e.end.Point.Equal(f.end.Point)
}

func (e *Edge) Equal(f *Edge) bool {
	return e.start.Point.Equal(f.start.Point) && e.end.Point.Equal(f.end.Point)
}

func (e *Edge) LengthSquared() *Integer {
	return DistanceSquared(e.start.Point, e.end.Point)
}

func (e *Edge) Line() *Line {
	return NewLine(e.start.Point, e.end.Point)
}

func (e *Edge) Heading() float64 {
	return math.Atan2(e.end.y.Sub(e.start.y).Float64(), e.end.x.Sub(e.start.x).Float64())
}

func (e *Edge) IntHeading() int {
	dx := e.end.X().Sub(e.start.X())
	dy := e.end.Y().Sub(e.start.Y())
	negeq := dx.Add(dy).Equal(Zero)
	eq := dx.Sub(dy).Equal(Zero)
	switch {
	case dy.Equal(Zero) && Zero.Less(dx):
		return 0
	case eq && Zero.Less(dy) && Zero.Less(dx):
		return 1
	case dx.Equal(Zero) && Zero.Less(dy):
		return 2
	case negeq && Zero.Less(dy) && dx.Less(Zero):
		return 3
	case dy.Equal(Zero) && dx.Less(Zero):
		return 4
	case eq && dy.Less(Zero) && dx.Less(Zero):
		return 5
	case dx.Equal(Zero) && dy.Less(Zero):
		return 6
	case negeq && dy.Less(Zero) && Zero.Less(dx):
		return 7
	}
	panic("no heading")
	return 0
}

func (e *Edge) Less(v interface{}) bool {
	f := v.(*Edge)
	if e.start.Equal(f.start.Point) && e.end.Equal(f.end.Point) {
		return false
	}

	return e.Heading() < f.Heading()
}

func (e *Edge) Parallel(f *Edge) bool {
	return math.Fabs(e.Heading()-f.Heading()) < FLOAT64_EPSILON || math.Fabs(e.twin.Heading()-f.Heading()) < FLOAT64_EPSILON
}

// Represents a face
type Face struct {
	boundary *Edge
	fromMap  *Map
	Value    interface{}
	copy     *Face
	Type     string
}

func (f *Face) DoEdges(D func(*Edge)) {
	D(f.boundary)
	for l := f.boundary.next; l != f.boundary; l = l.next {
		D(l)
	}
}

func (f *Face) Neighbors() []*Face {
	nh := make(map[*Face]bool)
	f.DoEdges(func (e *Edge) {
		nh[e.twin.face] = true
	})
	ret := make([]*Face,0)
	for n,t := range(nh) {
		if t {
			ret = append(ret,n)
		}
	}
	return ret
}

func (f *Face) Inner() bool {
	leastVtx := f.boundary.start
	for l := f.boundary.next; l != f.boundary; l = l.next {
		if l.start.Less(leastVtx) {
			leastVtx = l.start
		}
	}
	var e2 *Edge
	leastVtx.OutgoingEdges.Do(func(e interface{}) {
		if e.(*Edge).face == f {
			e2 = e.(*Edge)
		}
	})
	return ( e2.IntHeading() == 1 && e2.prev.IntHeading() == 6) ||
		( e2.IntHeading() == 0 && (e2.prev.IntHeading() == 6 || e2.prev.IntHeading() == 5)) ||
		( e2.IntHeading() == 7 && (e2.prev.IntHeading() == 6 || e2.prev.IntHeading() == 5 || e2.prev.IntHeading() == 4))
}

func (f *Face) Frame() *Map {
	verts := make([]*Point,0)
	f.DoEdges(func (e *Edge) {
		verts = append(verts,e.start.Point.Copy())
	})

	return PolygonMap(verts)
}

// Represents a planar map
type Map struct {
	Verticies, Edges, Faces *vector.Vector
	adjacencyMatrix         [][]bool
}

func (m *Map) Init() {
	m.Verticies.Do(func(v interface{}) {
		sort.Sort(v.(*Vertex).OutgoingEdges)
	})
	sort.Sort(m.Verticies)
	m.adjacencyMatrix = m.makeAdjacencyMatrix()
	m.Edges.Do(func (e interface{}) {
		e.(*Edge).fromMap = m
	})
	m.Faces.Do(func (f interface{}) {
		f.(*Face).fromMap = m
	})
}

func NewMap() *Map {
	return &Map{new(vector.Vector), new(vector.Vector), new(vector.Vector), make([][]bool, 1)}
}


// Given two verticies returns a twin pair of half edges between them,
// the first goes from start to end, the second is the twin. 
// Correctly sets next and prev based on the other incident edges, and adds them
// to the relevant set of pointers in the vertex.
func NewEdgePair(start, end *Vertex) (*Edge, *Edge) {
	e, f := new(Edge), new(Edge)
	e.start, e.end = start, end
	f.start, f.end = e.end, e.start

	e.twin, f.twin = f, e

	if start.OutgoingEdges.Len() == 0 {
		start.OutgoingEdges.Push(e)
		e.prev = f
		f.next = e
	} else {
		i := start.EdgeIndexCounterClockwiseOf(e)
		n := start.OutgoingEdges.Len()
		e.prev = start.OutgoingEdges.At(i % n).(*Edge).twin
		f.next = start.OutgoingEdges.At((i - 1 + n) % n).(*Edge)
		start.OutgoingEdges.Insert(i, e)
	}
	f.next.prev = f
	e.prev.next = e

	if end.OutgoingEdges.Len() == 0 {
		end.OutgoingEdges.Push(f)
		e.next = f
		f.prev = e
	} else {
		i := end.EdgeIndexCounterClockwiseOf(f)
		n := end.OutgoingEdges.Len()
		f.prev = end.OutgoingEdges.At(i % n).(*Edge).twin
		e.next = end.OutgoingEdges.At((i - 1 + n) % n).(*Edge)
		end.OutgoingEdges.Insert(i, f)
	}
	f.prev.next = f
	e.next.prev = e

	return e, f
}

// Next returns the next edge along a face
func (e *Edge) Next() *Edge {
	return e.next
}

// Joins f to e as the next of e, updating next and prev pointers for both edges. 
// If this operation would not make sense (i.e. e and f are not at the same vertex)
// it does nothing.
func (e *Edge) JoinNext(f *Edge) {
	if e.end == f.start {
		e.next = f
		f.prev = e
	}
}

// Returns the previous edge along a face
func (e *Edge) Prev() *Edge {
	return e.prev
}

// Joins f to e prior to e instead of after it. Also does not act if this
// operation does not make sense.
func (e *Edge) JoinPrev(f *Edge) {
	if f.end == e.start {
		e.prev = f
		f.next = e
	}
}

// Returns the twin of an edge. Note that twins cannot be set.
func (e *Edge) Twin() *Edge {
	return e.twin
}

func (m *Map) AddVertex(p *Point) *Vertex {
	V := m.Verticies
	for i := 0; i < V.Len(); i++ {
		if V.At(i).(*Vertex).Point.Equal(p) {
			return V.At(i).(*Vertex)
		}
	}
	nv := NewVertex(p)
	m.Verticies.Push(nv)
	return nv
}

func (m *Map) JoinVerticies(u *Vertex, v *Vertex) *Edge {
	e, eTwin := NewEdgePair(u, v)
	m.Edges.Push(e)
	m.Edges.Push(eTwin)
	e.fromMap,eTwin.fromMap = m,m
	return e
}

func (m *Map) Copy() *Map {
	c := NewMap()
	m.Verticies.Do(func(v interface{}) {
		cv := NewVertex(v.(*Vertex).Copy())
		v.(*Vertex).copy = cv
		c.Verticies.Push(cv)
	})
	m.Edges.Do(func(f interface{}) {
		e := f.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			e.copy = c.JoinVerticies(e.start.copy, e.end.copy)
			e.twin.copy = e.copy.twin
			e.copy.Generation = e.Generation
			e.copy.twin.Generation = e.twin.Generation
		}
	})
	m.Faces.Do(func(f interface{}) {
		F := f.(*Face)
		G := new(Face)
		G.fromMap = c
		G.Value = F.Value
		G.Type = F.Type
		G.boundary = F.boundary.copy
		F.copy = G
		c.Faces.Push(G)
	})
	c.Faces.Do(func(f interface{}) {
		F := f.(*Face)
		for e := F.boundary.Next(); e != F.boundary; e = e.Next() {
			e.face = F
		}
		F.boundary.face = F
	})
	c.Edges.Do(func(f interface{}) {
		if f.(*Edge).face == nil {
			panic("map copy horrific failure")
		}
	})
	c.Init()
	return c
}

func (m *Map) Translate(from, to *Vertex) *Map {
	T := MakeTranslation(from.Point, to.Point)
	m.Verticies.Do(func(v interface{}) {
		v.(*Vertex).Point = T(v.(*Vertex).Point)
	})
	return m
}

func (m *Map) RotatePi4(n int) *Map {
	m.Verticies.Do(func(v interface{}) {
		v.(*Vertex).RotatePi4(n)
	})
	m.Init()
	return m
}

func (m *Map) ReflectXAxis() *Map {
	m.Verticies.Do(func(v interface{}) {
		v.(*Vertex).ReflectXAxis()
	})
	m.Edges.Do(func(f interface{}) {
		e := f.(*Edge)
		e.newFace = e.twin.face
	})
	m.Edges.Do(func(f interface{}) {
		e := f.(*Edge)
		e.face = e.newFace
		e.newFace = nil
	})
	m.Init()
	return m
}

func (m *Map) makeAdjacencyMatrix() [][]bool {
	mat := make([][]bool, m.Verticies.Len())
	for i, _ := range mat {
		mat[i] = make([]bool, m.Verticies.Len())
	}
	sort.Sort(m.Verticies)
	for i := 0; i < m.Verticies.Len(); i++ {
		u := m.Verticies.At(i).(*Vertex)
		for j := 0; j < m.Verticies.Len(); j++ {
			if i == j {
				continue
			}
			v := m.Verticies.At(j).(*Vertex)
			u.OutgoingEdges.Do(func(e interface{}) {
				mat[i][j] = mat[i][j] || v.Equal(e.(*Edge).end.Point)
			})
		}
	}
	return mat
}

func (m *Map) AdjacencyMatrix() [][]bool {
	return m.adjacencyMatrix
}

func (m *Map) Isomorphic(n *Map) bool {
	if m.Verticies.Len() != n.Verticies.Len() || m.Edges.Len() != n.Edges.Len() {
		return false
	}

	for i := 0; i < m.Verticies.Len(); i++ {
		u := m.Verticies.At(i).(*Vertex)
		v := n.Verticies.At(i).(*Vertex)
		if u.OutgoingEdges.Len() != v.OutgoingEdges.Len() {
			return false
		}
		for j := 0; j < u.OutgoingEdges.Len(); j++ {
			if u.OutgoingEdges.At(j).(*Edge).IntHeading() != v.OutgoingEdges.At(j).(*Edge).IntHeading() ||
				!u.OutgoingEdges.At(j).(*Edge).LengthSquared().Equal(v.OutgoingEdges.At(j).(*Edge).LengthSquared()) {
				return false
			}
		}
	}

	return true
}

type Isomorphism func(*Map) *Map

func (m *Map) Isomorphism(n *Map) (Isomorphism,os.Error) {
	o := n.Copy()
	for rots :=0; rots < 8; rots++ {
		if m.Isomorphic(o) {
			return func (m *Map) *Map {
				return m.RotatePi4(rots)
			},nil
		}
		o.RotatePi4(1)
	}

	return nil,os.NewError("not isomorphic at any rotation")
}


func (m *Map) Equal(n *Map) bool {
	ret := true
	if !m.Isomorphic(n) {
		return false
	}
	for i:=0; i < m.Verticies.Len(); i++ {
		u := m.Verticies.At(i).(*Vertex)
		v := n.Verticies.At(i).(*Vertex)
		ret = ret && u.Equal(v.Point)
	}
	return ret
}
		

func (m *Map) Centroid() *Point {
	Xsum := new(Integer)
	Ysum := new(Integer)
	n := NewInteger(int64(m.Verticies.Len()),0)
	for i:=0; i < m.Verticies.Len(); i++ {
		v := m.Verticies.At(i).(*Vertex)
		Xsum = Xsum.Add(v.x)
		Ysum = Ysum.Add(v.y)
	}
	X,canX := Xsum.Div(n)
	Y,canY := Ysum.Div(n)

	if !canX || !canY {
		panic("centroid not at a quadratic integer")
	}

	return NewPoint(X,Y)
}

func (m *Map) SetGeneration(g int) {
	m.Edges.Do(func(E interface{}) {
		e := E.(*Edge)
		e.Generation = g
	})
}
