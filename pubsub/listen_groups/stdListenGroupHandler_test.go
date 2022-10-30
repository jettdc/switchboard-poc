package listen_groups

import (
	"github.com/go-playground/assert/v2"
	"strconv"
	"testing"
)

func TestNewStdListenGroupHandler(t *testing.T) {
	nlgh := NewStdListenGroupHandler()
	assert.Equal(t, len(nlgh.listenGroups), 0)
}

func TestStdListenGroupHandler_CreateListenGroup_NewTopic_NewListener(t *testing.T) {
	// Should create a new listen group for the given topic, with a single listener
	lg := NewStdListenGroupHandler()

	for i := 0; i < 100; i++ {
		id := "test" + strconv.Itoa(i)
		topic := "/test/topic" + strconv.Itoa(i)
		lg.CreateListenGroup(id, topic)

		// Listen group created for topic
		assert.Equal(t, lg.listenGroups[i].Topic, topic)

		// Within that listen group, there is a single listener with a message channel
		assert.Equal(t, len(lg.listenGroups[i].Listeners), 1)
		assert.Equal(t, lg.listenGroups[i].NumListeners, 1)

		_, ok := lg.listenGroups[i].Listeners[id]
		if !ok {
			t.Errorf("listener was not correctly added to the group")
		}
	}

	assert.Equal(t, len(lg.listenGroups), 100)
}

func TestStdListenGroupHandler_CreateListenGroup_ExistingTopic_NewListener(t *testing.T) {
	lg := NewStdListenGroupHandler()
	id1 := "test"
	topic := "/test/topic"

	lg.CreateListenGroup(id1, topic)

	for i := 0; i < 100; i++ {
		id2 := "test" + strconv.Itoa(i)
		lg.CreateListenGroup(id2, topic)

		// Should not create another listen group, but rather add to the group
		assert.Equal(t, len(lg.listenGroups), 1)
		assert.Equal(t, len(lg.listenGroups[0].Listeners), i+2)
		assert.Equal(t, lg.listenGroups[0].NumListeners, i+2)

		_, ok := lg.listenGroups[0].Listeners[id1]
		if !ok {
			t.Errorf("first listener was not correctly added to the group")
		}

		_, ok = lg.listenGroups[0].Listeners[id2]
		if !ok {
			t.Errorf("second listener was not correctly added to the group")
		}
	}

}

func TestStdListenGroupHandler_CreateListenGroup_ExistingTopic_ExistingListener(t *testing.T) {
	lg := NewStdListenGroupHandler()
	id1 := "test1"
	topic := "/test/topic"

	f, k := lg.CreateListenGroup(id1, topic)
	f2, k2 := lg.CreateListenGroup(id1, topic)

	// Should not create or add to the listen group
	assert.Equal(t, len(lg.listenGroups), 1)
	assert.Equal(t, len(lg.listenGroups[0].Listeners), 1)
	assert.Equal(t, lg.listenGroups[0].NumListeners, 1)

	// Should have returned the existing chans
	assert.Equal(t, &f, &f2)
	assert.Equal(t, &k, &k2)

	_, ok := lg.listenGroups[0].Listeners[id1]
	if !ok {
		t.Errorf("first listener was not correctly added to the group")
	}

}

func TestStdListenGroupHandler_CreateListenGroup_NewTopic_ExistingListener(t *testing.T) {
	lg := NewStdListenGroupHandler()
	id1 := "test1"

	for i := 0; i < 100; i++ {
		topic := "/test/topic" + strconv.Itoa(i)
		lg.CreateListenGroup(id1, topic)

		assert.Equal(t, len(lg.listenGroups), i+1)

		lisGroup := lg.listenGroups[i]

		assert.Equal(t, len(lisGroup.Listeners), 1)
		assert.Equal(t, lisGroup.NumListeners, 1)

		for listener, _ := range lisGroup.Listeners {
			assert.Equal(t, listener, id1)
		}

	}
}

func TestStdListenGroupHandler_JoinListenGroup_GroupDoesntExist(t *testing.T) {
	lg := NewStdListenGroupHandler()

	_, _, err := lg.JoinListenGroup("test", "/test/nil")
	if err == nil {
		t.Errorf("listener should not be able to join a listen group that doesn't exist")
	}
}

func TestStdListenGroupHandler_JoinListenGroup_(t *testing.T) {
	lg := NewStdListenGroupHandler()
	fmc1, dc1 := lg.CreateListenGroup("test", "/test")
	fmc2, dc2, err := lg.JoinListenGroup("test2", "/test")

	if err != nil {
		t.Errorf("joining an existing listen group should not throw an error")
	}

	assert.Equal(t, len(lg.listenGroups), 1)
	assert.Equal(t, len(lg.listenGroups[0].Listeners), 2)
	assert.Equal(t, lg.listenGroups[0].NumListeners, 2)
	assert.NotEqual(t, &fmc1, &fmc2) // Should get different message channels
	assert.Equal(t, &dc1, &dc2)      // Should get the same kill channel

}

func TestStdListenGroupHandler_LeaveListenGroup_GroupDoesntExist(t *testing.T) {
	lg := NewStdListenGroupHandler()
	_, err := lg.LeaveListenGroup("test", "/test")
	if err == nil {
		t.Errorf("should return an error if trying to leave a listen group that doesnt exist")
	}
}

func TestStdListenGroupHandler_LeaveListenGroup_MoreThanOneListener(t *testing.T) {
	lg := NewStdListenGroupHandler()
	lg.CreateListenGroup("test1", "/test")
	lg.JoinListenGroup("test2", "/test")

	left, err := lg.LeaveListenGroup("test1", "/test")
	if err != nil {
		t.Errorf("should be able to leave listen group without error")
	}

	assert.Equal(t, left, 1)
}
func TestStdListenGroupHandler_LeaveListenGroup_LastListener(t *testing.T) {
	lg := NewStdListenGroupHandler()
	lg.CreateListenGroup("test1", "/test")
	_, destroyed, _ := lg.JoinListenGroup("test2", "/test")

	go func() {
		select {
		case <-destroyed:
			return
		default:
			t.Errorf("should have messaged the destroyed channel")
			return
		}
	}()

	lg.LeaveListenGroup("test1", "/test")

	left, err := lg.LeaveListenGroup("test2", "/test")
	if err != nil {
		t.Errorf("should be able to leave listen group without error")
	}

	assert.Equal(t, left, 0)
}
