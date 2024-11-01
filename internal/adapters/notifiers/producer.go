package notifiers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fernandoocampo/players/internal/players"
)

type EventBus interface {
	Publish(ctx context.Context, message players.NewEvent) error
}

type NotifierSetup struct {
	Logger *slog.Logger
	// timeout to publish in seconds
	TimeoutToPublish int
	EventBus         EventBus
}

// Notifier defines logic to produce events and publish those events into an event bus.
type Notifier struct {
	logger           *slog.Logger
	timeoutToPublish time.Duration
	eventBus         EventBus
	eventStream      chan players.NewEvent
}

const (
	defaultBatchEvents      = 10
	defaultTimeoutToPublish = 3 // seconds
)

func NewNotifier(setup NotifierSetup) *Notifier {
	if setup.TimeoutToPublish < 1 {
		setup.TimeoutToPublish = defaultTimeoutToPublish
	}

	newNotifier := Notifier{
		logger:           setup.Logger,
		eventStream:      make(chan players.NewEvent, defaultBatchEvents),
		eventBus:         setup.EventBus,
		timeoutToPublish: time.Duration(setup.TimeoutToPublish) * time.Second,
	}

	return &newNotifier
}

func (n *Notifier) Start(ctx context.Context) {
	n.logger.Info("starting worker as a notifier")

	go func() {
		for {
			select {
			case <-ctx.Done():
				n.logger.Info("context was cancelled")

				return
			case newEvent, ok := <-n.eventStream:
				if !ok {
					n.logger.Info("event stream was closed, ending player events worker")

					return
				}

				newCTX, cancel := context.WithTimeout(ctx, n.timeoutToPublish)
				{
					err := n.publish(newCTX, newEvent)
					if err != nil {
						n.logger.Error("publishing event", slog.String("error", err.Error()))
					}
				}
				cancel()
			}
		}
	}()
}

func (n *Notifier) publish(ctx context.Context, event players.NewEvent) error {
	n.logger.Info("publishing event into hypotetical event bus", slog.Any("event", event))

	if n.eventBus == nil {
		return nil
	}

	err := n.eventBus.Publish(ctx, event)
	if err != nil {
		n.logger.Error("publishing event into eventbus",
			slog.String("error", err.Error()),
			slog.Any("event", event))

		return fmt.Errorf("unable to publish event: %w", err)
	}

	return nil
}

func (n *Notifier) Notify(event players.NewEvent) {
	n.logger.Info("notifying new player event", slog.Any("event", event))

	go func() {
		ctx, cancel := context.WithTimeout(context.TODO(), n.timeoutToPublish)
		defer cancel()

		select {
		case <-ctx.Done():
			n.logger.Info("context was cancelled in notifying new player event")

			return
		case n.eventStream <- event:
		}
	}()
}

func (n *Notifier) Health() (string, error) {
	// take advantage of heartbeat concurrency pattern for notifier worker
	// plus checking eventbus client connection.
	return "event-notifier", nil
}
