package listen_groups

import (
	"context"
)

// SubscriptionRoutine is used to subscribe to a topic and pass along messages.
// It is only called on the first request for a subscription, and is told finish when the last
// listener indicates that they are finished.
type SubscriptionRoutine = func(topic string, doneChannel ListenGroupDestroyedChan, messages chan<- ForwardedMessage, subscriptionDone chan<- bool, ctx context.Context)

// BaseForwarder handles the coordination between listen groups and a physical subscription. It has logic for only
// making one network subscription to the pubsub and sending those messages to a listen group.
//
// It also calls a [SubscriptionRoutine] if it creates a new listen group, which establishes the physical connection
// between the pubsub provider and the listen group. In this case, it is the leader of the listen group.
func BaseForwarder(ctx context.Context, topic string, lgh ListenGroupHandler, listenerId string, subscriptionRoutine SubscriptionRoutine) (chan ForwardedMessage, error) {
	incomingMessages, doneChannel, firstSubscriptionToTopic := GetListenGroup(lgh, listenerId, topic)

	// Stop listening when done
	// Note that if we are the last listener, then a message is sent to the doneChannel
	go LeaveListenGroupOnCtxDone(ctx, lgh, listenerId, topic)

	if firstSubscriptionToTopic {
		subscriptionCtx := context.Background()
		messagesFromSubscription, subscriptionDone := make(chan ForwardedMessage, 8), make(chan bool, 1)

		// Handle the actual network subscription + forward incoming messages to everyone in the group
		go subscriptionRoutine(topic, doneChannel, messagesFromSubscription, subscriptionDone, subscriptionCtx)
		go MultiplexMessages(subscriptionCtx, lgh, topic, messagesFromSubscription, subscriptionDone)
	}

	return incomingMessages, nil
}
