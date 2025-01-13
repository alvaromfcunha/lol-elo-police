package enum

type QueueType string

const (
	Flex QueueType = "RANKED_FLEX_SR"
	Solo QueueType = "RANKED_SOLO_5x5"
)

type QueueId int

const (
	NormalId    QueueId = 400
	SoloId      QueueId = 420
	FlexId      QueueId = 440
	AramId      QueueId = 450
	SwiftPlayId QueueId = 480
	QuickPlayId QueueId = 490
)

var QueueIdTypeMap = map[QueueId]QueueType{
	FlexId: Flex,
	SoloId: Solo,
}
