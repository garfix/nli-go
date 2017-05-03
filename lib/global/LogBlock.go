package global

// a log block is a section of a log whose contents is a distinct part of the cli process
type LogBlock struct {
    title string
    lines []string
    success bool
}

func NewLogBlock(title string) *LogBlock {
    return &LogBlock{ title: title, lines: []string{}, success: true }
}

func (block *LogBlock) AddLine(line string) {
    block.lines = append(block.lines, line)
}

func (block *LogBlock) Fail() {
    block.success = false
}

func (block *LogBlock) IsOk() bool {
    return block.success
}

func (block *LogBlock) String() string {

    s := block.title + "\n"

    if block.success == false {
        s += "> FAILED" + "\n"
    }

    for _, line := range block.lines {
        s += line + "\n"
    }

    return s
}
