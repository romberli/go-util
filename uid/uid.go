package uid

import (
	"hash/crc32"
	"math/rand"
	"net"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/romberli/go-util/constant"
)

const (
	defaultNodeAndStepBits = 22
)

var (
	defaultMaxNodeID = 1023
)

// SetEpoch sets the epoch of uid generator, remember to set it before calling NewNode() function
func SetEpoch(t time.Time) {
	snowflake.Epoch = t.Unix()
}

// SetNodeBits sets the node bits, remember to set it before calling NewNode() function
func SetNodeBit(nodeBits uint8) {
	snowflake.NodeBits = nodeBits
	snowflake.StepBits = defaultNodeAndStepBits - snowflake.NodeBits
	defaultMaxNodeID = -1 ^ (-1 << snowflake.NodeBits)
}

// SetStepBits sets the step bits, remember to set it before calling NewNode() function
func SetStepBit(stepBits uint8) {
	snowflake.StepBits = stepBits
	snowflake.NodeBits = defaultNodeAndStepBits - snowflake.StepBits
	defaultMaxNodeID = -1 ^ (-1 << snowflake.NodeBits)
}

// GetLocalIP gets the local ip, it will get rid of the loop back ip
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return constant.EmptyString
	}

	for _, address := range addrs {
		in, ok := address.(*net.IPNet)
		if ok && !in.IP.IsLoopback() {
			if in.IP.To4() == nil {
				continue
			}

			return in.IP.String()
		}
	}

	return constant.EmptyString
}

// GetIPWorkerID gets the worker id with ip
func GetIPWorkerID(ip string) int {
	return HashWorkerID([]byte(ip))
}

// GetRandWorkerID gets the worker id with random number
func GetRandWorkerID() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	return int(rand.Uint32()) & defaultMaxNodeID
}

// HashWorkerID gets the worker id with input data
func HashWorkerID(data []byte) int {
	h := crc32.NewIEEE()
	_, _ = h.Write(data)

	return int(h.Sum32()) & defaultMaxNodeID
}
