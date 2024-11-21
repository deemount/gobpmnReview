package main

import (
	"log"
	"reflect"
	"strings"
)

// quantity holds all the quantities of the BPMN elements
// in the BPMN model. It is used to count the number of elements
type quantity struct {
	Pool                   int
	Process                int
	Participant            int
	StartEvent             int // Note: I want to get rid of these from here to Flow
	EndEvent               int
	IntermediateCatchEvent int
	IntermediateThrowEvent int
	InclusiveGateway       int
	ExclusiveGateway       int
	ParallelGateway        int
	Task                   int
	UserTask               int
	ScriptTask             int
	Flow                   int // Note: all between should till here should be removed. ProcessElements holds all the values
	ProcessElements        map[int]map[string]int
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
	q.countFieldsInProcess(v)
}

// countFieldsInProcess counts the fields in a process and stores them in a map.
func (q *quantity) countFieldsInProcess(v *reflectValue) {
	if q.ProcessElements == nil {
		q.ProcessElements = make(map[int]map[string]int)
	}
	elements := map[string]string{
		"StartEvent":             "StartEvent",
		"EndEvent":               "EndEvent",
		"IntermediateCatchEvent": "IntermediateCatchEvent",
		"IntermediateThrowEvent": "IntermediateThrowEvent",
		"ParallelGateway":        "ParallelGateway",
		"ExclusiveGateway":       "ExclusiveGateway",
		"InclusiveGateway":       "InclusiveGateway",
		"Task":                   "Task", // Note: Task is a very general term, which are also used in the other task elements. strings.Contains are not very helpful here (look below)
		"UserTask":               "UserTask",
		"ScriptTask":             "ScriptTask",
		"From":                   "SequenceFlow",
	}

	for i, field := range v.ProcessName {
		t := v.Target.FieldByName(field)
		if q.ProcessElements[i] == nil {
			q.ProcessElements[i] = make(map[string]int)
		}
		if t.IsValid() && t.CanInterface() {
			for j := 0; j < t.NumField(); j++ {
				name := t.Type().Field(j).Name
				if t.Field(j).Kind() == reflect.Struct {
					for element, counter := range elements {
						log.Printf("Name: %v, Element: %v, Counter: %v", name, element, counter)
						if strings.Contains(name, element) && !strings.Contains(name, "From") {
							q.ProcessElements[i][counter]++ // Bug: it's counting wrong, so that in activities, I haven't the right number of activities
						}
					}
					if strings.HasPrefix(name, "From") {
						q.ProcessElements[i]["SequenceFlow"]++
					}
				}
			}
		}
	}
	log.Printf("Process elements: %v", q.ProcessElements)
	log.Print("-------------------------")
}

/*
 * @Note: used in a single process
 */

// countProcess counts all the processes in the BPMN model.
// Ruleset:
//   - If the field contains the word "Process" then it is a process.
func (q *quantity) countProcess(field string) {
	if strings.Contains(field, "Process") {
		q.Process++
	}
}

// countFlow counts all the flows in the BPMN model.
// Ruleset:
//   - If the field contains the word "From" then it is a flow.
func (q *quantity) countFlow(field string) {
	if strings.Contains(field, "From") {
		q.Flow++
	}
}

// countElement counts all the elements in the BPMN model
// and increments the counter for each element.
// Ruleset:
//   - If the field contains one of the words below and without the word "From"
//     then it is an element.
//   - If the field contains the word from one of the words below
func (q *quantity) countElement(field string) {

	// Define a mapping of element types to their corresponding counters.
	elementCounters := map[string]*int{
		"StartEvent":             &q.StartEvent,
		"EndEvent":               &q.EndEvent,
		"IntermediateCatchEvent": &q.IntermediateCatchEvent,
		"IntermediateThrowEvent": &q.IntermediateThrowEvent,
		"InclusiveGateway":       &q.InclusiveGateway,
		"ExclusiveGateway":       &q.ExclusiveGateway,
		"ParallelGateway":        &q.ParallelGateway,
		"UserTask":               &q.UserTask,
		"Task":                   &q.Task,
		"ScriptTask":             &q.ScriptTask,
	}

	// Check if the field is not a flow (i.e., it does not contain "From").
	if !strings.Contains(field, "From") {
		// Iterate over the element counters and increment the corresponding counter if the field matches.
		for element, counter := range elementCounters {
			if strings.Contains(field, element) {
				(*counter)++
				break
			}
		}
	}
}
