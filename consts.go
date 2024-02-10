package snowflake

const (
	maxEpoch    uint64 = 1<<41 - 1
	maxNode     uint64 = 1<<8 - 1
	maxSequence uint64 = 1<<14 - 1
	shiftEpoch  uint8  = 22
	shiftNode   uint8  = 14
)
