package zellij

import "../quadratic/quadratic"
import "os"
import "json"

type Database []*Motif

type Motif struct {
	Frame *quadratic.Map
	Variations []*quadratic.Map
	Symmetries []string
}

func (d Database) Matching(frame *quadratic.Map) (*Motif,quadratic.Isomorphism,os.Error) {
	for _,m := range(d) {
		iso,ok := m.Frame.Isomorphism(frame)
		if ok == nil {
			return m,iso,nil
		}
	}
	return nil,nil,os.NewError("no match in database")
}

func LoadDatabase(dir string) Database {
	d := make([]*Motif,1)
	dirHandle,ok := os.Open(dir)
	if ok != nil {
		panic("database failure" +ok.String())
	}

	dirHandle.Chdir()

	dbnames,_ := dirHandle.Readdirnames(0)

	for _,fname := range(dbnames) {
		entry,_ := os.Open(fname)
		data := make([]*quadratic.Map,2)
		data[0],data[1] = quadratic.NewMap(),quadratic.NewMap()
		json.NewDecoder(entry).Decode(data)
		
		for _,m := range(d) {
			if m.Frame.Isomorphic(data[0]) {
				m.Variations = append(m.Variations,data[1])
				break
			} else {
				d = append(d,&Motif{data[0],[]*quadratic.Map{data[1]},nil})
				break
			}
		}
	}

	return d
}
