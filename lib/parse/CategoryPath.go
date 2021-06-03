package parse

const NodeTypeRoot = "root"
const NodeTypePrev = "prev"
const NodeTypeCategory = "category"
const NodeTypeUp = ".."

type CategoryPath []CategoryPathNode

type CategoryPathNode struct {
	nodeType string
	value string
	variables []string
}

func NewCategoryPathNode(nodeType string, value string, variables []string) CategoryPathNode {
	return CategoryPathNode{
		nodeType:  nodeType,
		value:     value,
		variables: variables,
	}
}