package b_tree

type BTree struct {
	serialNumber int64
	root Node

}

type Node interface {
	Serial() int64
}

type InternalNode struct {
	serial int64
	nChild int64
	child []Node
}

func NewInternalNode(serial int64, nChild int64, child []Node) *InternalNode {
	return &InternalNode{serial: serial, nChild: nChild, child: child}
}
