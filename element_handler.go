package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// elementHandler handles BPMN process elements in a single process
type elementHandler struct {
	process reflect.Value
	method  string
	name    string
	typ     string
	hash    string
}

// newElementHandler creates a new elementHandler instance
func newElementHandler(p reflect.Value, n, t, h string) *elementHandler {
	return &elementHandler{
		process: p,
		method:  "Get" + n,
		name:    n,
		typ:     t,
		hash:    h,
	}
}

// handleElement processes a single BPMN element
func (h *elementHandler) handleElement() error {
	element := h.getElementType()
	if element == "" {
		return fmt.Errorf("unknown element type: %s", h.name)
	}

	el, err := h.invokeGetMethod(0)
	if err != nil {
		return fmt.Errorf("failed to invoke get method: %w", err)
	}

	if err := h.setCommonAttributes(el); err != nil {
		return fmt.Errorf("failed to set common attributes: %w", err)
	}

	switch element {
	case startEvent:
		return h.handleStartEvent(el)
	case endEvent:
		return h.handleEndEvent(el)
	default:
		return nil // Other elements don't need special handling
	}
}

// elementMatchType represents a BPMN element type with matching information
type elementMatchType struct {
	element processElement
	exact   bool
}

// getElementType determines the type of BPMN element using Two-Pass Matching
func (h *elementHandler) getElementType() processElement {

	// Add logging for debugging
	defer func() {
		log.Printf("Element type determination for %s completed", h.name)
	}()

	// Order matters: more specific types should come before general types
	elementTypes := []struct {
		name    string
		element elementMatchType
	}{
		// exact matches (specific types)
		{"UserTask", elementMatchType{userTask, true}},
		{"ScriptTask", elementMatchType{scriptTask, true}},
		{"StartEvent", elementMatchType{startEvent, true}},
		{"EndEvent", elementMatchType{endEvent, true}},
		{"IntermediateCatchEvent", elementMatchType{intermediateCatchEvent, true}},
		{"IntermediateThrowEvent", elementMatchType{intermediateThrowEvent, true}},
		{"InclusiveGateway", elementMatchType{inclusiveGateway, true}},
		{"ExclusiveGateway", elementMatchType{exclusiveGateway, true}},
		{"ParallelGateway", elementMatchType{parallelGateway, true}},
		// partial matches (general types)
		{"Task", elementMatchType{task, false}}, // generic task should be checked last
	}

	// First try exact matches
	for _, et := range elementTypes {
		if et.element.exact && h.name == et.name {
			return et.element.element
		}
	}

	// Then try partial matches
	for _, et := range elementTypes {
		if !et.element.exact && strings.Contains(h.name, et.name) {
			return et.element.element
		}
	}

	return ""
}

// invokeGetMethod invokes the Get method for an element
func (h *elementHandler) invokeGetMethod(index int) (reflect.Value, error) {
	method := h.process.MethodByName(h.method)
	if !method.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid method: %s", h.method)
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(index)})
	if len(results) == 0 {
		return reflect.Value{}, fmt.Errorf("method %s returned no results", h.method)
	}

	return results[0], nil
}

// setCommonAttributes sets common attributes for all elements
func (h *elementHandler) setCommonAttributes(el reflect.Value) error {
	if err := h.invokeMethod(el, "SetID", h.typ, h.hash); err != nil {
		return err
	}
	return h.invokeMethod(el, "SetName", h.name)
}

// handleStartEvent processes start event specific attributes
func (h *elementHandler) handleStartEvent(el reflect.Value) error {
	if err := h.invokeMethod(el, "SetOutgoing", 1); err != nil {
		return err
	}

	out, err := h.invokeMethodWithReturn(el, "GetOutgoing", 0)
	if err != nil {
		return err
	}

	return h.invokeMethod(out, "SetFlow", h.hash)
}

// handleEndEvent processes end event specific attributes
func (h *elementHandler) handleEndEvent(el reflect.Value) error {
	if err := h.invokeMethod(el, "SetIncoming", 1); err != nil {
		return err
	}

	in, err := h.invokeMethodWithReturn(el, "GetIncoming", 0)
	if err != nil {
		return err
	}

	return h.invokeMethod(in, "SetFlow", h.hash)
}

// invokeMethod invokes a method with parameters
func (h *elementHandler) invokeMethod(el reflect.Value, methodName string, args ...interface{}) error {
	method := el.MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("invalid method: %s", methodName)
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	method.Call(reflectArgs)
	return nil
}

// invokeMethodWithReturn invokes a method and returns its result
func (h *elementHandler) invokeMethodWithReturn(el reflect.Value, methodName string, index int) (reflect.Value, error) {
	method := el.MethodByName(methodName)
	if !method.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid method: %s", methodName)
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(index)})
	if len(results) == 0 {
		return reflect.Value{}, fmt.Errorf("method %s returned no results", methodName)
	}

	return results[0], nil
}

// CallMethods is the main entry point for processing BPMN elements
func CallMethods(p reflect.Value, n, t, h string) error {
	handler := newElementHandler(p, n, t, h)
	return handler.handleElement()
}

// CallFlows handles flow-related processing (to be implemented)
func CallFlows() error {
	// TODO: Implement flow handling
	return nil
}
