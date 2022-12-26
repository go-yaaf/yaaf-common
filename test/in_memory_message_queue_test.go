// Copyright 2022. Motty Cohen
//
// Test in memory message queue implementation tests
package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-common/messaging"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// region Domain model for the Test ------------------------------------------------------------------------------------

type HeroMessage struct {
	MsgQueue string `json:"queue"`
	Hero     *Hero  `json:"hero"`
}

func (m *HeroMessage) Queue() string { return m.MsgQueue }
func (m *HeroMessage) Payload() any  { return m.Hero }

func newHeroMessage(queue string, hero *Hero) IQueueMessage {
	return &HeroMessage{
		MsgQueue: queue,
		Hero:     hero,
	}
}

// endregion

// region Tests --------------------------------------------------------------------------------------------------------

func getInitializedMessageQueue() (IMessageQueue, error) {
	mq, err := NewInMemoryMessageQueue()
	if err != nil {
		return nil, err
	}

	// Push messages to 4 queues (queue_0, queue_1, queue_2, queue_3)
	for idx, hero := range list_of_heroes {
		queue := fmt.Sprintf("queue_%d", idx%4)
		if er := mq.Push(newHeroMessage(queue, hero.(*Hero))); er != nil {
			return nil, er
		}
	}
	return mq, nil
}

func TestInMemoryMessageQueue_Pop(t *testing.T) {

	mq, fe := getInitializedMessageQueue()
	assert.Nil(t, fe, "error initializing Message queue")

	for {
		if msg, err := mq.Pop(nil, 0, "queue_0"); err == nil {
			hero := msg.Payload().(*Hero)
			fmt.Println(msg.Queue(), hero.Id, hero.Name)
		} else {
			break
		}
	}
	fmt.Println("done")
}

func TestInMemoryMessageQueue_PopWithTimeout(t *testing.T) {

	mq, fe := getInitializedMessageQueue()
	assert.Nil(t, fe, "error initializing Message queue")

	// Push message to queue_y after 10 seconds
	go func() {
		time.Sleep(time.Second * 5)
		mq.Push(newHeroMessage("queue_x", &Hero{
			BaseEntity: entity.BaseEntity{},
			Key:        100,
			Name:       "Delayed hero",
		}))
	}()

	if msg, err := mq.Pop(nil, time.Second*12, "queue_x", "queue_y", "queue_z"); err != nil {
		fmt.Println(err.Error())
	} else {
		hero := msg.Payload().(*Hero)
		fmt.Println(msg.Queue(), hero.Id, hero.Name)
	}

	fmt.Println("done")
}
