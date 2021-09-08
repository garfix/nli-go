package mentalese

const NodeTypePrevSentence = "prev_sentence"
const NodeTypeNextSibling = "next_sibling"
const NodeTypePrevSibling = "prev_sibling"
const NodeTypeSibling = "sibling"
const NodeTypeChild = "child"
const NodeTypeParent = "parent"

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

func (c CategoryPathNode) BindSingle(binding Binding) CategoryPathNode {

	newVariables := []string{}

	for _, variable := range c.Variables {
		value, found := binding.Get(variable)
		if found {
			newVariables = append(newVariables, value.TermValue)
		} else {
			newVariables = append(newVariables, variable)
		}
	}

	newNode := CategoryPathNode{
		NodeType:      c.NodeType,
		Category:      c.Category,
		Variables:     newVariables,
		AllowIndirect: c.AllowIndirect,
	}

	return newNode
}

func (c CategoryPathNode) Copy() CategoryPathNode {
	return CategoryPathNode{
		NodeType:      c.NodeType,
		Category:      c.Category,
		Variables:     c.Variables,
		AllowIndirect: c.AllowIndirect,
	}
}