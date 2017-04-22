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

func (log *systemLog) String() string {

    s := ""

    for _, block := range log.blocks {
        s += block.String()
    }

    return s
}