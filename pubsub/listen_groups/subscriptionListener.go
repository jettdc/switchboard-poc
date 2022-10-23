/*
Package listen_groups provides functionality for administering "listen groups" used primarily by pubsub implementations.

# Listen Groups

A "listen group" is a concept for managing multiple processes that want to "listen" to a certain topic.
Because there may be many processes that want to listen to a given topic, and we don't need more than one
actual stream of messages from the pubsub provider, we can provide many "logical" subscriptions corresponding
to a single stream of messages from a "physical" subscription.

Note that a listen group does not actually deal with the logic of subscribing through a pubsub provider, but
rather provides a simple abstraction for the single physical subscriber to forward the messages to all of
the logical subscribers.

That single physical subscriber will be henceforth referred to as the "leader" of the listen group.

Essentially, a listen group is just an association between a topic name and a list of channels that processes
will receive messages from that topic on.
*/
package listen_groups

import (
	"context"
	"fmt"
	"github.com/jettdc/switchboard/u"
)

// ListenerId uniquely identifies a process that is listening in a group
type ListenerId = string

// ListenGroupDestroyedChan receives a message when the final listener is removed from a group, indicating that any
// associated subscriptions to pubsub topics should be destroyed.
type ListenGroupDestroyedChan = <-chan bool

// ListenGroupHandler can be implemented with different strategies for administrating listen groups.
//
// # Methods
//
//	CreateListenGroup(s ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan)
//
// Should add a new listen group to the handler, associated with the topic and with a single initial
// user denoted by the ListenerId
//
// Note that if a user tries to create a listen group and one already exists for that topic, the listener should just
// join the existing listen group
//
//	JoinListenGroup(s ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan, error)
//
// Should add the listener to the group for the topic, as well as return the message channel and a
// [ListenGroupDestroyedChan] channel.
//
// An error should be indicated if no listen group currently exists for the topic.
//
//	LeaveListenGroup(s ListenerId, topic string) (int, error)
//
// Should remove the listener from the group, as well as send a cancel request if there are no more listeners in the group.
//
//	MessageGroup(msg ForwardedMessage, topic string)
//
// Should send a message to the channel for each listener in the group
type ListenGroupHandler interface {
	CreateListenGroup(s ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan)
	JoinListenGroup(s ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan, error)
	LeaveListenGroup(s ListenerId, topic string) (int, error)
	MessageGroup(msg ForwardedMessage, topic string)
}

// GetListenGroup either joins a listen group or creates a new one if none exist, as well as returns the relevant
// listen group channels and a boolean indicating whether they were the first listener in the group (so that the client
// can then establish itself as the leader and make the physical pubsub subscription)
func GetListenGroup(lgh ListenGroupHandler, id ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan, bool) {
	var messages chan ForwardedMessage
	var doneChannel ListenGroupDestroyedChan

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

// LeaveListenGroupOnCtxDone can be used to clean up resources when the context finishes.
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

// MultiplexMessages can be called by the leader to forward messages from the pubsub subscription to the listen group.
func MultiplexMessages(ctx context.Context, lgh ListenGroupHandler, topic string, messages <-chan ForwardedMessage, subscriptionDone <-chan bool) {
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
