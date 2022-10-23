package listen_groups

import (
	"fmt"
)

type stdListenGroupHandler struct {
	listenGroups []*listenGroup
}

type listenGroup struct {
	Topic            string
	Listeners        map[ListenerId]chan ForwardedMessage
	NumListeners     int
	KillSubscription chan bool
}

// NewStdListenGroupHandler returns the standard implementation of [ListenGroupHandler]
func NewStdListenGroupHandler() *stdListenGroupHandler {
	return &stdListenGroupHandler{[]*listenGroup{}}
}

func (s *stdListenGroupHandler) CreateListenGroup(id ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan) {
	if msg, kill, err := s.JoinListenGroup(id, topic); err == nil {
		return msg, kill
	}

	subs := make(map[ListenerId]chan ForwardedMessage)
	subs[id] = make(chan ForwardedMessage, 8)
	newSub := &listenGroup{
		topic,
		subs,
		1,
		make(chan bool, 1),
	}
	s.listenGroups = append(s.listenGroups, newSub)
	return newSub.Listeners[id], newSub.KillSubscription
}

func (s *stdListenGroupHandler) JoinListenGroup(id ListenerId, topic string) (chan ForwardedMessage, ListenGroupDestroyedChan, error) {
	for _, sub := range s.listenGroups {
		if sub.Topic == topic {
			existing, ok := sub.Listeners[id]
			if ok {
				return existing, sub.KillSubscription, nil
			}

			sub.NumListeners += 1
			sub.Listeners[id] = make(chan ForwardedMessage, 8)
			return sub.Listeners[id], sub.KillSubscription, nil
		}
	}
	return nil, nil, fmt.Errorf("no subscription currently exists for the given topic")

}

func (s *stdListenGroupHandler) LeaveListenGroup(id ListenerId, topic string) (int, error) {
	for i, sub := range s.listenGroups {
		if sub.Topic == topic {
			sub.NumListeners -= 1
			delete(sub.Listeners, id)

			if sub.NumListeners == 0 {
				s.listenGroups = removeListenerFromGroup(s.listenGroups, i)
				sub.KillSubscription <- true
				return 0, nil
			} else {
				return sub.NumListeners, nil
			}

		}
	}

	return -1, fmt.Errorf("no subscription currently exists for the given topic")
}

func (s *stdListenGroupHandler) MessageGroup(msg ForwardedMessage, topic string) {
	for _, sub := range s.listenGroups {
		if sub.Topic == topic {
			for _, msgChan := range sub.Listeners {
				msgChan <- msg
			}
			return
		}
	}
}
