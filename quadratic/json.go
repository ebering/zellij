package quadratic

import (
	"json"
	"os"
	"strconv"
	"fmt"
	"container/vector"
)

type mapJSON struct {
	Verts []*Vertex
	Edges []edgeJSON
	Faces []faceJSON
}

type edgeJSON struct {
	Start,End string
}

type faceJSON struct {
	Start,End string
	Value interface{}
	Type string
}

type intJSON struct {
	A,B string
}

func (p *Point) MarshalJSON() ([]byte,os.Error) {
	x,ex := json.Marshal(p.x)
	if ex != nil {
		return nil,ex
	}
	y,ey := json.Marshal(p.y)
	if ey !=nil {
		return nil,ey
	}
	return []byte(fmt.Sprintf("{\"x\":%s,\"y\":%s}",x,y)),nil
}

func (i *Integer) MarshalJSON() ([]byte,os.Error) {
	return json.Marshal(intJSON{strconv.Itoa64(i.a),strconv.Itoa64(i.b)})
}

func (m *Map) MarshalJSON() ([]byte,os.Error) {
	jv := new(mapJSON)
	jv.Verts = make([]*Vertex,m.Verticies.Len())
	jv.Edges = make([]edgeJSON,0,m.Edges.Len()/2)
	jv.Faces = make([]faceJSON,0,m.Faces.Len())
	vertindex := make(map[*Vertex]string)
	for i:=0; i < m.Verticies.Len(); i++ {
		jv.Verts[i] = m.Verticies.At(i).(*Vertex)
		vertindex[m.Verticies.At(i).(*Vertex)] = strconv.Itoa(i)
	}
	m.Edges.Do(func (e interface{}) {
		E := e.(*Edge)
		if E.start.Less(E.end) {
			jv.Edges = append(jv.Edges,edgeJSON{vertindex[E.start],vertindex[E.end]})
		}
	})
	m.Faces.Do(func (f interface{}) {
		jv.Faces = append(jv.Faces,faceJSON{vertindex[f.(*Face).boundary.start],vertindex[f.(*Face).boundary.end],f.(*Face).Value,f.(*Face).Type})
	})
	return json.Marshal(jv)
}

func (m *Map) UnmarshalJSON(js []byte) os.Error {
	jv := new(mapJSON)
	err := json.Unmarshal(js,jv)
	if err != nil {
		return err
	}
	for _,v := range(jv.Verts) {
		m.Verticies.Push(v)
	}
	for _,e := range(jv.Edges) {
		st,_ := strconv.Atoi64(e.Start)
		en,_ := strconv.Atoi64(e.End)
		m.JoinVerticies(jv.Verts[st],jv.Verts[en])
	}
	for _,f := range(jv.Faces) {
		st,_ := strconv.Atoi64(f.Start)
		en,_ := strconv.Atoi64(f.End)
		nf := new(Face)
		nf.fromMap = m
		nf.Value = f.Value
		nf.Type = f.Type
		bdyStart := jv.Verts[st]
		bdyStart.OutgoingEdges.Do(func (e interface{}) {
			if e.(*Edge).end == jv.Verts[en] {
				nf.boundary = e.(*Edge)
			}
		})
		nf.DoEdges(func (e *Edge) {
			e.face = nf
		})
		m.Faces.Push(nf)
	}

	if len(jv.Faces) == 0 {
		edges := m.Edges.Copy()

		for edges.Len() > 0 {
			F := new(Face)
			F.boundary = edges.Pop().(*Edge)
			F.boundary.face = F
			for e := F.boundary.next; e != F.boundary; e = e.next {
				e.face = F
				for i := 0; i < edges.Len();  {
					if edges.At(i).(*Edge) == e {
						edges.Delete(i)
					} else {
						i++
					}
				}
			}
			m.Faces.Push(F)
		}
	}

	m.Init()
	m.Edges.Do(func(f interface{}) {
		if f.(*Edge).face == nil {
			panic("unmarshal horrific failure")
		}
	})
	
	return nil
}

func (v *Vertex) UnmarshalJSON(js []byte) os.Error {
	vtx := make(map[string]intJSON)
	err := json.Unmarshal(js,&vtx)
	if err != nil {
		return err
	}
	v.Point = &Point{vtx["x"].Integer(),vtx["y"].Integer()}
	v.OutgoingEdges = new(vector.Vector)
	return nil
}

func (i intJSON) Integer() *Integer {
	a, _ := strconv.Atoi64(i.A)
	b, _ := strconv.Atoi64(i.B)
	return &Integer{a,b}
}

