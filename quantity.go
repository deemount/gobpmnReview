package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// quantity holds all the quantities of the BPMN elements
// in the BPMN model. It is used to count the number of elements
type quantity struct {
	Pool            int
	Process         int
	Participant     int
	ProcessElements map[int]map[string]int
}

// countPool counts a pool in the process model.
func (q *quantity) countPool(field string) {
	if strings.Contains(field, "Pool") {
		q.Pool++
	}
}

// countFieldsInPool counts the number of participants in the BPMN model, which are defined in a Pool.
// A pool is structured and have configurations, a collaboration, processes and ID and can have messages.
// Ruleset:
//   - If the reflection field contains the word "Process", then count a process.
//   - If the reflection field contains the word "Participant", then count a participant.
func (q *quantity) countFieldsInPool(v *reflectValue) {
	typ := reflect.TypeOf(v.Pool.Interface())
	if typ.Kind() != reflect.Struct {
		panic("Input data must be a struct")
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if strings.Contains(field.Name, "Process") {
			v.ProcessName = append(v.ProcessName, extractPrefixBeforeProcess(field.Name))
			q.Process++
		}
		if strings.Contains(field.Name, "Participant") {
			v.ParticipantName = append(v.ParticipantName, field.Name)
			q.Participant++
		}
	}
}

// countFieldsInProcess counts the fields in a process and stores them in a map.
func (q *quantity) countFieldsInProcess(v *reflectValue) {

	// Initialize the ProcessElements map if nil
	if q.ProcessElements == nil {
		q.ProcessElements = make(map[int]map[string]int)
	}

	// Define element types with more specific matching rules
	elements := map[string]struct {
		name       string
		exactMatch bool
	}{
		"StartEvent":             {"StartEvent", true},
		"EndEvent":               {"EndEvent", true},
		"IntermediateCatchEvent": {"IntermediateCatchEvent", true},
		"IntermediateThrowEvent": {"IntermediateThrowEvent", true},
		"ParallelGateway":        {"ParallelGateway", true},
		"ExclusiveGateway":       {"ExclusiveGateway", true},
		"InclusiveGateway":       {"InclusiveGateway", true},
		"UserTask":               {"UserTask", true},
		"ScriptTask":             {"ScriptTask", true},
		"Task":                   {"Task", false}, // Will only match if no other task type matches
		"From":                   {"SequenceFlow", false},
	}

	// process each process name
	for processIdx, processName := range v.ProcessName {
		// Initialize the inner map for this process
		if q.ProcessElements[processIdx] == nil {
			q.ProcessElements[processIdx] = make(map[string]int)
		}

		// check for multiple processes
		if q.Process > 1 {

			// handle multiple processes
			if err := q.countMultipleProcessElements(v, processIdx, processName, elements); err != nil {
				log.Printf("Error processing multiple processes: %v", err)
				continue
			}

		} else {

			// handle single process
			q.countSingleProcessElements(v, processIdx, elements)

		}

	}

}

// countMultipleProcessElements handles counting for multiple processes
func (q *quantity) countMultipleProcessElements(v *reflectValue, processIndex int, procName string, elements map[string]struct {
	name       string
	exactMatch bool
}) error {
	field := v.Target.FieldByName(procName)
	if !field.IsValid() {
		return fmt.Errorf("invalid field for process %s", procName)
	}

	return q.countFieldElements(field, processIndex, elements)
}

// countSingleProcessElements handles counting for a single process
func (q *quantity) countSingleProcessElements(v *reflectValue, processIndex int, elements map[string]struct {
	name       string
	exactMatch bool
}) {
	for j := 0; j < len(v.Fields); j++ {
		fieldName := v.Fields[j].Name

		// Handle sequence flows first
		if strings.HasPrefix(fieldName, "From") {
			q.ProcessElements[processIndex]["SequenceFlow"]++
			continue
		}

		q.matchAndCountElement(processIndex, fieldName, elements)
	}
}

// countFieldElements counts elements in a reflected struct field
func (q *quantity) countFieldElements(field reflect.Value, processIndex int, elements map[string]struct {
	name       string
	exactMatch bool
}) error {
	for j := 0; j < field.NumField(); j++ {
		fieldName := field.Type().Field(j).Name

		// handle sequence flows first
		if strings.HasPrefix(fieldName, "From") {
			q.ProcessElements[processIndex]["SequenceFlow"]++
			continue
		}

		q.matchAndCountElement(processIndex, fieldName, elements)
	}
	return nil
}

// matchAndCountElement matches and counts a single element.
// Note: seperated exact and partial matching
func (q *quantity) matchAndCountElement(processIndex int, fieldName string, elements map[string]struct {
	name       string
	exactMatch bool
}) {
	// First try to find an exact match
	for element, info := range elements {
		if info.exactMatch && fieldName == element {
			q.ProcessElements[processIndex][info.name]++
			return // Found an exact match, no need to continue; early return
		}
	}

	// If no exact match found, look for partial matches
	for element, info := range elements {
		if !info.exactMatch && strings.Contains(fieldName, element) {
			q.ProcessElements[processIndex][info.name]++
			return // Found a partial match, no need to continue; early return
		}
	}
}
