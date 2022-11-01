package listen_groups

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jettdc/switchboard/u"
	"testing"
)

func TestGetListenGroup_NoExisting(t *testing.T) {
	u.InitializeLogger("testing")
	ctrl := gomock.NewController(t)

	fmc := make(chan ForwardedMessage, 1)
	lgdc := make(ListenGroupDestroyedChan, 1)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	lgMock := NewMockListenGroupHandler(ctrl)
	lgMock.EXPECT().
		JoinListenGroup(gomock.Eq("id"), gomock.Eq("/topic")).
		Return(fmc, lgdc, fmt.Errorf("test err"))
	lgMock.EXPECT().
		CreateListenGroup(gomock.Eq("id"), gomock.Eq("/topic")).
		Return(fmc, lgdc)

	fmcRes, lgdcRes, first := GetListenGroup(lgMock, "id", "/topic")

	assert.Equal(t, fmc, fmcRes)
	assert.Equal(t, lgdc, lgdcRes)
	assert.Equal(t, first, true)
}

func TestGetListenGroup_Existing(t *testing.T) {}
