package main

import (
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

	// Process each process name
	for i, procName := range v.ProcessName {
		// Initialize the inner map for this process
		if q.ProcessElements[i] == nil {
			q.ProcessElements[i] = make(map[string]int)
		}

		// check for multiple processes
		if q.Process > 1 {
			// Get the fields of multiple processes
			field := v.Target.FieldByName(procName)
			if !field.IsValid() {
				log.Printf("Warning: Invalid field for process %s", procName)
				continue
			}
			// Count elements in the process
			countElements(field, i, elements, q.ProcessElements)
		} else {

			// Get the fields of a single process process
			matched := false
			for j := 0; j < len(v.Fields); j++ {
				fieldName := v.Fields[j].Name
				for element, info := range elements {

					if strings.HasPrefix(fieldName, "From") {
						q.ProcessElements[i]["SequenceFlow"]++
					}

					if info.exactMatch {
						if fieldName == element {
							q.ProcessElements[i][info.name]++
							matched = true
							break
						}
					} else if !matched && strings.Contains(fieldName, element) {
						q.ProcessElements[i][info.name]++
						matched = true
						break
					}

				}

			}
		}

	}

	log.Printf("Process elements: %v", q.ProcessElements)
	log.Print("-------------------------")
}

// countElements counts the BPMN elements in a reflected struct field
func countElements(field reflect.Value, processIndex int, elements map[string]struct {
	name       string
	exactMatch bool
}, processElements map[int]map[string]int) {

	for j := 0; j < field.NumField(); j++ {
		fieldName := field.Type().Field(j).Name

		// Only process struct fields
		if field.Field(j).Kind() != reflect.Struct {
			continue
		}

		// Handle sequence flows first
		if strings.HasPrefix(fieldName, "From") {
			processElements[processIndex]["SequenceFlow"]++
			continue
		}

		// Handle other elements
		matched := false
		for element, info := range elements {
			if info.exactMatch {
				if fieldName == element {
					processElements[processIndex][info.name]++
					matched = true
					break
				}
			} else if !matched && strings.Contains(fieldName, element) {
				processElements[processIndex][info.name]++
				matched = true
				break
			}
		}
	}
}
