package main

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

// ElementProcessor handles BPMN element processing
type ElementProcessor struct {
	value    *reflectValue
	quantity *quantity
}

// NewElementProcessor creates a new ElementProcessor instance
func NewElementProcessor(v *reflectValue, q *quantity) *ElementProcessor {
	return &ElementProcessor{
		value:    v,
		quantity: q,
	}
}

// ProcessingConfig holds the configuration for element processing
type ProcessingConfig struct {
	field     reflect.Value
	numFields int
	indices   map[string]int
	counts    map[string]int
}

// createProcessingConfig initializes processing configuration
func (ep *ElementProcessor) createProcessingConfig(processIndex int, field reflect.Value) *ProcessingConfig {
	return &ProcessingConfig{
		field:     field,
		numFields: field.NumField(),
		indices:   initializeIndices(),
		counts:    ep.quantity.ProcessElements[processIndex],
	}
}

/*
 * @multiple processes
 */

// ProcessMultipleElement processes all elements across multiple processes
func (ep *ElementProcessor) ProcessMultipleElement() {
	for i, processName := range ep.value.ProcessName {
		field := ep.value.Target.FieldByName(processName)
		if !field.IsValid() {
			log.Printf("Invalid field for process: %s", processName)
			continue
		}
		ep.processProcess(i, field)
	}
}

// processProcess handles processing for a single process
func (ep *ElementProcessor) processProcess(processIndex int, field reflect.Value) {
	config := ep.createProcessingConfig(processIndex, field)

	for j := 0; j < config.numFields; j++ {
		ep.processField(processIndex, j, config)
	}
}

// processField handles processing of individual fields
func (ep *ElementProcessor) processField(processIndex, fieldIndex int, config *ProcessingConfig) {
	fieldInfo := extractFieldInfo(config.field, fieldIndex) // helper function in internal.go

	if !isValidField(fieldInfo) {
		return
	}

	ep.handleElement(processIndex, fieldInfo, config)
}

/*
 * @single process
 */

// ProcessSingleElement processes elements for a single process
func (ep *ElementProcessor) ProcessSingleElement() {
	for i := 0; i < ep.value.TargetNumField; i++ {

		field := ep.value.Target.Field(i)
		fieldType := ep.value.Target.Type().Field(i)

		if field.Kind() == reflect.Bool {
			continue
		}

		if field.Kind() == reflect.Struct {
			ep.processSingleStructField(field, fieldType)
		}
	}
}

// processSingleStructField processes a single struct field
func (ep *ElementProcessor) processSingleStructField(field reflect.Value, fieldType reflect.StructField) {

	typ := field.FieldByName("Type").String()
	hash := field.FieldByName("Hash").String()

	if strings.HasPrefix(fieldType.Name, "From") {
		callFlows() // Note: helper function in internal.go (not ready yet)
		return
	}

	handler := newElementHandler(ep.value.Process[0], fieldType.Name, typ, hash)
	handler.handleElement()

}

// callFlows this method is not completed
// (Note: this method is used for a single process in model)
// ... in progress
func callFlows() {
	log.Printf("Flow element detected")
}

/*
 * @handlers
 */

// handleElement processes different types of BPMN elements
func (ep *ElementProcessor) handleElement(processIndex int, info FieldInfo, config *ProcessingConfig) {
	switch info.typ {
	case "startevent":
		ep.handleStartEvent(processIndex, info, config)
	case "event":
		ep.handleEvent(processIndex, info, config)
	case "flow":
		ep.handleFlow(processIndex, info, config)
	case "gateway":
		ep.handleGateway(processIndex, info, config)
	case "activity":
		ep.handleActivity(processIndex, info, config)
	}
}

// StartEventIndices holds indices for start events
type StartEventIndices struct {
	startEvent int
}

// handleStartEvent ...
func (ep *ElementProcessor) handleStartEvent(processIndex int, info FieldInfo, config *ProcessingConfig) {

	// Create local variables for the indices
	indices := StartEventIndices{
		startEvent: config.indices["startEventIndex"],
	}

	idx := &indices.startEvent

	if *idx < config.counts["StartEvent"] {

		if processIndex > 0 && *idx == 0 {
			info.typ = "event"
		}

		el := ep.value.Process[processIndex].MethodByName("GetStartEvent").Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		el.MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(info.name)})
		// a startevent has only one outgoing
		el.MethodByName("SetOutgoing").Call([]reflect.Value{reflect.ValueOf(1)})
		out := el.MethodByName("GetOutgoing").Call([]reflect.Value{reflect.ValueOf(0)})[0]
		out.MethodByName("SetFlow").Call([]reflect.Value{reflect.ValueOf(info.nextHash)})
	}
}

// EventIndices holds indices for different types of events
type EventIndices struct {
	catch int
	throw int
	end   int
}

