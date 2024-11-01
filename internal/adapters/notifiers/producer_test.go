package notifiers_test

import (
	"context"
	"testing"
	"time"

	"github.com/fernandoocampo/players/internal/adapters/notifiers"
	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/assert"
)

func TestNotify(t *testing.T) {
	t.Parallel()
	// Given
	playerID := unittests.NewPlayerID().String()
	newEvents := []players.NewEvent{
		{
			PlayerID: playerID,
			Event:    "new event 1",
		},
		{
			PlayerID: playerID,
			Event:    "new event 2",
		},
	}

	want := map[string]players.NewEvent{
		"new event 1": {
			PlayerID: playerID,
			Event:    "new event 1",
		},
		"new event 2": {
			PlayerID: playerID,
			Event:    "new event 2",
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	eventbusMock := newEvenbusMock()
	defer close(eventbusMock.events)

	setup := notifiers.NotifierSetup{
		Logger:           unittests.NewLogger(),
		TimeoutToPublish: 5,
		EventBus:         eventbusMock,
	}

	notifier := notifiers.NewNotifier(setup)

	notifier.Start(ctx)

	// When
	for _, event := range newEvents {
		notifier.Notify(event)
	}

	for range 2 {
		select {
		case <-ctx.Done():
			t.Errorf("unexpected context cancelled: %s", ctx.Err().Error())
			t.FailNow()
		case got, ok := <-eventbusMock.events:
			if !ok {
				t.Errorf("event bus channel was closed unexpectedly")
				t.FailNow()
			}

			wantValue, ok := want[got.Event]
			assert.True(t, ok)
			assert.Equal(t, wantValue, got)
		}
	}
}

type evenbusMock struct {
	events chan players.NewEvent
}

func newEvenbusMock() *evenbusMock {
	return &evenbusMock{
		events: make(chan players.NewEvent),
	}
}

func (e *evenbusMock) Publish(ctx context.Context, message players.NewEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e.events <- message:
	}
	return nil
}
