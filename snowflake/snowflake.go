package snowflake

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	/*
		+--------------------------------------------------------------------------+
		| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
		+--------------------------------------------------------------------------+
	*/

	// Epoch is starting timestamp of generated ids
	// and it is set to the twitter snowflake epoch of
	// Nov 04 2010 01:42:54 UTC in milliseconds
	Epoch     int64 = 1288834974657
	EpochUnix       = time.Unix(Epoch/1000, (Epoch%1000)*1000*1000)

	NodeIDBits uint8 = 10
	SeqIDBits  uint8 = 12

	seqIDMask      int64 = -1 ^ (-1 << SeqIDBits)
	nodeIDShift          = SeqIDBits
	maxNodeID      int64 = -1 ^ (-1 << NodeIDBits)
	nodeIDMask           = maxNodeID << nodeIDShift
	timestampShift       = NodeIDBits + SeqIDBits
)

type Host struct {
	mu sync.Mutex

	epoch     time.Time
	timestamp int64
	nodeID    int64
	seqID     int64
}

// NewHost can create a new instance of id generator
func NewHost(nodeID int64) (*Host, error) {
	if NodeIDBits+SeqIDBits > 22 {
		return nil, fmt.Errorf("NodeIDBits and SeqIDBits overflow")
	}

	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("invalid nodeID, need [0, %v]", maxNodeID)
	}

	host := &Host{}

	currTime := time.Now()
	// rollback the time to Epoch to get a monotoic id generation
	// Reference from github.com/bwmarrin/snowflake
	host.epoch = currTime.Add(EpochUnix.Sub(currTime))

	host.timestamp = 0
	host.nodeID = nodeID
	host.seqID = 0
	return host, nil
}

type ID int64

// Generate can generate a new ID
func (h *Host) Generate() ID {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Since(h.epoch).Milliseconds()
	if now == h.timestamp {
		h.seqID = (h.seqID + 1) & seqIDMask

		if h.seqID == 0 {
			for now <= h.timestamp {
				now = time.Since(h.epoch).Milliseconds()
			}
			h.timestamp = now
		}
	} else {
		h.timestamp = now
		h.seqID = 0
	}

	r := ID((now)<<int64(timestampShift) |
		(h.nodeID << int64(nodeIDShift)) |
		(h.seqID),
	)

	return r
}

func (id ID) ToInt64() int64 {
	return int64(id)
}

func (id ID) ToString() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id ID) Timestamp() int64 {
	return (int64(id) >> timestampShift) + Epoch
}

func (id ID) NodeID() int64 {
	return (int64(id) & nodeIDMask) >> int64(nodeIDShift)
}

func (id ID) SeqID() int64 {
	return (int64(id) & seqIDMask)
}
