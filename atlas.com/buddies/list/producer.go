package list

import (
	list2 "atlas-buddies/kafka/message/list"
	"github.com/Chronicle20/atlas-kafka/producer"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/segmentio/kafka-go"
)

func UpdateCapacityCommandProvider(characterId uint32, worldId byte, capacity byte) model.Provider[[]kafka.Message] {
	key := producer.CreateKey(int(characterId))
	value := &list2.Command[list2.UpdateCapacityCommandBody]{
		WorldId:     worldId,
		CharacterId: characterId,
		Type:        list2.CommandTypeUpdateCapacity,
		Body: list2.UpdateCapacityCommandBody{
			Capacity: capacity,
		},
	}
	return producer.SingleMessageProvider(key, value)
}