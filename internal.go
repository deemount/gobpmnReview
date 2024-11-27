package main

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"
)

/*
 * @M hash and typ
 */

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