// handleEvent ...
func (ep *ElementProcessor) handleEvent(processIndex int, info FieldInfo, config *ProcessingConfig) {

	// Create local variables for the indices
	indices := EventIndices{
		catch: config.indices["intermediateCatchEventIndex"],
		throw: config.indices["intermediateThrowEventIndex"],
		end:   config.indices["endEventIndex"],
	}

	// Handle the gateway based on its type
	switch info.extName {
	case "IntermediateCatchEvent":
		ep.handleIntermediateCatchEvent(processIndex, info, config, &indices.catch)
	case "IntermediateThrowEvent":
		ep.handleIntermediateThrowEvent(processIndex, info, config, &indices.throw)
	case "EndEvent":
		ep.handleEndEvent(processIndex, info, config, &indices.end)
	}

	// Update the indices in the config
	config.indices["intermediateCatchEventIndex"] = indices.catch
	config.indices["intermediateThrowEventIndex"] = indices.throw
	config.indices["endEventIndex"] = indices.end

}

// handleIntermediateCatchEvent ...
func (ep *ElementProcessor) handleIntermediateCatchEvent(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["IntermediateCatchEvent"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// handleIntermediateThrowEvent ...
func (ep *ElementProcessor) handleIntermediateThrowEvent(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["IntermediateThrowEvent"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// handleEndEvent ...
func (ep *ElementProcessor) handleEndEvent(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["EndEvent"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// GatewayIndices holds indices for different types of gateways
type GatewayIndices struct {
	inclusive int
	exclusive int
	parallel  int
}

// handleGateway ...
// (Note: no gateway is actually handled in this function)
func (ep *ElementProcessor) handleGateway(processIndex int, info FieldInfo, config *ProcessingConfig) {

	// Create local variables for the indices
	indices := GatewayIndices{
		inclusive: config.indices["inclusiveGatewayIndex"],
		exclusive: config.indices["exclusiveGatewayIndex"],
		parallel:  config.indices["parallelGatewayIndex"],
	}

	// Handle the gateway based on its type
	switch info.extName {
	case "InclusiveGateway":
		ep.handleInclusiveGateway(processIndex, info, config, &indices.inclusive)
	case "ExclusiveGateway":
		ep.handleExclusiveGateway(processIndex, info, config, &indices.exclusive)
	case "ParallelGateway":
		ep.handleParallelGateway(processIndex, info, config, &indices.parallel)
	}

	// Update the indices in the config
	config.indices["inclusiveGatewayIndex"] = indices.inclusive
	config.indices["exclusiveGatewayIndex"] = indices.exclusive
	config.indices["parallelGatewayIndex"] = indices.parallel

}

// handleInclusiveGateway processes inclusive gateway elements
func (ep *ElementProcessor) handleInclusiveGateway(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["InclusiveGateway"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// handleExclusiveGateway processes exclusive gateway elements
func (ep *ElementProcessor) handleExclusiveGateway(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["ExclusiveGateway"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// handleParallelGateway processes parallel gateway elements
func (ep *ElementProcessor) handleParallelGateway(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["ParallelGateway"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// ActivityIndices holds indices for different types of activities
type ActivityIndices struct {
	task       int
	userTask   int
	scriptTask int
}

// handleActivity ...
func (ep *ElementProcessor) handleActivity(processIndex int, info FieldInfo, config *ProcessingConfig) {

	// Create local variables for the indices
	indices := ActivityIndices{
		task:       config.indices["taskIndex"],
		userTask:   config.indices["userTaskIndex"],
		scriptTask: config.indices["scriptTaskIndex"],
	}

	// Handle the activity based on its type
	switch info.extName {
	case "Task":
		ep.handleTask(processIndex, info, config, &indices.task)
	case "UserTask":
		ep.handleUserTask(processIndex, info, config, &indices.userTask)
	case "ScriptTask":
		ep.handleScriptTask(processIndex, info, config, &indices.scriptTask)
	}

	// Update the indices in the config
	config.indices["taskIndex"] = indices.task
	config.indices["userTaskIndex"] = indices.userTask
	config.indices["scriptTaskIndex"] = indices.scriptTask

}

// handleTask ...
func (ep *ElementProcessor) handleTask(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["Task"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// handleUserTask ...
func (ep *ElementProcessor) handleUserTask(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["UserTask"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// handleScriptTask ...
func (ep *ElementProcessor) handleScriptTask(processIndex int, info FieldInfo, config *ProcessingConfig, idx *int) {
	if *idx < config.counts["ScriptTask"] {
		methodName := fmt.Sprintf("Get%s", info.extName)
		el := ep.value.Process[processIndex].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(*idx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*idx)++
	}
}

// FlowIndices holds indices for flow elements
type FlowIndices struct {
	flow int
}

// handleFlow ...
func (ep *ElementProcessor) handleFlow(processIndex int, info FieldInfo, config *ProcessingConfig) {

	// Create local variables for the indices
	indices := FlowIndices{
		flow: config.indices["flowIndex"],
	}

	flowIdx := &indices.flow

	if *flowIdx < config.counts["Flow"] {
		el := ep.value.Process[processIndex].MethodByName("GetSequenceFlow").Call([]reflect.Value{reflect.ValueOf(*flowIdx)})[0]
		el.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(info.typ), reflect.ValueOf(info.hash)})
		(*flowIdx)++
	}

	config.indices["flowIndex"] = indices.flow

}

/*
 * @helpers
 */

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

// getNextHash ...
// (Note: maybe redudant? look at reflectValue method "next". Refactor?)
func getNextHash(field reflect.Value, j, numFields int) string {
	if j+1 < numFields {
		return field.Field(j + 1).FieldByName("Hash").String()
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
