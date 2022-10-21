package pubsub

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jettdc/switchboard/u"
)

type ListenerId = string

type ListenGroupHandler interface {
	CreateListenGroup(s ListenerId, topic string) (chan Message, chan bool)
	JoinListenGroup(s ListenerId, topic string) (chan Message, chan bool, error)
	LeaveListenGroup(s ListenerId, topic string) (int, error)
	MessageGroup(msg Message, topic string)
}

func NewListenerId() string {
	return uuid.NewString()
}

func GetListenGroup(lgh ListenGroupHandler, id ListenerId, topic string) (chan Message, chan bool, bool) {
	messages := make(chan Message, 8)
	doneChannel := make(chan bool, 1)
	firstSubscriptionToTopic := false
	// Someone has already initiated a listen group, so join it
	if existingMessageChannel, killSubscription, err := lgh.JoinListenGroup(id, topic); err == nil {
		u.Logger.Debug(fmt.Sprintf("Joining listen group for %s", topic))
		messages = existingMessageChannel
		doneChannel = killSubscription

		// We are the first to subscribe, so create a new group
	} else {
		u.Logger.Debug(fmt.Sprintf("Creating listen group for %s", topic))
		messages, doneChannel = lgh.CreateListenGroup(id, topic)
		firstSubscriptionToTopic = true
	}

	return messages, doneChannel, firstSubscriptionToTopic
}

func LeaveListenGroupOnCtxDone(ctx context.Context, lgh ListenGroupHandler, id ListenerId, topic string) {
	for {
		select {
		case <-ctx.Done():
			numListening, _ := lgh.LeaveListenGroup(id, topic)
			u.Logger.Debug(fmt.Sprintf("No longer listening on topic %s. %d listeners remain.", topic, numListening))
			return
		}
	}
}

func MultiplexMessages(ctx context.Context, lgh ListenGroupHandler, topic string, messages <-chan Message, subscriptionDone <-chan bool) {
	for {
		select {
		case msg := <-messages:
			lgh.MessageGroup(msg, topic)
		case <-subscriptionDone:
			return
		case <-ctx.Done():
			return
		}
	}
}

func removeListenerFromGroup(slice []*listenGroup, s int) []*listenGroup {
	return append(slice[:s], slice[s+1:]...)
}
