package tests

import (
	"fmt"
	"github.com/tidwall/pinhole"
	"image/color"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"strconv"
	"testing"
)

// Mimics some of SHRDLU's functions, but in the nli-go way

func TestBlocksWorld(t *testing.T) {

	var tests = [][]struct {
		question      string
		answer        string
	}{
		{
				{"Does the table support the big red block?", "Yes"},
			{"Pick up a big red block", "OK"},
				{"Does the table support the big red block?", "No"},

			// original "I don't understand which pyramid you mean"
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

				{"Do I own the blue pyramid?", "No"},
			{"The blue pyramid is mine", "I understand"},
				{"Do I own the blue pyramid?", "Yes"},

				{"Do I own a green block?", "No"},
			{"I own blocks which are not red, but I don't own anything which supports a pyramid", "I understand"},
				{"Do I own a green block?", "Yes"},
				{"Do I own all green blocks?", "No"},

			{"Do I own the box?", "No"},

			{"Do I own anything in the box?", "Yes, two things: the blue block and the blue pyramid"},

				{"Does a green block support a pyramid?", "Yes"},
			{"Will you please stack up both of the red blocks and either a green cube or a pyramid?", "OK"},
				{"Is the small red block supported by a green block?", "Yes"},
				{"Is a green block supported by the big red block?", "Yes"},
				{"Does a green block support a pyramid?", "Yes"},

			{"Which cube is sitting on the table?", "The large green one which supports the red pyramid"},
			{"Is there a large block behind a pyramid?", "Yes, three of them: a large red one, a large green cube and the blue one"},
			{"Put a small one onto the green cube which supports a pyramid", "OK"},

				{"Does the small red block support the green pyramid?", "No"},
			{"Put the littlest pyramid on top of it", "OK"},
				{"Does the small red block support the green pyramid?", "Yes"},

			// original: I'm not sure what you mean by "on top of" in the phrase "on top of green cubes", do you mean:
			{"How many things are on top of green cubes?", " [0] Directly on the surface [1] Anywhere on top of"},
			{"1", "Three of them"},

			{"Had you touched any pyramid before you put the green one on the little cube?", "Yes, the green one"},

			// original: While I was stacking up the red cube, a large red block and a large green cube
			{"When did you pick it up?", "While I was stacking up a large red block, a large green cube and the red cube"},
		},
		{
			//{"Stack up 2 green blocks and a small red block", "OK"},
			//{"stack up a blue block and a blue pyramid", "OK1"},
			//{"Put a blue block into the box", "OK"},
			//{"Will you please stack up both of the red blocks and either a green cube or a pyramid?", "OK"},
			//{"Stack up 3 objects", "OK"},
			//{"Stack up 3 objects", "OK"},
		},
	}

	log := common.NewSystemLog()

	//log.SetDebug(true)
	//log.SetPrint(true)
	system := global.NewSystem(common.Dir() + "/../../resources/blocks", "blocks-demo", common.Dir() + "/../../var", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	system.ResetSession()

	for _, session := range tests {

		for _, test := range session {

			log.Clear()

			if test.question == "Put the littlest pyramid on top of it" {
				//log.SetDebug(true)
				//log.SetPrint(true)
			}

			answer, options := system.Answer(test.question)

			if options.HasOptions() {
				answer += options.String()
			}

			if answer != test.answer {
				fmt.Printf(test.question)
				t.Errorf("Test relationships:\nGOT:\n  %v\nWANT:\n  %v", answer, test.answer)
				t.Errorf("\n%s", log.String())
				//t.Errorf("\n%s", log.GetErrors())
			}
		}
	}

	//createGrid(system)

	createImage(system)
}

func createGrid(system *global.System) {
	g := [20][20]string{}
	for _, binding := range system.Query("dom:grid(fixed, H, V, 1)").GetAll() {
		//fmt.Println(binding.String())
		h, _ := strconv.Atoi(binding.MustGet("H").TermValue)
		v, _ := strconv.Atoi(binding.MustGet("V").TermValue)
		g[v][h]="x"
	}
	for h := 9; h < 20; h++ {
		for v := 0; v < 10; v++ {
			if g[v][19-h] == "x" {
				fmt.Print("x")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

// Using Pinhole https://github.com/tidwall/pinhole to render the scene to a png
//
// go get -u github.com/tidwall/pinhole
//
func createImage(system *global.System) {

	p := pinhole.New()

	data := system.Query("dom:at(E, X, Z, Y) dom:type(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)")

	p.DrawCube(-.99, -.99, -.99, .99, .99, .99)

	scale := 500.0
	zScale := 1200.0

	for _, binding := range data.GetAll() {

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