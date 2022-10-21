package pubsub

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jettdc/switchboard/u"
)

// SubscriptionRoutine is used to subscribe to a topic and pass along messages
// It is only called on the first request for a subscription, and is told finish when the last
// listener indicates that they are finished.
type SubscriptionRoutine = func(topic string, doneChannel <-chan bool, messages chan<- Message, subscriptionDone chan<- bool, ctx context.Context)

// Handles the logic for only making one network subscription to the pubsub, but multiplexing all the messages
// from that subscription to all listeners.
func baseSubscribe(ctx context.Context, topic string, subscriptionRoutine SubscriptionRoutine) (chan Message, error) {
	// If a subscription already exists, listen, otherwise make one
	listenerId := uuid.NewString()

	var messages chan Message
	var doneChannel chan bool
	firstSubscriptionToTopic := false

	// Someone has already initiated a subscription, so just listen to it
	if existingMessageChannel, killSubscription, err := ListenToExistingSubscription(listenerId, topic); err == nil {
		u.Logger.Debug(fmt.Sprintf("Listening on existing topic subscription to %s", topic))
		messages = existingMessageChannel
		doneChannel = killSubscription

		// We are the first to subscribe, so create a new subscription
	} else {
		u.Logger.Debug(fmt.Sprintf("Subscribing to topic %s", topic))
		messages, doneChannel = CreateSubscriptionListener(listenerId, topic)
		firstSubscriptionToTopic = true
	}

	// Stop listening when done
	// Note that if we are the last listener, then a message is sent to the doneChannel
	go func() {
		for {
			select {
			case <-ctx.Done():
				numListening, _ := StopListeningToExistingSubscription(listenerId, topic)
				u.Logger.Debug(fmt.Sprintf("No longer listening on topic %s. %d listeners remain.", topic, numListening))
				return
			}
		}
	}()

	// Only subscribe if no existing subscriptions
	if firstSubscriptionToTopic {

		messagesFromSubscription := make(chan Message, 8)
		subscriptionDone := make(chan bool, 1)
		ctx := context.Background()

		// Handle the actual network subscription
		go subscriptionRoutine(topic, doneChannel, messagesFromSubscription, subscriptionDone, ctx)

		// Multiplex messages
		go func() {
			for {
				select {
				case msg := <-messagesFromSubscription:
					SendMessageToAllListeners(msg, topic)
				case <-subscriptionDone:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return messages, nil
}
