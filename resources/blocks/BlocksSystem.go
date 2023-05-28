package blocks

import (
	"image/color"
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"

	"github.com/tidwall/pinhole"
)

const MESSAGE_DESCRIBE = "describe"
const MESSAGE_DESCRIPTION = "description"
const MESSAGE_CREATE_PNG = "create_png"

type BlocksSystem struct {
	base api.System
}

func CreateBlocksSystem(base api.System) *BlocksSystem {
	return &BlocksSystem{
		base: base,
	}
}

func (system *BlocksSystem) HandleRequest(request mentalese.Request) {
	switch request.MessageType {
	case MESSAGE_DESCRIBE:
		scene := "dom:at(E, X, Z, Y) go:has_sort(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)"
		bindings := system.base.RunRelationSetString(central.NO_RESOURCE, scene)
		system.base.GetClientConnector().SendToClient(central.NO_RESOURCE, MESSAGE_DESCRIPTION, bindings.AsSimple())
	case MESSAGE_CREATE_PNG:
		createImage(system.base)
		system.base.GetClientConnector().SendToClient(central.NO_RESOURCE, mentalese.MessageAcknowledge, "")
	default:
		system.base.HandleRequest(request)
	}
}

func (system *BlocksSystem) RunRelationSet(resource string, relationSet mentalese.RelationSet) mentalese.BindingSet {
	return system.base.RunRelationSet(resource, relationSet)
}

func (system *BlocksSystem) RunRelationSetString(resource string, relationSet string) mentalese.BindingSet {
	return system.base.RunRelationSetString(resource, relationSet)
}

func (system *BlocksSystem) GetClientConnector() api.ClientConnector {
	return system.base.GetClientConnector()
}

//
// =======================================================
//

// Using Pinhole https://github.com/tidwall/pinhole to render the scene to a png
// go get -u github.com/tidwall/pinhole
func createImage(system api.System) {

	p := pinhole.New()

	//data := system.Query("dom:at(E, X, Z, Y) go:has_sort(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)")
	scene := "dom:at(E, X, Z, Y) go:has_sort(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)"
	bindings := system.RunRelationSetString(central.NO_RESOURCE, scene)

	p.DrawCube(-.99, -.99, -.99, .99, .99, .99)

	scale := 500.0
	zScale := 1200.0

	for _, binding := range bindings.GetAll() {

		p.Begin()

		x, _ := strconv.ParseFloat(binding.MustGet("X").TermValue, 64)
		y, _ := strconv.ParseFloat(binding.MustGet("Y").TermValue, 64)
		z, _ := strconv.ParseFloat(binding.MustGet("Z").TermValue, 64)
		theType := binding.MustGet("Type").TermValue
		theColor := binding.MustGet("Color").TermValue
		width, _ := strconv.ParseFloat(binding.MustGet("Width").TermValue, 64)
		length, _ := strconv.ParseFloat(binding.MustGet("Length").TermValue, 64)
		height, _ := strconv.ParseFloat(binding.MustGet("Height").TermValue, 64)

		x1 := (x - 500) / scale
		y1 := (y - 500) / scale
		z1 := (z + 50) / zScale

		x2 := x1 + width/scale
		y2 := y1 + height/scale
		z2 := z1 + length/zScale

		if theType == "pyramid" {
			drawPyramid(p, x1, y1, z1, width/scale, height/scale, length/zScale)
		} else {
			p.DrawCube(x1, y1, z1, x2, y2, z2)
		}

		switch theColor {
		case "red":
			p.Colorize(color.RGBA{200, 0, 0, 255})
		case "green":
			p.Colorize(color.RGBA{0, 200, 0, 255})
		case "blue":
			p.Colorize(color.RGBA{0, 0, 200, 255})
		default:
			p.Colorize(color.RGBA{0, 0, 0, 200})
		}

		p.End()
	}

	p.SavePNG(common.Dir()+"/blocksworld.png", 800, 800, nil)
}

func drawPyramid(p *pinhole.Pinhole, x float64, y float64, z float64, width float64, height float64, length float64) {
	topX := x + width/2
	topY := y + height
	topZ := z + length/2

	p.DrawLine(x, y, z, x+width, y, z)
	p.DrawLine(x+width, y, z, x+width, y, z+length)
	p.DrawLine(x+width, y, z+length, x, y, z+length)
	p.DrawLine(x, y, z+length, x, y, z)

	p.DrawLine(x, y, z, topX, topY, topZ)
	p.DrawLine(x+width, y, z, topX, topY, topZ)
	p.DrawLine(x+width, y, z+length, topX, topY, topZ)
	p.DrawLine(x, y, z+length, topX, topY, topZ)
}
