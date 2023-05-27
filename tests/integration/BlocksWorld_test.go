package tests

import (
	"nli-go/lib/server"
	"testing"
	"time"
)

// Mimics Terry Winograd's SHRDLU dialog, but in the NLI-GO way

func TestBlocksWorld(t *testing.T) {

	srv := server.NewServer("3334")
	srv.RunInBackground()
	defer srv.Close()

	time.Sleep(500 * time.Millisecond)

	client := server.CreateTestClient("blocks")
	defer client.Close()

	client.Run([]server.Test{
		{H: "Does the table support the big red block?", C: "Yes", Clarifications: []string{}},
		{H: "Pick up a big red block", C: "OK", Clarifications: []string{}},
		{H: "Does the table support the big red block?", C: "No", Clarifications: []string{}},
		// original "I don't understand which pyramid you mean"
		{H: "Grasp the pyramid", C: "I don't understand which one you mean", Clarifications: []string{}},

		{H: "Is the blue block in the box?", C: "No", Clarifications: []string{}},
		// original "By "it", I assume you mean the block which is taller than the one I am holding"
		{H: "Find a block which is taller than the one you are holding and put it into the box.", C: "OK", Clarifications: []string{}},
		{H: "Is the blue block in the box?", C: "Yes", Clarifications: []string{}},

		{H: "What does the box contain?", C: "The blue pyramid and the blue block", Clarifications: []string{}},
		{H: "What is the pyramid supported by?", C: "The box", Clarifications: []string{}},
		{H: "How many blocks are not in the box?", C: "Four of them", Clarifications: []string{}},
		{H: "Is at least one of them narrower than the one which I told you to pick up?", C: "Yes, the red cube", Clarifications: []string{}},
		{H: "Is it supported?", C: "Yes, by the table", Clarifications: []string{}},
		// original: No
		{H: "Can the table pick up blocks?", C: "It doesn't make sense", Clarifications: []string{}},
		{H: "Can a pyramid be supported by a block?", C: "Yes", Clarifications: []string{}},
		{H: "Can a pyramid support a pyramid?", C: "I don't know", Clarifications: []string{}},
		{H: "Stack up two pyramids.", C: "I can't", Clarifications: []string{}},

		{H: "Do I own the blue pyramid?", C: "No", Clarifications: []string{}},
		{H: "The blue pyramid is mine", C: "I understand", Clarifications: []string{}},
		{H: "Do I own the blue pyramid?", C: "Yes", Clarifications: []string{}},

		{H: "Do I own a green block?", C: "No", Clarifications: []string{}},
		{H: "I own blocks which are not red, but I don't own anything which supports a pyramid", C: "I understand", Clarifications: []string{}},
		{H: "Do I own a green block?", C: "Yes", Clarifications: []string{}},
		{H: "Do I own all green blocks?", C: "No", Clarifications: []string{}},

		{H: "Do I own the box?", C: "No", Clarifications: []string{}},

		{H: "Do I own anything in the box?", C: "Yes, two things: the blue block and the blue pyramid", Clarifications: []string{}},

		{H: "Does a green block support a pyramid?", C: "Yes", Clarifications: []string{}},
		{H: "Will you please stack up both of the red blocks and either a green cube or a pyramid?", C: "OK", Clarifications: []string{}},
		{H: "Is the small red block supported by a green block?", C: "Yes", Clarifications: []string{}},
		{H: "Is a green block supported by the big red block?", C: "Yes", Clarifications: []string{}},
		{H: "Does a green block support a pyramid?", C: "Yes", Clarifications: []string{}},

		{H: "Which cube is sitting on the table?", C: "The large green one which supports the red pyramid", Clarifications: []string{}},
		{H: "Is there a large block behind a pyramid?", C: "Yes, three of them: a large red one, a large green cube and the blue one", Clarifications: []string{}},
		{H: "Put a small one onto the green cube which supports a pyramid", C: "OK", Clarifications: []string{}},

		{H: "Does the small red block support the green pyramid?", C: "No", Clarifications: []string{}},
		{H: "Put the littlest pyramid on top of it", C: "OK", Clarifications: []string{}},
		{H: "Does the small red block support the green pyramid?", C: "Yes", Clarifications: []string{}},

		// original: I'm not sure what you mean by "on top of" in the phrase "on top of green cubes", do you mean:
		{H: "How many things are on top of green cubes?", C: "Three of them", Clarifications: []string{"1"}},

		{H: "Had you touched any pyramid before you put the green one on the little cube?", C: "Yes, the green one", Clarifications: []string{}},

		// original: While I was stacking up the red cube, a large red block and a large green cube
		{H: "When did you pick it up?", C: "While I was stacking up a large red block, a large green cube and the red cube", Clarifications: []string{}},
		{H: "Why?", C: "To get rid of it", Clarifications: []string{}},
		{H: "Why did you do that?", C: "To clean off the red cube", Clarifications: []string{}},
		{H: "Why did you clear off that cube?", C: "To put it on a large green cube", Clarifications: []string{}},
		// original: To stack up the red cube, a large red block and a large green cube
		{H: "Why did you do that?", C: "To stack up a large red block, a large green cube and the red cube", Clarifications: []string{}},
		{H: "Why did you do that?", C: "Because you asked me to", Clarifications: []string{}},
		// original: By putting a large red block on the table; then letting go of it; then putting a large green cube on it; then letting go of that cube; then putting the red cube on top of that cube; then letting go of that cube
		{H: "How did you do it?", C: "By putting a large red block on the table ; then putting a large green cube on that large red block ; then putting the red cube on that large green cube", Clarifications: []string{}},
		{H: "How many objects did you touch while you were doing it?", C: "Four of them", Clarifications: []string{}},

		{H: "What did the red cube support before you started to clean it off?", C: "The green pyramid", Clarifications: []string{}},
		// original: No, only four of them: the red cube, two large green cubes and a large red block
		{H: "There were five blocks to the left of the box then.", C: "No, only four of them: the red cube, two large green cubes and the large red block", Clarifications: []string{}},
		{H: "Put the blue pyramid on the block in the box.", C: "OK", Clarifications: []string{}},
		{H: "Is there anything which is bigger than every pyramid but is not as wide as the thing that supports it?", C: "Yes, the blue block", Clarifications: []string{}},
		// original: Respond immediately after typeing "steeple" (without the return), and say 'Sorry, I don't know the word "steeple"'
		{H: "Does a steeple", C: "Sorry, I don't know the word \"steeple\"", Clarifications: []string{}},
		{H: "A \"steeple\" is a stack which contains two green cubes and a pyramid.", C: "I understand", Clarifications: []string{}},
		{H: "Are there any steeples now?", C: "No", Clarifications: []string{}},
		{H: "Build one", C: "OK", Clarifications: []string{}},
		{H: "Are there any steeples now?", C: "Yes, a large green cube", Clarifications: []string{}},
		{H: "Does a green cube support a green cube?", C: "Yes", Clarifications: []string{}},
		{H: "Call the biggest block \"superblock\".", C: "OK", Clarifications: []string{}},
		{H: "Have you picked up superblock since we began?", C: "Yes", Clarifications: []string{}},
		// original: By "it", I assume you mean the superblock.
		{H: "Why did you drop it?", C: "I can't explain a nonexistent event", Clarifications: []string{}},
		// original: Yes, four things: the box, the red cube, the blue pyramid and the block named superblock
		{H: "Is there anything to the right of the red pyramid?", C: "Yes, four things: the red cube, a large red block, the block named superblock and the box", Clarifications: []string{}},
		// original: You're welcome!
		{H: "Thank you", C: "You're welcome !", Clarifications: []string{}},
	})

}
