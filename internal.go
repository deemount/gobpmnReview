package main

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"reflect"
	"regexp"
	"strings"
)

/*
 * @Single Process
 */

// callMethods ...
// Note: this method is uesed for a single process in a model and needs to refactored.
// ... in progress (the solution is in multiple processes)
// Note: if you put more than one task, event and so on in a process, the method is not working correctly.
func callMethods(p reflect.Value, n, t, h string) {
	method := "Get" + n
	switch true {
	// events
	case strings.Contains(n, "StartEvent"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
		// a startevent has only one outgoing
		el.MethodByName("SetOutgoing").Call([]reflect.Value{reflect.ValueOf(1)})
		out := el.MethodByName("GetOutgoing").Call([]reflect.Value{reflect.ValueOf(0)})[0]
		outFlowMethod := out.MethodByName("SetFlow")
		outFlowMethod.Call([]reflect.Value{reflect.ValueOf(h)}) // Note: the h value must be shown to the next element, which is refering to (like )
	case strings.Contains(n, "EndEvent"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
		// - An EndEvent usually has exactly one incoming sequence flow,
		//   as it often represents the end of a straightforward or clearly structured path.
		// - An EndEvent can have multiple incoming sequence flows,
		//   if it is the endpoint of paths that are merged at a Parallel Gateway or Inclusive Gateway.
		//   This is useful for bundling different process paths at a single endpoint.
		el.MethodByName("SetIncoming").Call([]reflect.Value{reflect.ValueOf(1)})
		in := el.MethodByName("GetIncoming").Call([]reflect.Value{reflect.ValueOf(0)})[0] // get the first incoming sequence flow
		inFlowMethod := in.MethodByName("SetFlow")
		inFlowMethod.Call([]reflect.Value{reflect.ValueOf(h)}) // Note: the h value must be shown to the next element, which is refer to
	case strings.Contains(n, "IntermediateCatchEvent"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	case strings.Contains(n, "IntermediateThrowEvent"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	// gateways
	case strings.Contains(n, "InclusiveGateway"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	case strings.Contains(n, "ExclusiveGateway"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	case strings.Contains(n, "ParallelGateway"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	// Note: the following elements are activities. All cases must be refactored.
	//       The contains method is buggy, because Task is always in other strings, too.
	// tasks
	case strings.Contains(n, "Task"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	// user tasks
	case strings.Contains(n, "UserTask"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	// script tasks
	case strings.Contains(n, "ScriptTask"):
		el := p.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(0)})[0]
		callSetters(el, method, n, t, h)
	}
}

// callFlows this method is not completed
// (Note: this method is used for a single process in model)
// ... in progress
func callFlows() {
	// not completed
}

// callSetters can be a helper funtion to reduce the code in callMethods.
func callSetters(el reflect.Value, method, name, typ, hash string) {
	el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
	el.MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(name)})
}

// methods is a helper function to call the methods of a single process.
// The function returns a struct slice of methods and their arguments.
// The arguments are the quantities of the BPMN elements.
// If the quantity is greater than 0, then the method is called with the
// corresponding argument. Eitherwise, the method is not called and
// the struct slice value in the model will be nil.

// rebuild this function
func methods() []struct {
	name string
	arg  int
} {
	return []struct {
		name string
		arg  int
	}{
		{"SetStartEvent", 1},
		{"SetTask", 1},
		{"SetEndEvent", 1},
		{"SetSequenceFlow", 2},
	}
}

/*
 * @M hash and typ
 */

// getNextHash ...
// (Note: maybe redudant? look at reflectValue method "next". Refactor?)
func getNextHash(field reflect.Value, j, numFields int) string {
	if j+1 < numFields {
		return field.Field(j + 1).FieldByName("Hash").String()
	}
	return ""
}

// typ ...
func typ(n string) string {

	// define a mapping of BPMN element types to their corresponding names.
	bpmnTypes := map[string]string{
		"IsExecutable":           "config",
		"Process":                "process",
		"StartEvent":             "startevent",
		"EndEvent":               "event",
		"IntermediateCatchEvent": "event",
		"IntermediateThrowEvent": "event",
		"InclusiveGateway":       "gateway",
		"ExclusiveGateway":       "gateway",
		"ParallelGateway":        "gateway",
		"Task":                   "activity",
		"UserTask":               "activity",
		"ScriptTask":             "activity",
	}

	// check if the name contains "From" to determine if it's a flow.
	if strings.Contains(n, "From") {
		return "flow"
	}

	// iterate over the BPMN types and return the type if the name contains the type.
	for t, bpmnType := range bpmnTypes {
		// Note: the contains method in this case is buggy (or not), because Task is always in other strings.
		if strings.Contains(n, t) {
			return bpmnType
		}
	}

	// Return an empty string if no matching type is found.
	return ""
}

// hash returns a hash value of a given string.
// It uses the FNV-1a algorithm to generate a hash value.
// The method returns a BPMN struct with the hash value.
// The argument typ is the type of the BPMN element.
func hash(typ string) (BPMN, error) {

	n := 8
	b := make([]byte, n)
	c := fnv.New32a()

	if _, err := rand.Read(b); err != nil {
		return BPMN{}, err
	}
	s := fmt.Sprintf("%x", b)

	if _, err := c.Write([]byte(s)); err != nil {
		return BPMN{}, err
	}
	defer c.Reset()

	result := BPMN{
		Type: typ, // Note: needs to be reconsidered. Is this the right method, to put the typ value into the Type field?
		Hash: fmt.Sprintf("%x", string(c.Sum(nil))),
	}

	return result, nil
}

/*
 * @ handlers for multiple processes
 */

// handleStartEvent ...
func handleStartEvent(v *reflectValue, i int, name, extName, typ, hash, nextHash string, numStartEvent int, startEventIndex *int) {
	if !strings.HasPrefix(name, "From") && *startEventIndex < numStartEvent {
		if i > 0 && *startEventIndex == 0 {
			typ = "event"
		}
		el := v.Process[i].MethodByName("GetStartEvent").Call([]reflect.Value{reflect.ValueOf(*startEventIndex)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
		el.MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(name)})
		// a startevent has only one outgoing
		el.MethodByName("SetOutgoing").Call([]reflect.Value{reflect.ValueOf(1)})
		out := el.MethodByName("GetOutgoing").Call([]reflect.Value{reflect.ValueOf(0)})[0]
		out.MethodByName("SetFlow").Call([]reflect.Value{reflect.ValueOf(nextHash)})
		(*startEventIndex)++
	}
}

// handleEvent ...
func handleEvent(v *reflectValue, i int, name, extName, typ, hash string, numIntermediateCatchEvent, numIntermediateThrowEvent, numEndEvent int, intermediateCatchEventIndex, intermediateThrowEventIndex, endEventIndex *int) {
	if !strings.HasPrefix(name, "From") {
		switch extName {
		case "EndEvent":
			if *endEventIndex < numEndEvent {
				el := v.Process[i].MethodByName("GetEndEvent").Call([]reflect.Value{reflect.ValueOf(*endEventIndex)})[0]
				el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
				(*endEventIndex)++
			}
		case "CatchEvent":
			if *intermediateCatchEventIndex < numIntermediateCatchEvent {
				el := v.Process[i].MethodByName("GetIntermediateCatchEvent").Call([]reflect.Value{reflect.ValueOf(*intermediateCatchEventIndex)})[0]
				el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
				(*intermediateCatchEventIndex)++
			}
		case "ThrowEvent":
			if *intermediateThrowEventIndex < numIntermediateThrowEvent {
				el := v.Process[i].MethodByName("GetIntermediateThrowEvent").Call([]reflect.Value{reflect.ValueOf(*intermediateThrowEventIndex)})[0]
				el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
				(*intermediateThrowEventIndex)++
			}
		}
	}
}

// handleActivity ...
func handleActivity(v *reflectValue, i int, name, extName, typ, hash string, numTask, numUserTask, numScriptTask int, taskIndex, userTaskIndex, scriptTaskIndex *int) {
	if !strings.Contains(name, "From") {
		switch extName {
		case "UserTask":
			// test it for many of them
			// Note: look at the bug in default
			if *userTaskIndex < numUserTask {
				el := v.Process[i].MethodByName("GetUserTask").Call([]reflect.Value{reflect.ValueOf(*userTaskIndex)})[0]
				el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
				(*userTaskIndex)++
			}
		case "ScriptTask":
			if *scriptTaskIndex < numScriptTask {
				el := v.Process[i].MethodByName("GetScriptTask").Call([]reflect.Value{reflect.ValueOf(*scriptTaskIndex)})[0]
				el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
				(*scriptTaskIndex)++
			}
		default:
			if *taskIndex < numTask {
				el := v.Process[i].MethodByName("GetTask").Call([]reflect.Value{reflect.ValueOf(*taskIndex)})[0]
				el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
				(*taskIndex)++
			}
		}
	}
}

// FieldInfo holds information about a field being processed
type FieldInfo struct {
	name     string
	typ      string
	hash     string
	nextHash string
	extName  string
}

// extractFieldInfo gathers all necessary field information
func extractFieldInfo(field reflect.Value, index int) FieldInfo {
	return FieldInfo{
		name:     field.Type().Field(index).Name,
		typ:      field.Field(index).FieldByName("Type").String(),
		hash:     field.Field(index).FieldByName("Hash").String(),
		nextHash: getNextHash(field, index, field.NumField()),
		extName:  extractLastTwoWords(field.Type().Field(index).Name),
	}
}

func initializeIndices() map[string]int {
	return map[string]int{
		"startEventIndex":             0,
		"endEventIndex":               0,
		"intermediateCatchEventIndex": 0,
		"intermediateThrowEventIndex": 0,
		"taskIndex":                   0,
		"userTaskIndex":               0,
		"scriptTaskIndex":             0,
		"flowIndex":                   0,
		"inclusiveGatewayIndex":       0,
		"parallelGatewayIndex":        0,
	}
}

func isValidField(info FieldInfo) bool {
	return info.name != "" && info.typ != ""
}

/*
 * @ strings
 */

// extractPrefixBeforeProcess extracts the prefix before the word "Process" in a string.
// The method returns the prefix as a string, which is used as the name of the process.
// (Note: the mthod is used two times: first in newReflectValue (reflect_di.go) and second countFieldsInPool (quantities.go))
func extractPrefixBeforeProcess(input string) string {
	re := regexp.MustCompile(`([A-Za-z]+)Process$`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// extractLastTwoWords extracts the last two words from a string
// (Note: this method is used to get the elements out of a string)
func extractLastTwoWords(input string) string {
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	words := re.FindAllString(input, -1)
	if len(words) < 2 {
		return strings.Join(words, "")
	}
	lastTwo := words[len(words)-2:]
	return strings.Join(lastTwo, "")
}

// getLastPart gibt den letzten Teil eines Snake-Case-Strings zurück.
// If no Snake-Case-String, return the original string.
func getLastPart(input string) string {
	if strings.Contains(input, "_") {
		parts := strings.Split(input, "_")
		// Gib den letzten Teil zurück
		return parts[len(parts)-1]
	}
	// Kein Snake-Case-String, gib den Original-String zurück
	return input
}

// trimDoublePart cuts the double part of a string.
func trimDoublePart(input string) string {
	parts := strings.Split(input, "_")
	if len(parts) >= 3 && parts[0] == parts[1] {
		// if the first two parts are equal, we return the rest of the string
		return strings.Join(parts[1:], "_")
	}
	// if the first two parts are not equal, we return the input string
	return input
}
