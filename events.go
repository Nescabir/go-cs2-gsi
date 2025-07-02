package cs2gsi

import (
	"sync"

	models "github.com/nescabir/go-cs2-gsi/models"
)

// Event represents a typed event with a name and data
type Event[T any] struct {
	Name string
	Data T
}

// EventHandler represents a function that handles a specific event type
type EventHandler[T any] func(event Event[T])

// eventHandlers stores handlers for each event type
var eventHandlers = make(map[string][]interface{})
var handlersMutex sync.RWMutex

// EventName is a type-safe wrapper for event names
type EventName[T any] string

// Event names with their associated types
var (
	Data              EventName[*models.State]     = EventName[*models.State](string(models.Data))
	RoundEnd          EventName[*models.Score]     = EventName[*models.Score](string(models.RoundEnd))
	Kill              EventName[*models.KillEvent] = EventName[*models.KillEvent](string(models.Kill))
	Hurt              EventName[*models.HurtEvent] = EventName[*models.HurtEvent](string(models.Hurt))
	TimeoutStart      EventName[*models.Team]      = EventName[*models.Team](string(models.TimeoutStart))
	TimeoutEnd        EventName[*models.Team]      = EventName[*models.Team](string(models.TimeoutEnd))
	Mvp               EventName[*models.Player]    = EventName[*models.Player](string(models.Mvp))
	FreezetimeStart   EventName[*models.Player]    = EventName[*models.Player](string(models.FreezetimeStart))
	FreezetimeEnd     EventName[*models.Player]    = EventName[*models.Player](string(models.FreezetimeEnd))
	IntermissionStart EventName[*models.Player]    = EventName[*models.Player](string(models.IntermissionStart))
	IntermissionEnd   EventName[*models.Player]    = EventName[*models.Player](string(models.IntermissionEnd))
	DefuseStart       EventName[*models.Player]    = EventName[*models.Player](string(models.DefuseStart))
	DefuseEnd         EventName[*models.Player]    = EventName[*models.Player](string(models.DefuseEnd))
	BombPlantStart    EventName[*models.Player]    = EventName[*models.Player](string(models.BombPlantStart))
	BombPlantStop     EventName[*models.Player]    = EventName[*models.Player](string(models.BombPlantStop))
	BombPlanted       EventName[*models.Player]    = EventName[*models.Player](string(models.BombPlanted))
	BombDefused       EventName[*models.Player]    = EventName[*models.Player](string(models.BombDefused))
	BombExploded      EventName[*models.Player]    = EventName[*models.Player](string(models.BombExploded))
	MapEnd            EventName[*models.Score]     = EventName[*models.Score](string(models.MapEnd))
	MapStart          EventName[*models.Score]     = EventName[*models.Score](string(models.MapStart))
	MatchEnd          EventName[*models.Score]     = EventName[*models.Score](string(models.MatchEnd))
)

// Subscribe registers a handler for a specific event type
// The type parameter T is automatically inferred from the event name
func Subscribe[T any](eventName EventName[T], handler EventHandler[T]) {
	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	eventHandlers[string(eventName)] = append(eventHandlers[string(eventName)], handler)
}

// Publish sends an event to all registered handlers
func Publish[T any](event Event[T]) {
	handlersMutex.RLock()
	handlers, found := eventHandlers[event.Name]
	handlersMutex.RUnlock()

	if !found {
		return
	}

	for _, handler := range handlers {
		// Type assertion to call the handler with the correct type
		if typedHandler, ok := handler.(EventHandler[T]); ok {
			typedHandler(event)
		}
	}
}

// Helper functions for type-safe event publishing
func PublishData(data *models.State) {
	Publish(Event[*models.State]{
		Name: string(models.Data),
		Data: data,
	})
}

func PublishRoundEnd(data *models.Score) {
	Publish(Event[*models.Score]{
		Name: string(models.RoundEnd),
		Data: data,
	})
}

func PublishKill(data *models.KillEvent) {
	Publish(Event[*models.KillEvent]{
		Name: string(models.Kill),
		Data: data,
	})
}

func PublishHurt(data *models.HurtEvent) {
	Publish(Event[*models.HurtEvent]{
		Name: string(models.Hurt),
		Data: data,
	})
}

func PublishTimeoutStart(data *models.Team) {
	Publish(Event[*models.Team]{
		Name: string(models.TimeoutStart),
		Data: data,
	})
}

func PublishTimeoutEnd(data *models.Team) {
	Publish(Event[*models.Team]{
		Name: string(models.TimeoutEnd),
		Data: data,
	})
}

func PublishMvp(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.Mvp),
		Data: data,
	})
}

func PublishFreezetimeStart(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.FreezetimeStart),
		Data: data,
	})
}

func PublishFreezetimeEnd(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.FreezetimeEnd),
		Data: data,
	})
}

func PublishIntermissionStart(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.IntermissionStart),
		Data: data,
	})
}

func PublishIntermissionEnd(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.IntermissionEnd),
		Data: data,
	})
}

func PublishDefuseStart(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.DefuseStart),
		Data: data,
	})
}

func PublishDefuseEnd(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.DefuseEnd),
		Data: data,
	})
}

func PublishBombPlantStart(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.BombPlantStart),
		Data: data,
	})
}

func PublishBombPlantStop(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.BombPlantStop),
		Data: data,
	})
}

func PublishBombPlanted(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.BombPlanted),
		Data: data,
	})
}

func PublishBombDefused(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.BombDefused),
		Data: data,
	})
}

func PublishBombExploded(data *models.Player) {
	Publish(Event[*models.Player]{
		Name: string(models.BombExploded),
		Data: data,
	})
}

func PublishMapEnd(data *models.Score) {
	Publish(Event[*models.Score]{
		Name: string(models.MapEnd),
		Data: data,
	})
}

func PublishMapStart(data *models.Score) {
	Publish(Event[*models.Score]{
		Name: string(models.MapStart),
		Data: data,
	})
}

func PublishMatchEnd(data *models.Score) {
	Publish(Event[*models.Score]{
		Name: string(models.MatchEnd),
		Data: data,
	})
}
