package main

import "cairo"
import "http"
import "log"
import "os"

import "./quadratic/quadratic"
import "./zellij"

func init() {
	ZellijTilings,reset = zellij.TileRegion(quadratic.NewInteger(-40,0),quadratic.NewInteger(40,0),quadratic.NewInteger(-40,0),quadratic.NewInteger(40,0))
}

var ZellijTilings <-chan *quadratic.Map
var reset chan<- int

func main() {
	http.HandleFunc("/",MainScreen)
	http.HandleFunc("/start",StartTiling)
	http.HandleFunc("/tiles",RenderTiles)
	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		log.Fatal("ListenAndServe: ",err.String())
	}
}

func MainScreen(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w,req,"ui/ui.html")
}

func RenderTiles(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String()+"\n")
	}	
	
	t := <-ZellijTilings

	/*t := zellij.TileMap(zellij.Tiles[0])
	t.Edges.Do(func (f interface{}) {
		fmt.Fprintf(os.Stderr,"%v heading: %v\n",f,f.(*quadratic.Edge).IntHeading())
	})
	/*u := zellij.TileMap(zellij.Tiles[0]).Translate(quadratic.NewVertex(zellij.Points["s"]),quadratic.NewVertex(zellij.Points["h"]))
	
	v,ok := t.Overlay(u,zellij.Overlay)
	if ok != nil {
		//os.Stderr.WriteString(ok.String()+"\n")
	}
	/*v,ok = v.Overlay(t.Translate(quadratic.NewVertex(zellij.Points["s"]),quadratic.NewVertex(zellij.Points["e"])),zellij.Overlay)
	if ok != nil {
		os.Stderr.WriteString(ok.String()+"\n")
	}
	t = v*/
	
	image := cairo.NewSurface("svg/test-surface.svg",72*4,72*4)
	image.SetSourceRGB(0.,0.,0.)
	image.SetLineWidth(.1)
	image.Translate(72*2.,72*2.)
	image.Scale(4.,4.)
	t.ColourFaces(image)
	image.SetSourceRGBA(0.,0.,0.,1.)
	t.DrawEdges(image)

	image.Finish()
	http.ServeFile(w,req,"svg/test-surface.svg")
}

func StartTiling(w http.ResponseWriter, req *http.Request) {
	if reset != nil {
		reset <- 1
	}

	ZellijTilings,reset = zellij.TileRegion(new(quadratic.Integer),quadratic.NewInteger(40,0),new(quadratic.Integer),quadratic.NewInteger(40,0))
	
	w.WriteHeader(http.StatusOK)
}
	
