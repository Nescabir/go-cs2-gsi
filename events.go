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

// eventHandler represents a function that handles a specific event type
type eventHandler[T any] func(event Event[T])

// eventHandlers stores handlers for each event type
var eventHandlers = make(map[string][]interface{})
var handlersMutex sync.RWMutex

// eventName is a type-safe wrapper for event names
type eventName[T any] string

// Event names with their associated types
var (
	Data     eventName[*models.State] = eventName[*models.State](string(models.Data))
	RoundEnd eventName[*models.Score] = eventName[*models.Score](string(models.RoundEnd))
	// Kill              eventName[*models.KillEvent] = eventName[*models.KillEvent](string(models.Kill))
	// Hurt              eventName[*models.HurtEvent] = eventName[*models.HurtEvent](string(models.Hurt))
	TimeoutStart      eventName[*models.Team]   = eventName[*models.Team](string(models.TimeoutStart))
	TimeoutEnd        eventName[*models.Team]   = eventName[*models.Team](string(models.TimeoutEnd))
	Mvp               eventName[*models.Player] = eventName[*models.Player](string(models.Mvp))
	FreezetimeStart   eventName[*models.Player] = eventName[*models.Player](string(models.FreezetimeStart))
	FreezetimeEnd     eventName[*models.Player] = eventName[*models.Player](string(models.FreezetimeEnd))
	IntermissionStart eventName[*models.Player] = eventName[*models.Player](string(models.IntermissionStart))
	IntermissionEnd   eventName[*models.Player] = eventName[*models.Player](string(models.IntermissionEnd))
	DefuseStart       eventName[*models.Player] = eventName[*models.Player](string(models.DefuseStart))
	DefuseEnd         eventName[*models.Player] = eventName[*models.Player](string(models.DefuseEnd))
	BombPlantStart    eventName[*models.Player] = eventName[*models.Player](string(models.BombPlantStart))
	BombPlantStop     eventName[*models.Player] = eventName[*models.Player](string(models.BombPlantStop))
	BombPlanted       eventName[*models.Player] = eventName[*models.Player](string(models.BombPlanted))
	BombDefused       eventName[*models.Player] = eventName[*models.Player](string(models.BombDefused))
	BombExploded      eventName[*models.Player] = eventName[*models.Player](string(models.BombExploded))
	// MapEnd            eventName[*models.Score]     = eventName[*models.Score](string(models.MapEnd))
	// MapStart          eventName[*models.Score]     = eventName[*models.Score](string(models.MapStart))
	MatchEnd eventName[*models.Score] = eventName[*models.Score](string(models.MatchEnd))
)

// Subscribe registers a handler for a specific event type
// The type parameter T is automatically inferred from the event name
func Subscribe[T any](eventName eventName[T], handler eventHandler[T]) {
	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	eventHandlers[string(eventName)] = append(eventHandlers[string(eventName)], handler)
}

// publish sends an event to all registered handlers
func publish[T any](event Event[T]) {
	handlersMutex.RLock()
	handlers, found := eventHandlers[event.Name]
	handlersMutex.RUnlock()

	if !found {
		return
	}

	for _, handler := range handlers {
		// Type assertion to call the handler with the correct type
		if typedHandler, ok := handler.(eventHandler[T]); ok {
			typedHandler(event)
		}
	}
}

// Helper functions for type-safe event publishing
func publishData(data *models.State) {
	publish(Event[*models.State]{
		Name: string(models.Data),
		Data: data,
	})
}

func publishRoundEnd(data *models.Score) {
	publish(Event[*models.Score]{
		Name: string(models.RoundEnd),
		Data: data,
	})
}

// func publishKill(data *models.KillEvent) {
// 	publish(Event[*models.KillEvent]{
// 		Name: string(models.Kill),
// 		Data: data,
// 	})
// }

// func publishHurt(data *models.HurtEvent) {
// 	publish(Event[*models.HurtEvent]{
// 		Name: string(models.Hurt),
// 		Data: data,
// 	})
// }

func publishTimeoutStart(data *models.Team) {
	publish(Event[*models.Team]{
		Name: string(models.TimeoutStart),
		Data: data,
	})
}

func publishTimeoutEnd(data *models.Team) {
	publish(Event[*models.Team]{
		Name: string(models.TimeoutEnd),
		Data: data,
	})
}

func publishMvp(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.Mvp),
		Data: data,
	})
}

func publishFreezetimeStart(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.FreezetimeStart),
		Data: data,
	})
}

func publishFreezetimeEnd(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.FreezetimeEnd),
		Data: data,
	})
}

func publishIntermissionStart(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.IntermissionStart),
		Data: data,
	})
}

func publishIntermissionEnd(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.IntermissionEnd),
		Data: data,
	})
}

func publishDefuseStart(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.DefuseStart),
		Data: data,
	})
}

func publishDefuseEnd(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.DefuseEnd),
		Data: data,
	})
}

func publishBombPlantStart(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.BombPlantStart),
		Data: data,
	})
}

func publishBombPlantStop(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.BombPlantStop),
		Data: data,
	})
}

func publishBombPlanted(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.BombPlanted),
		Data: data,
	})
}

func publishBombDefused(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.BombDefused),
		Data: data,
	})
}

func publishBombExploded(data *models.Player) {
	publish(Event[*models.Player]{
		Name: string(models.BombExploded),
		Data: data,
	})
}

// func publishMapEnd(data *models.Score) {
// 	publish(Event[*models.Score]{
// 		Name: string(models.MapEnd),
// 		Data: data,
// 	})
// }

// func publishMapStart(data *models.Score) {
// 	publish(Event[*models.Score]{
// 		Name: string(models.MapStart),
// 		Data: data,
// 	})
// }

func publishMatchEnd(data *models.Score) {
	publish(Event[*models.Score]{
		Name: string(models.MatchEnd),
		Data: data,
	})
}
