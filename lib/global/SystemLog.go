package global

type systemLog struct {
    blocks []*LogBlock
}

func NewSystemLog() *systemLog {
    return &systemLog{ blocks: []*LogBlock{} }
}

func (log *systemLog) AddBlock(block *LogBlock) {
    log.blocks = append(log.blocks, block)
}

func (log *systemLog) IsOk() bool {
    ok := true

    for _, block := range log.blocks {
        ok = ok && block.IsOk()
    }

    return ok
}

func (log *systemLog) String() string {

    s := ""

    for _, block := range log.blocks {
        s += block.String()
    }

    return s
}

func (log *systemLog) GetLogLines() []string {

    s := []string{}

    for _, block := range log.blocks {
        s = append(s, block.lines...)
    }

    return s
}