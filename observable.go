package riot

import (
	//"fmt"
	"reflect"
)

type Listener struct {
	id         int
	callback   interface{}
	event_type reflect.Type
	busy       bool
}

var _newListenerId int = 0

func newListener(callback interface{}) Listener {
	callback_type := reflect.TypeOf(callback)

	if callback_type.NumIn() != 1 {
		panic("listener call back must take precisely one argument")
	}

	event_type := callback_type.In(0)

	_newListenerId += 1

	return Listener{
		id:         _newListenerId,
		callback:   callback,
		event_type: event_type,
	}
}

type ListenerList []Listener

func (listeners ListenerList) dispatch(event interface{}) {
	for _, listener := range listeners {
		if listener.busy {
			continue // prevent recursive listener dispatch
		}

		arguments := []reflect.Value{reflect.ValueOf(event)}

		listener.busy = true
		reflect.ValueOf(listener.callback).Call(arguments)
		listener.busy = false
	}
}

type ListenerMap map[reflect.Type]ListenerList

type Sink struct {
	listeners ListenerMap
	once      ListenerMap
}

func (listeners ListenerMap) add(callback interface{}) Listener {
	listener := newListener(callback)

	if _, present := listeners[listener.event_type]; !present {
		listeners[listener.event_type] = make(ListenerList, 0, 3)
	}

	listeners[listener.event_type] = append(listeners[listener.event_type], listener)

	return listener
}

func (listeners ListenerMap) remove(listener Listener) {
	if list, present := listeners[listener.event_type]; present {
		for i, item := range list {
			if item.id == listener.id {
				// TODO this doesn't remove the item, instead it seems to add one!?!
				list = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
}

func (sink *Sink) On(callback interface{}) Listener {
	return sink.listeners.add(callback)
}

func (sink *Sink) Once(callback interface{}) Listener {
	return sink.once.add(callback)
}

func (sink *Sink) Send(event interface{}) {
	event_type := reflect.TypeOf(event)

	if listeners, present := sink.once[event_type]; present {
		listeners.dispatch(event)
		sink.once[event_type] = nil
	}

	if listeners, present := sink.listeners[event_type]; present {
		listeners.dispatch(event)
	}
}

func (sink *Sink) Off(listener Listener) {
	sink.listeners.remove(listener)
}

func NewSink() Sink {
	return Sink{
		listeners: make(ListenerMap),
		once:      make(ListenerMap),
	}
}
