package mentalese

const NodeTypePrevSentence = "prev_sentence"
const NodeTypeNextSibling = "next_sibling"
const NodeTypePrevSibling = "prev_sibling"
const NodeTypeSibling = "sibling"
const NodeTypeChild = "child"
const NodeTypeParent = "parent"

type CategoryPath []CategoryPathNode

type CategoryPathNode struct {
	NodeType      string
	Category      string
	Variables     []string
	AllowIndirect bool
}

func NewCategoryPathNode(nodeType string, value string, variables []string, allowIndirect bool) CategoryPathNode {
	return CategoryPathNode{
		NodeType:      nodeType,
		Category:      value,
		Variables:     variables,
		AllowIndirect: allowIndirect,
	}
}

func (c CategoryPathNode) GetNodeType() string {
	return c.NodeType
}

func (c CategoryPathNode) GetCategory() string {
	return c.Category
}

func (c CategoryPathNode) GetVariables() []string {
	return c.Variables
}

func (c CategoryPathNode) DoesAllowIndirect() bool {
	return c.AllowIndirect
}