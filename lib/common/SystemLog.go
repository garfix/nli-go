package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type LogMessage struct {
	MessageType string
	Message     string
}

type LogListener func(message LogMessage)

type LogNode struct {
	Label    string
	Bindings string
	Children []*LogNode
}

type SystemLog struct {
	listeners   []LogListener
	active      bool
	nodes       map[string][]*LogNode
	nodesById   map[string]map[string]*LogNode
	productions []string
	debugLines  []string
	debugDepth  int
	errors      []string
	mutex       sync.Mutex
	ok          bool
}

func NewSystemLog() *SystemLog {
	log := SystemLog{
		active:    false,
		nodes:     map[string][]*LogNode{},
		nodesById: map[string]map[string]*LogNode{},
	}
	log.Clear()

	return &log
}

func (log *SystemLog) AddListener(listener LogListener) {
	log.listeners = append(log.listeners, listener)
}

func (log *SystemLog) IsActive() bool {
	return log.active
}

func (log *SystemLog) Clear() {
	log.productions = []string{}
	log.debugLines = []string{}
	log.debugDepth = 0
	log.errors = []string{}
	log.ok = true
}

func (log *SystemLog) SetDebug(on bool) {
	log.active = on
}

func (log *SystemLog) AddFrame(message string, messageType string, action string, id string, parentId string) {

	if !log.active {
		return
	}

	log.mutex.Lock()

	_, found := log.nodes[messageType]
	if !found {
		log.nodes[messageType] = []*LogNode{}
		log.nodesById[messageType] = map[string]*LogNode{}
	}

	if action == "create" {

		node := LogNode{
			Label:    message,
			Bindings: "",
			Children: []*LogNode{},
		}

		parent, found := log.nodesById[messageType][parentId]

		if found {
			parent.Children = append(parent.Children, &node)
		} else {
			log.nodes[messageType] = append(log.nodes[messageType], &node)
		}

		log.nodesById[messageType][id] = &node

	} else if action == "append" {

		node, found := log.nodesById[messageType][id]

		if found {
			node.Bindings = message
		}

		if parentId == "root" {
			log.PushFrames(messageType, node)
		}

	}

	log.mutex.Unlock()
}

func (log *SystemLog) PushFrames(messageType string, node *LogNode) {
	for _, listener := range log.listeners {

		json, _ := json.Marshal(node)

		listener(LogMessage{
			MessageType: messageType,
			Message:     string(json),
		})
	}
}

func (log *SystemLog) AddProduction(name string, production string) {
	stmt := name + ": " + production + " "
	log.productions = append(log.productions, stmt)
	for _, listener := range log.listeners {
		listener(LogMessage{
			MessageType: name,
			Message:     production,
		})
	}
}

func (log *SystemLog) AddDebug(name string, production string) {
	stmt := strings.Repeat("| ", log.debugDepth) + name + ": " + production + " "
	log.debugLines = append(log.debugLines, stmt)
}

func (log *SystemLog) AddError(error string) {
	log.ok = false
	log.errors = append(log.errors, error)
	log.AddDebug("ERROR", error)
}

func (log *SystemLog) IsOk() bool {
	return log.ok
}

func (log *SystemLog) GetProductions() []string {
	return log.productions
}

func (log *SystemLog) GetErrors() []string {
	return log.errors
}

func (log *SystemLog) String() string {
	s := ""

	if !log.IsOk() {
		s += "\n"
		for _, error := range log.errors {
			s += "ERROR: " + error + "\n"
		}
		s += "\n"
	}

	for _, production := range log.GetProductions() {
		s += fmt.Sprintln(production)
	}
	s += "\n"

	for _, debugLine := range log.debugLines {
		s += debugLine + "\n"
	}

	return s
}
