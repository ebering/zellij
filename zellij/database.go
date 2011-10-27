package zellij

import "../quadratic/quadratic"
import "os"
import "json"
import "strings"
import "path"
import "fmt"

type Database []*Motif

type Motif struct {
	Name string
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
	d := make([]*Motif,0)
	dirHandle,ok := os.Open(dir)
	if ok != nil {
		panic("database failure" +ok.String())
	}

	dbnames,_ := dirHandle.Readdirnames(0)

	LOAD: for _,fname := range(dbnames) {
		if !strings.HasSuffix(fname,".zellij") {
			continue
		}
		entry,_ := os.Open(path.Join(dir,fname))
		data := make([]*quadratic.Map,2)
		data[0],data[1] = quadratic.NewMap(),quadratic.NewMap()
		json.NewDecoder(entry).Decode(&data)
		
		for _,m := range(d) {
			if m.Frame.Isomorphic(data[0]) {
				m.Variations = append(m.Variations,data[1])
				continue LOAD
			}
		}
		d = append(d,&Motif{fname,data[0],[]*quadratic.Map{data[1]},nil})
	}

	fmt.Fprintf(os.Stderr,"Loaded database with %v frames\n",len(d))

	return d
}
