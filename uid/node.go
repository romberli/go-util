package uid

import (
	"github.com/bwmarrin/snowflake"
)

type Node struct {
	*snowflake.Node
}

// NewNode returns *Node with given node id
func NewNode(nodeID int) (*Node, error) {
	node, err := snowflake.NewNode(int64(nodeID))
	if err != nil {
		return nil, err
	}

	return &Node{node}, nil
}

// NewNodeWithDefault returns *Node with default node id
func NewNodeWithDefault() (*Node, error) {
	return NewNode(GetRandWorkerID())
}

// NewNodeWithIP returns *Node with ip
func NewNodeWithIP(ip string) (*Node, error) {
	return NewNode(GetIPWorkerID(ip))
}

// GetNodeID returns the node id of the node
func (n *Node) GetNodeID() int {
	return int(n.GenerateID().Node())
}

// Generate generates a unique id
func (n *Node) GenerateID() snowflake.ID {
	return n.Node.Generate()
}

// GenerateInt generates a unique id with int type
func (n *Node) GenerateInt() int {
	return int(n.Generate().Int64())
}

// GenerateString generates a unique id with string type
func (n *Node) GenerateString() string {
	return n.Generate().String()
}
