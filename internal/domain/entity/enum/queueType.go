package enum

type QueueType string

const (
	Flex QueueType = "RANKED_FLEX_SR"
	Solo QueueType = "RANKED_SOLO_5x5"
)
