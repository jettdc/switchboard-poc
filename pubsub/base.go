package pubsub

import (
	"context"
	"github.com/jettdc/switchboard/pubsub/listen_groups"
)

// SubscriptionRoutine is used to subscribe to a topic and pass along messages
// It is only called on the first request for a subscription, and is told finish when the last
// listener indicates that they are finished.
type SubscriptionRoutine = func(topic string, doneChannel <-chan bool, messages chan<- listen_groups.ForwardedMessage, subscriptionDone chan<- bool, ctx context.Context)

// Handles the logic for only making one network subscription to the pubsub, but multiplexing all the messages
// from that subscription to all listeners.
func baseSubscribe(ctx context.Context, topic string, lgh listen_groups.ListenGroupHandler, listenerId string, subscriptionRoutine SubscriptionRoutine) (chan listen_groups.ForwardedMessage, error) {
	incomingMessages, doneChannel, firstSubscriptionToTopic := listen_groups.GetListenGroup(lgh, listenerId, topic)

	// Stop listening when done
	// Note that if we are the last listener, then a message is sent to the doneChannel
	go listen_groups.LeaveListenGroupOnCtxDone(ctx, lgh, listenerId, topic)

	if firstSubscriptionToTopic {
		subscriptionCtx := context.Background()
		messagesFromSubscription, subscriptionDone := make(chan listen_groups.ForwardedMessage, 8), make(chan bool, 1)

		// Handle the actual network subscription + forward incoming messages to everyone in the group
		go subscriptionRoutine(topic, doneChannel, messagesFromSubscription, subscriptionDone, subscriptionCtx)
		go listen_groups.MultiplexMessages(subscriptionCtx, lgh, topic, messagesFromSubscription, subscriptionDone)
	}

	return incomingMessages, nil
}
