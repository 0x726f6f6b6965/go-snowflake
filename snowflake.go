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
	// onceCloseGenerator guarantee close generator only once
	onceCloseGenerator sync.Once

	// ErrInvalidNode - invalid node id
	ErrInvalidNode = fmt.Errorf("invalid node id; must be 0 ≤ id < %d", maxNode)

	// ErrStartZero - the error of starting time is zero
	ErrStartZero = errors.New("the start time cannot be a zero value")

	// ErrStartFuture - the error of starting time is in the future
	ErrStartFuture = errors.New("the start time cannot be greater than the current millisecond")

	// ErrStartExceed - the start time is more than 69 years ago.
	ErrStartExceed = errors.New("the maximum life cycle of the snowflake algorithm is 69 years")

	// ErrGeneratorClosed
	ErrGeneratorClosed = errors.New("generator is closed")
)

type Generator interface {
	// Next - get an unused sequence
	Next() (*big.Int, error)
	// Close - close the generator
	Close()
}

type generator struct {
	// nodeID is the node ID that the Snowflake generator will use for the next 8 bits
	nodeID int64
	// sequence is the last 14 bits.
	sequence chan int64
	// baseEpoch is the start time.
	baseEpoch int64
	// stop is the signal to close the generator.
	stop chan struct{}
}

func NewGenerator(node int64, start time.Time) (Generator, error) {
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
			sequence:  make(chan int64),
			baseEpoch: start.UnixMilli(),
			stop:      make(chan struct{}, 1),
		}
		go func() {
			var (
				signal chan struct{}
			)

			for {
				var reset <-chan time.Time
				if signal == nil {
					reset = time.After(time.Millisecond)
				}
				select {
				case <-reset:
					signal = make(chan struct{}, 1)
					signal <- struct{}{}
				case <-signal:
					signal = nil
					var i int64
					for i = 0; i <= maxSequence; i++ {
						rootGenerator.sequence <- i
						select {
						case <-rootGenerator.stop:
							close(rootGenerator.sequence)
							close(rootGenerator.stop)
							return
						default:
							continue
						}
					}
				}
			}
		}()
	})

	return rootGenerator, nil
}

func (g *generator) Next() (*big.Int, error) {
	seq, ok := <-g.sequence
	if !ok {
		return nil, ErrGeneratorClosed
	}

	current := time.Now().UnixMilli()
	if uint64(current-g.baseEpoch) > maxEpoch {
		return nil, ErrStartExceed
	}

	result := (current-g.baseEpoch)<<shiftEpoch | g.nodeID<<shiftNode | seq
	num := big.NewInt(result)
	return num, nil
}

func (g *generator) Close() {
	onceCloseGenerator.Do(func() {
		g.stop <- struct{}{}
		// clean sequence channel
		<-g.sequence
	})
}
