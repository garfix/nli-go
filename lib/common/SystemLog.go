package common

import (
	"fmt"
	"strings"
)

type SystemLog struct {
	debugOn     bool
	productions []string
	debugLines  []string
	debugDepth  int
	errors      []string
	clarificationQuestion string
	clarificationOptions *Options
	ok          bool
}

func NewSystemLog(debugOn bool) *SystemLog {
	log := SystemLog{debugOn: debugOn}
	log.Clear()

	return &log
}

func (log *SystemLog) Clear() {
	log.productions = []string{}
	log.debugLines = []string{}
	log.debugDepth = 0
	log.errors = []string{}
	log.clarificationQuestion = ""
	log.clarificationOptions = &Options{}
	log.ok = true
}

func (log *SystemLog) EnableDebug() {
	log.debugOn = true
}

func (log *SystemLog) DisableDebug() {
	log.debugOn = true
}

func (log *SystemLog) ToggleDebug() {
	log.debugOn = !log.debugOn
}

func (log *SystemLog) AddProduction(name string, production string) {
	log.productions = append(log.productions, name + ": " + production)
}

func (log *SystemLog) AddError(error string) {
	log.ok = false
	log.errors = append(log.errors, error)
}

func (log *SystemLog) SetClarificationRequest(question string, options *Options) {
	log.clarificationQuestion = question
	log.clarificationOptions = options
}

func (log *SystemLog) GetClarificationQuestion() string {
	return log.clarificationQuestion
}

func (log *SystemLog) GetClarificationOptions() *Options {
	return log.clarificationOptions
}

func (log *SystemLog) IsOk() bool {
	return log.ok
}

func (log *SystemLog) IsDone() bool {
	return !log.ok || log.clarificationQuestion != ""
}

func (log *SystemLog) StartDebug(text string, vals ...interface{}) {

	if !log.debugOn {
		return
	}

	stmt := strings.Repeat("  ", log.debugDepth) + text + " "
	for _, val := range vals {
		stmt += fmt.Sprintf("%v", val) + " "
	}

	log.debugLines = append(log.debugLines, stmt)
	log.debugDepth++
}

func (log *SystemLog) EndDebug(text string, vals ...interface{}) {

	if !log.debugOn {
		return
	}

	if log.debugDepth > 0 {
		log.debugDepth--
	}

	stmt := strings.Repeat("  ", log.debugDepth) + text + " "
	for _, val := range vals {
		stmt += fmt.Sprintf("%v", val) + " "
	}

	log.debugLines = append(log.debugLines, stmt)
}

func (log *SystemLog) GetDebugLines() []string {
	return log.debugLines
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

	for _, debugLine := range log.debugLines {
		s += debugLine + "\n"
	}

	return s
}
