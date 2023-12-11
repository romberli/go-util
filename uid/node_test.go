package uid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testNode *Node
)

func init() {
	ip := GetLocalIP()
	node, err := NewNodeWithIP(ip)
	if err != nil {
		panic(err)
	}

	testNode = node
}

func TestNode_All(t *testing.T) {
	TestNode_GenerateID(t)
	TestNode_GenerateInt(t)
	TestNode_GenerateString(t)
}

func TestNode_GenerateID(t *testing.T) {
	asst := assert.New(t)

	id1 := testNode.GenerateID()
	asst.NotNil(id1, "test Node.GenerateID() failed")
	id2 := testNode.GenerateID()
	asst.NotNil(id2, "test Node.GenerateID() failed")
	asst.NotEqual(id1, id2, "test Node.GenerateID() failed")
	t.Logf("id1: %d, id2: %d", id1, id2)
}

func TestNode_GenerateInt(t *testing.T) {
	asst := assert.New(t)

	id1 := testNode.GenerateID()
	asst.NotNil(id1, "test Node.GenerateID() failed")
	id2 := testNode.GenerateID()
	asst.NotNil(id2, "test Node.GenerateID() failed")
	asst.NotEqual(id1, id2, "test Node.GenerateID() failed")
	t.Logf("id1: %d, id2: %d", id1, id2)
}

func TestNode_GenerateString(t *testing.T) {
	asst := assert.New(t)

	id1 := testNode.GenerateID()
	asst.NotNil(id1, "test Node.GenerateID() failed")
	id2 := testNode.GenerateID()
	asst.NotNil(id2, "test Node.GenerateID() failed")
	asst.NotEqual(id1, id2, "test Node.GenerateID() failed")
	t.Logf("id1: %d, id2: %d", id1, id2)
}
