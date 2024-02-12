package snowflake

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"
)

var (
	// onceInitGenerator guarantee initialize generator only once
	onceInitGenerator sync.Once
	// rootGenerator - the root generator
	rootGenerator *generator

	// ErrInvalidNode - invalid node id
	ErrInvalidNode = fmt.Errorf("invalid node id; must be 0 ≤ id < %d", maxNode)

	// ErrStartZero - the error of starting time is zero
	ErrStartZero = errors.New("the start time cannot be a zero value")

	// ErrStartFuture - the error of starting time is in the future
	ErrStartFuture = errors.New("the start time cannot be greater than the current millisecond")

	// ErrStartExceed - the start time is more than 69 years ago.
	ErrStartExceed = errors.New("the maximum life cycle of the snowflake algorithm is 69 years")
)

type Generator interface {
	// Next - get an unused sequence
	Next() (*big.Int, error)
}

type generator struct {
	// nodeID is the node ID that the Snowflake generator will use for the next 8 bits
	nodeID uint64
	// sequence is the last 14 bits.
	sequence chan uint64
	// baseEpoch is the start time.
	baseEpoch int64
}

func NewGenerator(node uint64, start time.Time) (Generator, error) {
	if node > maxNode {
		return nil, ErrInvalidNode
	}
	start = start.UTC()

	if start.IsZero() {
		return nil, ErrStartZero
	}

	if start.After(time.Now().UTC()) {
		return nil, ErrStartFuture
	}

	if uint64(time.Now().UnixMilli()-start.UnixMilli()) > maxEpoch {
		return nil, ErrStartExceed
	}
	// singleton
	onceInitGenerator.Do(func() {
		rootGenerator = &generator{
			nodeID:    node,
			sequence:  make(chan uint64),
			baseEpoch: start.UnixMilli(),
		}
		go func() {
			var (
				seq chan uint64
			)

			for {
				var reset <-chan time.Time
				if seq == nil {
					reset = time.After(time.Millisecond)
				}
				select {
				case <-reset:
					seq = make(chan uint64, 1)
					seq <- 0
				case current := <-seq:
					seq = nil
					for i := current; current <= maxSequence; i++ {
						rootGenerator.sequence <- i
					}
				}
			}
		}()
	})

	return rootGenerator, nil
}

func (g *generator) Next() (*big.Int, error) {
	current := time.Now().UnixMilli()
	if uint64(current-g.baseEpoch) > maxEpoch {
		return nil, ErrStartExceed
	}

	seq := <-g.sequence

	nodeId := g.nodeID << shiftNode
	result := uint64(current-g.baseEpoch)<<shiftEpoch + nodeId + seq
	num := big.NewInt(0)
	return num.SetUint64(result), nil
}
