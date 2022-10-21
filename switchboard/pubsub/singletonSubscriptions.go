package pubsub

import (
	"fmt"
)

type ListenerId = string

var topicSubscriptions = make([]*topicSubscription, 0)

type topicSubscription struct {
	Topic            string
	Subscriptions    map[ListenerId]chan Message
	NumListeners     int
	KillSubscription chan bool
}

func ListenToExistingSubscription(s ListenerId, topic string) (chan Message, chan bool, error) {
	for _, sub := range topicSubscriptions {
		if sub.Topic == topic {
			sub.NumListeners += 1
			sub.Subscriptions[s] = make(chan Message, 8)
			return sub.Subscriptions[s], sub.KillSubscription, nil
		}
	}
	return nil, nil, fmt.Errorf("no subscription currently exists for the given topic")
}

func CreateSubscriptionListener(s ListenerId, topic string) (chan Message, chan bool) {
	subs := make(map[ListenerId]chan Message)
	subs[s] = make(chan Message, 8)
	newSub := &topicSubscription{
		topic,
		subs,
		1,
		make(chan bool, 1),
	}
	topicSubscriptions = append(topicSubscriptions, newSub)
	return newSub.Subscriptions[s], newSub.KillSubscription
}

func StopListeningToExistingSubscription(s ListenerId, topic string) (int, error) {
	for i, sub := range topicSubscriptions {
		if sub.Topic == topic {
			sub.NumListeners -= 1
			delete(sub.Subscriptions, s)

			if sub.NumListeners == 0 {
				topicSubscriptions = remove(topicSubscriptions, i)
				sub.KillSubscription <- true
				return 0, nil
			} else {
				return sub.NumListeners, nil
			}

		}
	}

	return -1, fmt.Errorf("no subscription currently exists for the given topic")
}

func SendMessageToAllListeners(msg Message, topic string) {
	for _, sub := range topicSubscriptions {
		if sub.Topic == topic {
			for _, msgChan := range sub.Subscriptions {
				msgChan <- msg
			}
			return
		}
	}
}

func remove(slice []*topicSubscription, s int) []*topicSubscription {
	return append(slice[:s], slice[s+1:]...)
}
