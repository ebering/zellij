package main

import "cairo"
import "flag"
import "json"
import "fmt"
import "path"
import "os"

import "runtime"

import "./quadratic/quadratic"
import "./zellij/zellij"

func init() {
	runtime.GOMAXPROCS(4)
	zellij.Workers = 2
	flag.StringVar(&skeleton,"skeleton","0246","the zellij skeleton as a string of headings")
	flag.StringVar(&tileSymmetry,"sym","d4","the minimum symmetry")
	flag.StringVar(&dir,"dir","tiles","output directory")
}

var ZellijTilings <-chan *quadratic.Map
var reset chan int
var skeleton string
var dir string
var tileSymmetry string

func main() {	
	flag.Parse()
	ZellijTilings, reset = zellij.TileSkeleton(skeleton, tileSymmetry, false)
	Frame := zellij.SkeletonFrame(skeleton)
	symmetryCounts := make(map[string]int)
	for T := range(ZellijTilings) {
		symmetryGroup := zellij.DetectSymmetryGroup(T)
		filename := fmt.Sprintf("%v-%v-%d",skeleton,symmetryGroup,symmetryCounts[symmetryGroup]+1)
		symmetryCounts[symmetryGroup]++
		save, err := os.Create(path.Join(dir,filename+".zellij"))
		if err != nil {
			panic("file error")
		}
		enc := json.NewEncoder(save)
		enc.Encode([]*quadratic.Map{Frame,T})
		save.Close()
		image := cairo.NewSurface(path.Join(dir,"svg",filename+".svg"),72*5,72*5)
		image.SetSourceRGB(0., 0., 0.)
		image.SetLineWidth(.1)
		image.Translate(72*2.5, 72*2.5)
		image.Scale(4., 4.)
		T.ColourFaces(image, zellij.PlainBrush)
		image.Finish()
	}
}

