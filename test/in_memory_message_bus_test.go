// Test in memory message queue implementation tests
package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-common/messaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// region Domain model for the Test ------------------------------------------------------------------------------------

type HeroMessage struct {
	BaseMessage
	Hero *Hero `json:"hero"`
}

func (m *HeroMessage) Payload() any { return m.Hero }

func NewHeroMessage() IMessage {
	return &HeroMessage{}
}

func newHeroMessage(topic string, hero *Hero) IMessage {
	message := &HeroMessage{
		Hero: hero,
	}
	message.MsgTopic = topic
	message.MsgOpCode = int(time.Now().Unix())
	message.MsgSessionId = entity.NanoID()
	return message
}

// endregion

func getInitializedMessageBus() (IMessageBus, error) {
	mq, err := NewInMemoryMessageBus()
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

func TestInMemoryMessageBus_Pop(t *testing.T) {
	skipCI(t)

	mq, fe := getInitializedMessageBus()
	assert.Nil(t, fe, "error initializing Message queue")

	for {
		if msg, err := mq.Pop(nil, 0, "queue_0"); err == nil {
			hero := msg.Payload().(*Hero)
			fmt.Println(msg.Topic(), hero.Id, hero.Name)
		} else {
			break
		}
	}
	fmt.Println("done")
}

func TestInMemoryMessageBus_PopWithTimeout(t *testing.T) {
	skipCI(t)

	mq, fe := getInitializedMessageBus()
	assert.Nil(t, fe, "error initializing Message queue")

	// Push message to queue_y after 10 seconds
	go func() {
		time.Sleep(time.Second * 5)
		_ = mq.Push(newHeroMessage("queue_x", &Hero{
			BaseEntity: entity.BaseEntity{},
			Key:        100,
			Name:       "Delayed hero",
		}))
	}()

	if msg, err := mq.Pop(nil, time.Second*12, "queue_x", "queue_y", "queue_z"); err != nil {
		fmt.Println(err.Error())
	} else {
		hero := msg.Payload().(*Hero)
		fmt.Println(msg.Topic(), hero.Id, hero.Name)
	}

	fmt.Println("done")
}

func TestInMemoryMessageBus_PubSub(t *testing.T) {
	skipCI(t)

	bus, fe := getInitializedMessageBus()
	assert.Nil(t, fe, "error initializing Message queue")

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Publish messages for 1 minute on 2 topics (heroes_1, heroes_2)
	go publishMessages(wg, bus, "heroes_1", time.Second)
	go publishMessages(wg, bus, "heroes_2", time.Second)

	sub, err := bus.Subscribe("subscriber", NewHeroMessage, subscriberCallback, "heroes_1")
	require.NoError(t, err, "subscription error")
	require.NotEmptyf(t, sub, "subscription is empty")

	wg.Wait()
	fmt.Println("done")
}

// run publisher for limited time and publish message every minute
func publishMessages(wg *sync.WaitGroup, bus IMessageBus, topic string, timeout time.Duration) {

	// run publisher for timeout
	after := time.After(timeout)
	idx := 0
	for {
		select {
		case _ = <-time.Tick(time.Millisecond):
			hero := list_of_heroes[idx]
			//message := newHeroMessage(topic, hero.(*Hero))
			message := GetMessage[*Hero](topic, hero.(*Hero))
			_ = bus.Publish(message)
			if idx == len(list_of_heroes)-1 {
				idx = 0
			} else {
				idx += 1
			}

		case <-after:
			wg.Done()
			return
		}
	}
}

// subscriber function callback
func subscriberCallback(msg IMessage) bool {
	if msg == nil {
		return false
	}
	if msg.Payload() == nil {
		return false
	}

	hero := msg.Payload().(*Hero)
	if hero == nil {
		return false
	}
	fmt.Println(msg.Topic(), msg.OpCode(), msg.SessionId(), hero.Id, hero.Name)
	return true
}
