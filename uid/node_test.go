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
	t.Logf("id1 base2: %s, id2 base2: %s", id1.Base2(), id2.Base2())
	t.Logf("id1 base32: %s, id2 base32: %s", id1.Base32(), id2.Base32())
	t.Logf("id1 base36: %s, id2 base2: %s", id1.Base36(), id2.Base36())
	t.Logf("id1 base58: %s, id2 base2: %s", id1.Base58(), id2.Base58())
	t.Logf("id1 base64: %s, id2 base2: %s", id1.Base64(), id2.Base64())
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
