package common

import (
	"fmt"
	"strings"
)

type SystemLog struct {
	debugOn     bool
	printOn		bool
	productions []string
	debugLines  []string
	debugDepth  int
	errors      []string
	clarificationQuestion string
	clarificationOptions *Options
	ok          bool
	depth		int
}

func NewSystemLog() *SystemLog {
	log := SystemLog{
		debugOn: false,
		printOn: false,
	}
	log.Clear()

	return &log
}

func (log *SystemLog) SetPrint(on bool) {
	log.printOn = on
}

func (log *SystemLog) Active() bool {
	return log.debugOn
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

func (log *SystemLog) AddProduction(name string, production string) {
	stmt := name + ": " + production + " "
	log.productions = append(log.productions, stmt)
	if log.printOn { fmt.Println(stmt) }
}

func (log *SystemLog) AddDebug(name string, production string) {
	stmt := strings.Repeat("| ", log.debugDepth) + name + ": " + production + " "
	log.debugLines = append(log.debugLines, stmt)
	if log.printOn { fmt.Println(stmt) }
}

func (log *SystemLog) StartDebug(name string, production string) bool {
	log.AddDebug("╭ " + name, production)
	log.debugDepth++
	return true
}

func (log *SystemLog) EndDebug(name string, production string) bool {
	log.debugDepth--
	log.AddDebug("╰ " + name, production)
	return true
}

func (log *SystemLog) AddError(error string) {
	log.ok = false
	log.errors = append(log.errors, error)
	log.AddDebug("ERROR", error)
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
	s += "\n"

	for _, debugLine := range log.debugLines {
		s += debugLine + "\n"
	}

	return s
}
