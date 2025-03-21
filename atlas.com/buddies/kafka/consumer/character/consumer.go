package character

import (
	consumer2 "atlas-buddies/kafka/consumer"
	"atlas-buddies/list"
	"context"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("character_status_event")(EnvEventTopicCharacterStatus)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(db *gorm.DB) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(db *gorm.DB) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(rf func(topic string, handler handler.Handler) (string, error)) {
			var t string
			t, _ = topic.EnvProvider(l)(EnvEventTopicCharacterStatus)()
			_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventCreated(db))))
			_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventLogin(db))))
			_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventLogout(db))))
			_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventChannelChanged(db))))
		}
	}
}

func handleStatusEventCreated(db *gorm.DB) message.Handler[statusEvent[statusEventCreatedBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, e statusEvent[statusEventCreatedBody]) {
		if e.Type != EventCharacterStatusTypeCreated {
			return
		}
		_, _ = list.Create(l)(ctx)(db)(e.CharacterId, 30)
	}
}

func handleStatusEventLogin(db *gorm.DB) func(l logrus.FieldLogger, ctx context.Context, event statusEvent[statusEventLoginBody]) {
	return func(l logrus.FieldLogger, ctx context.Context, event statusEvent[statusEventLoginBody]) {
		if event.Type != EventCharacterStatusTypeLogin {
			return
		}
		err := list.UpdateChannel(l)(ctx)(db)(event.CharacterId, event.WorldId, int8(event.Body.ChannelId))
		if err != nil {
			l.WithError(err).Errorf("Unable to process login for character [%d].", event.CharacterId)
		}
	}
}

func handleStatusEventLogout(db *gorm.DB) func(l logrus.FieldLogger, ctx context.Context, event statusEvent[statusEventLogoutBody]) {
	return func(l logrus.FieldLogger, ctx context.Context, event statusEvent[statusEventLogoutBody]) {
		if event.Type != EventCharacterStatusTypeLogout {
			return
		}
		err := list.UpdateChannel(l)(ctx)(db)(event.CharacterId, event.WorldId, -1)
		if err != nil {
			l.WithError(err).Errorf("Unable to process logout for character [%d].", event.CharacterId)
		}
	}
}

func handleStatusEventChannelChanged(db *gorm.DB) func(l logrus.FieldLogger, ctx context.Context, event statusEvent[statusEventChannelChangedBody]) {
	return func(l logrus.FieldLogger, ctx context.Context, event statusEvent[statusEventChannelChangedBody]) {
		if event.Type != EventCharacterStatusTypeChannelChanged {
			return
		}
		if event.Body.ChannelId == event.Body.OldChannelId {
			return
		}
		err := list.UpdateChannel(l)(ctx)(db)(event.CharacterId, event.WorldId, int8(event.Body.ChannelId))
		if err != nil {
			l.WithError(err).Errorf("Unable to process change channel for character [%d].", event.CharacterId)
		}
	}
}
