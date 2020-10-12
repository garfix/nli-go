package tests

import (
	"github.com/tidwall/pinhole"
	"image/color"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"os"
	"strconv"
	"testing"
)

// Mimics some of SHRDLU's functions, but in the nli-go way

// Using Pinhole https://github.com/tidwall/pinhole to render the scene to a png
//
// go get -u github.com/tidwall/pinhole
//
func TestBlocksWorld(t *testing.T) {
	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/blocks", log)
	sessionId := "1"
	actualSessionPath := common.AbsolutePath(common.Dir(), "sessions/" + sessionId + ".json")

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = [][]struct {
		question      string
		answer        string
	}{
		{
				{"Does the table support the big red block?", "Yes"},
			{"Pick up a big red block", "OK"},
				{"Does the table support the big red block?", "No"},

			// todo "I don't understand which pyramid you mean"
			{"Grasp the pyramid", "I don't understand which one you mean"},

				{"Is the blue block in the box?", "No"},
			// todo "By "it", I assume you mean the block which is taller than the one I am holding"
			{"Find a block which is taller than the one you are holding and put it into the box.", "OK"},
				{"Is the blue block in the box?", "Yes"},

			{"What does the box contain?", "The blue pyramid and the blue block"},
			{"What is the pyramid supported by?", "The box"},
			{"How many blocks are not in the box?", "Four of them"},
			{"Is at least one of them narrower than the one which I told you to pick up?", "Yes, the red cube"},
			{"Is it supported?", "Yes, by the table"},
			{"Can the table pick up blocks?", "No"},
			{"Can a pyramid be supported by a block?", "Yes"},
			// todo: must be: I don't know
			{"Can a pyramid support a pyramid?", "No"},
			{"Stack up two pyramids.", "I can't"},

				//{"Do I own the blue pyramid?", "No"},
			{"The blue pyramid is mine", "I understand"},
				{"Do I own the blue pyramid?", "Yes"},

				{"Do I own a green block?", "No"},
			{"I own blocks which are not red, but I don't own anything which supports a pyramid", "I understand"},
				{"Do I own a green block?", "Yes"},
				{"Do I own all green blocks?", "No"},

			{"Do I own the box?", "No"},

			// todo: must be: Yes, two things: the blue block and the blue pyramid
			{"Do I own anything in the box?", "Yes, the blue block and the blue pyramid"},

				{"Does a green block support a pyramid?", "Yes"},
			{"Will you please stack up both of the red blocks and either a green cube or a pyramid?", "OK"},
				{"Is the small red block supported by a green block?", "Yes"},
				{"Is a green block supported by the big red block?", "Yes"},
				{"Does a green block support a pyramid?", "Yes"},

			{"Which cube is sitting on the table?", "The large green one which supports the red pyramid"},

			//{"Is there a large block behind a pyramid?", "Yes, three of them: a large red one, a large green cube and a blue one"},
		},
		{
		},
	}

	_ = os.Remove(actualSessionPath)

	for _, session := range tests {

		for _, test := range session {

			log.Clear()

			system.PopulateDialogContext(actualSessionPath, false)

			answer, options := system.Answer(test.question)

			if options.HasOptions() {
				answer += options.String()
			}

			system.StoreDialogContext(actualSessionPath)

			if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				t.Error(log.String())
			}
		}
	}

	createImage(system)
}

func createImage(system *global.System) {

	p := pinhole.New()

	data := system.Query("dom:at(E, X, Z, Y) dom:type(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)")

	p.DrawCube(-.99, -.99, -.99, .99, .99, .99)

	scale := 500.0
	zScale := 1200.0

	for _, binding := range data {

		p.Begin()

		x, _ := strconv.ParseFloat(binding["X"].TermValue, 64)
		y, _ := strconv.ParseFloat(binding["Y"].TermValue, 64)
		z, _ := strconv.ParseFloat(binding["Z"].TermValue, 64)
		theType := binding["Type"].TermValue
		theColor := binding["Color"].TermValue
		width, _ := strconv.ParseFloat(binding["Width"].TermValue, 64)
		length, _ := strconv.ParseFloat(binding["Length"].TermValue, 64)
		height, _ := strconv.ParseFloat(binding["Height"].TermValue, 64)

		x1 := (x - 500) / scale
		y1 := (y - 500) / scale
		z1 := (z + 50) / zScale

		x2 := x1 + width / scale
		y2 := y1 + height / scale
		z2 := z1 + length / zScale

		if theType == "pyramid" {
			drawPyramid(p, x1, y1, z1, width / scale, height / scale, length / zScale)
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

	p.SavePNG(common.Dir() + "/blocksworld.png", 800, 800, nil)
}

func drawPyramid(p *pinhole.Pinhole, x float64, y float64, z float64, width float64, height float64, length float64) {
	topX := x + width / 2
	topY := y + height
	topZ := z + length / 2

	p.DrawLine(x, y, z, x + width, y, z)
	p.DrawLine(x + width, y, z, x + width, y, z + length)
	p.DrawLine(x + width, y, z + length, x, y, z + length)
	p.DrawLine(x, y, z + length, x, y, z)

	p.DrawLine(x, y, z, topX, topY, topZ)
	p.DrawLine(x + width, y, z, topX, topY, topZ)
	p.DrawLine(x + width, y, z + length, topX, topY, topZ)
	p.DrawLine(x, y, z + length, topX, topY, topZ)
}