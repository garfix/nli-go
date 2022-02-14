package central

import "nli-go/lib/mentalese"

type MessageListener func(message mentalese.RelationSet)

type MessageManager struct {
	listeners []MessageListener
	messages  []mentalese.RelationSet
}

func NewMessageManager() *MessageManager {
	return &MessageManager{
		// note! the system actually supports only a single listener
		listeners: []MessageListener{},
	}
}

// returns any pending messages
func (mm *MessageManager) AddListener(l MessageListener) []mentalese.RelationSet {
	mm.listeners = append(mm.listeners, l)

	messages := mm.messages

	mm.messages = []mentalese.RelationSet{}

	return messages
}

func (mm *MessageManager) RemoveListener(l MessageListener) {
	mm.listeners = []MessageListener{}
}

func (mm *MessageManager) NotifyListeners(message mentalese.RelationSet) {

	if len(mm.listeners) == 0 {
		mm.messages = append(mm.messages, message)
		return
	}

	for _, l := range mm.listeners {
		l(message)
	}
}
