package pubsub

import "fmt"

type StdListenGroupHandler struct {
	listenGroups []*listenGroup
}

type listenGroup struct {
	Topic            string
	Listeners        map[ListenerId]chan Message
	NumListeners     int
	KillSubscription chan bool
}

func NewStdListenGroupHandler() *StdListenGroupHandler {
	return &StdListenGroupHandler{[]*listenGroup{}}
}

func (s *StdListenGroupHandler) CreateListenGroup(id ListenerId, topic string) (chan Message, chan bool) {
	if msg, kill, err := s.JoinListenGroup(id, topic); err == nil {
		return msg, kill
	}

	subs := make(map[ListenerId]chan Message)
	subs[id] = make(chan Message, 8)
	newSub := &listenGroup{
		topic,
		subs,
		1,
		make(chan bool, 1),
	}
	s.listenGroups = append(s.listenGroups, newSub)
	return newSub.Listeners[id], newSub.KillSubscription
}

func (s *StdListenGroupHandler) JoinListenGroup(id ListenerId, topic string) (chan Message, chan bool, error) {
	for _, sub := range s.listenGroups {
		if sub.Topic == topic {
			sub.NumListeners += 1
			sub.Listeners[id] = make(chan Message, 8)
			return sub.Listeners[id], sub.KillSubscription, nil
		}
	}
	return nil, nil, fmt.Errorf("no subscription currently exists for the given topic")

}

func (s *StdListenGroupHandler) LeaveListenGroup(id ListenerId, topic string) (int, error) {
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

func (s *StdListenGroupHandler) MessageGroup(msg Message, topic string) {
	for _, sub := range s.listenGroups {
		if sub.Topic == topic {
			for _, msgChan := range sub.Listeners {
				msgChan <- msg
			}
			return
		}
	}
}
