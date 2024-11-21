package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// BPMN ...
type BPMN struct {
	Pos, Width, Height, X, Y int
	Type, Hash               string
}

// NewReflectDI ...
// (Note: looks a little bit chaotically and should be refactored)
func NewReflectDI(p any) any {
	v := newReflectValue(p)
	v.TargetNumField = v.Target.NumField()
	v.Def = v.Target.FieldByName("Def")
	if !v.Def.IsValid() {
		fmt.Println("Def field is not valid")
		return nil
	}
	if v.Def.CanSet() {
		definitions := NewDefinitions() // Note: I'm not sure if I should use the NewDefinitions function here. Maybe reflecting it?
		definitions.SetDefaultAttributes()
		v.Def.Set(reflect.ValueOf(definitions))
	}
	m := new(mapping) // Note: Mapping and Quantity is a helper structure and can be put together in a single structure?
	q := new(quantity)
	m.Assign(v)
	if len(m.Anonym) > 0 {
		v.pool(q, m)
	} else {
		q.Process += 1 // Note: Must be counted and set here. Actually it is used to create a single process (e.g. In m.BPMNType the first element at index 1 is process).
	}
	v.reflectProcess(q)
	di := v.target(q, m) // Note: any solutions with generics here?
	v.process(q, m)      // Note: with the maps it is possible to build the collaboration after the processes
	v.collaboration(q)

	return di
}

// reflectValue
// (Note: since I'm using the extreme programming methodology,
// I'm having a lot of (my not needed) slices to handle the reflection)
type reflectValue struct {
	Name            string
	Fields          []reflect.StructField
	Target          reflect.Value
	TargetNumField  int
	Def             reflect.Value
	Pool            reflect.Value
	Process         []reflect.Value
	ProcessName     []string
	ProcessType     []string
	ProcessHash     []string
	ProcessExec     []bool
	ParticipantName []string
	ParticipantHash []string
}

// newReflectValue ...
func newReflectValue(p interface{}) *reflectValue {
	typeOf := reflect.TypeOf(p)
	return &reflectValue{
		Name:   extractPrefixBeforeProcess(typeOf.Name()),
		Fields: reflect.VisibleFields(typeOf),
		Target: reflect.New(typeOf).Elem(),
	}
}

// target
func (v *reflectValue) target(q *quantity, m *mapping) any {
	if len(m.Anonym) > 0 {
		for _, anonymField := range m.Anonym {
			targetFieldByName := v.Target.FieldByName(anonymField)
			targetNumField := targetFieldByName.NumField()
			for i := 0; i < targetNumField; i++ {
				name := targetFieldByName.Type().Field(i).Name
				switch targetFieldByName.Field(i).Kind() {
				case reflect.Bool:
					v.config(name, i, targetFieldByName)
				case reflect.Struct:
					q.countFlow(name)    // Note: should I count here?
					q.countElement(name) // Note: should I count here?
					v.current(i, targetFieldByName)
					v.next(i, targetFieldByName)
				}
			}
		}
	} else {
		v.nonAnonym(q, m)
	}
	return v.Target.Interface()
}

// config sets the IsExecutable field to true if the name contains "IsExecutable" and index is 0.
func (v *reflectValue) config(name string, index int, target reflect.Value) {
	if strings.Contains(name, "IsExecutable") && index == 0 {
		target.Field(0).SetBool(true)
	}
}

// current sets the hash value of the current field if it is empty.
func (v *reflectValue) current(index int, target reflect.Value) {
	h := fmt.Sprintf("%s", target.Field(index).FieldByName("Hash"))
	if h == "" {
		typ := typ(target.Type().Field(index).Name)
		hash, _ := hash(typ)
		target.Field(index).Set(reflect.ValueOf(hash))
	}
}

// next sets the hash value of the next field.
func (v *reflectValue) next(index int, target reflect.Value) {
	if index+1 < target.NumField() {
		typ := typ(target.Type().Field(index + 1).Name)
		hash, _ := hash(typ)
		target.Field(index + 1).Set(reflect.ValueOf(hash))
	}
}

/*
 * @Note: The following methods are used in a single process.
 */

// nonAnonym sets the IsExecutable field to true and populates the reflection fields with hash values.
// (Note: the method is used in a single process).
func (v *reflectValue) nonAnonym(q *quantity, m *mapping) {
	v.setIsExecutable(m.Config)
	v.populateReflectionFields(q, m.BPMNType)
}

// setIsExecutable sets the IsExecutable field to true (Note: maybe redudant; look at method config).
func (v *reflectValue) setIsExecutable(configFields map[int]string) {
	for _, configField := range configFields {
		f := v.Target.FieldByName(configField)
		f.SetBool(true)
		break
	}
}

// populateReflectionFields populates the reflection fields with hash values.
// (Note: this method is used in a single process (Note: maybe redudant;).
func (v *reflectValue) populateReflectionFields(q *quantity, reflectionFields map[int]string) {
	for _, field := range reflectionFields {
		q.countProcess(field)
		q.countFlow(field)
		q.countElement(field)
		f := v.Target.FieldByName(field)
		fName, _ := v.Target.Type().FieldByName(field)
		typ := typ(fName.Name)
		hash, _ := hash(typ)
		f.Set(reflect.ValueOf(hash))
	}
}

// singleProcess ...
// (Note: this method is used in a single process and not in usage. Do I need it?)
func (v *reflectValue) singleProcess(q *quantity, m *mapping) {
	for i := 0; i < v.TargetNumField; i++ {
		field := v.Target.Field(i)
		fieldType := v.Target.Type().Field(i)
		if strings.Contains(fieldType.Name, "IsExecutable") {
			v.Process[0].MethodByName("SetIsExecutable").Call([]reflect.Value{reflect.ValueOf(field.Bool())})
		}
		if fieldType.Name == "Process" {
			hash := field.FieldByName("Hash").String()
			v.Process[0].MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(strings.ToLower(fieldType.Name)), reflect.ValueOf(hash)})
			v.Process[0].MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(fieldType.Name)})
		}
	}
	methodCalls := methods(q)
	for _, call := range methodCalls {
		method := v.Process[0].MethodByName(call.name)
		if method.IsValid() && call.arg > 0 {
			method.Call([]reflect.Value{reflect.ValueOf(call.arg)})
		}
	}
	v.invokeElementsFromProcess()
}

// invokeElementFromProcess ...
// (Note: this method is used in a single process)
func (v *reflectValue) invokeElementsFromProcess() {
	for i := 0; i < v.TargetNumField; i++ {
		fieldValue := v.Target.Field(i)
		fieldType := v.Target.Type().Field(i)
		switch fieldValue.Kind() {
		case reflect.Bool:
			continue
		case reflect.Struct:
			t := fieldValue.FieldByName("Type").String()
			h := fieldValue.FieldByName("Hash").String()
			if strings.Contains(fieldType.Name, "From") {
				callFlows() // Note: build the flows. Function is incomplete (look at internal.go).
			} else {
				callMethods(v.Process[0], fieldType.Name, t, h)
			}
		}
	}
}

/*
 * @Note: The following methods are used in multiple processes.
 */

// collaboration sets up the collaboration in the BPMN model.
func (v *reflectValue) collaboration(q *quantity) {
	if q.Participant > 0 && (q.Process == q.Participant) {
		v.Def.MethodByName("SetCollaboration").Call([]reflect.Value{})
		collaboration := v.Def.MethodByName("GetCollaboration").Call([]reflect.Value{})[0]
		collaboration.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf("collaboration"), reflect.ValueOf("collaboration")})
		collaboration.MethodByName("SetParticipant").Call([]reflect.Value{reflect.ValueOf(q.Participant)})
		for i := 0; i < q.Participant; i++ {
			participant := collaboration.MethodByName("GetParticipant").Call([]reflect.Value{reflect.ValueOf(i)})[0]
			participant.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf("participant"), reflect.ValueOf(v.ParticipantHash[i])})
			participant.MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(v.ParticipantName[i])})
			participant.MethodByName("SetProcessRef").Call([]reflect.Value{reflect.ValueOf("process"), reflect.ValueOf(v.ProcessHash[i])})
		}
	}
}

// pool assigns the field Pool to a reflected value.
// (Note: walk through m.Anonym  is maybe redudant)
func (v *reflectValue) pool(q *quantity, m *mapping) {
	for _, anonymField := range m.Anonym {
		if strings.Contains(anonymField, "Pool") {
			v.Pool = v.Target.FieldByName(anonymField)
			q.countFieldsInPool(v)
			q.Pool++
			break
		}
	}
}

// reflectProcess reflects the processes in the BPMN model.
// (Note: instead of using two functions to fullfill the reflection for a process,
// it is possible to use a single function).
func (v *reflectValue) reflectProcess(q *quantity) {
	v.Process = make([]reflect.Value, q.Process)
	v.Def.MethodByName("SetProcess").Call([]reflect.Value{reflect.ValueOf(q.Process)})
	for j := 0; j < q.Process; j++ {
		v.Process[j] = v.Def.MethodByName("GetProcess").Call([]reflect.Value{reflect.ValueOf(j)})[0]
	}
}

// process sets the maps and process parameters in the BPMN model.
func (v *reflectValue) process(q *quantity, m *mapping) {
	v.ProcessExec = make([]bool, q.Process)
	v.ProcessType = make([]string, q.Process)
	v.ProcessHash = make([]string, q.Process)
	v.ParticipantHash = make([]string, q.Participant)
	l, j, n := 0, 0, 0
	if q.Pool > 0 {
		for i := 0; i < v.Pool.NumField(); i++ {
			name := v.Pool.Type().Field(i).Name
			field := v.Pool.Field(i)
			if strings.Contains(name, "IsExecutable") {
				if field.IsValid() {
					v.ProcessExec[j] = field.Bool()
					j++
				}
			}
			if strings.Contains(name, "Process") {
				if field.IsValid() {
					v.ProcessHash[l] = field.FieldByName("Hash").String()
					v.ProcessType[l] = field.FieldByName("Type").String()
					l++
				}
			}
			if strings.Contains(name, "Participant") {
				if field.IsValid() {
					v.ParticipantHash[n] = field.FieldByName("Hash").String()
					n++
				}
			}
		}
		for i := 0; i < q.Process; i++ {
			if i == 0 {
				v.Process[i].MethodByName("SetIsExecutable").Call([]reflect.Value{reflect.ValueOf(v.ProcessExec[i])})
			}
			v.Process[i].MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(v.ProcessType[i]), reflect.ValueOf(v.ProcessHash[i])})
			v.Process[i].MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(v.ProcessName[i])})
			for method, arg := range q.ProcessElements[i] {
				methodName := fmt.Sprintf("Set%s", method)
				v.Process[i].MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(arg)})
			}
		}
		v.processElements(q)
	} else {
		v.singleProcess(q, m)
	}
}

// invokeElementFromProcesses processes elements from the quantity and updates the reflectValue.
// (Note: handlers having a lot of args)
func (v *reflectValue) processElements(q *quantity) {

	for i, processName := range v.ProcessName {

		field := v.Target.FieldByName(processName) // Note: this scheme is also used in injectTarget
		numFields := field.NumField()              // Note: this scheme is also used in injectTarget

		// Get the number of elements to build the process
		numFlows := q.ProcessElements[i]["SequenceFlow"]
		numTask := q.ProcessElements[i]["Task"]
		numScriptTask := q.ProcessElements[i]["ScriptTask"]
		numUserTask := q.ProcessElements[i]["UserTask"]
		numStartEvent := q.ProcessElements[i]["StartEvent"]
		numEndEvent := q.ProcessElements[i]["EndEvent"]
		numIntermediateCatchEvent := q.ProcessElements[i]["IntermediateCatchEvent"]
		numIntermediateThrowEvent := q.ProcessElements[i]["IntermediateThrowEvent"]
		numInclusiveGateway := q.ProcessElements[i]["InclusiveGateway"]
		numParallelGateway := q.ProcessElements[i]["ParallelGateway"]

		// Initialize the indices
		indices := map[string]int{
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

		for j := 0; j < numFields; j++ {

			name := field.Type().Field(j).Name
			typ := field.Field(j).FieldByName("Type").String()
			hash := field.Field(j).FieldByName("Hash").String()
			nextHash := getNextHash(field, j, numFields)
			extName := extractLastTwoWords(name)

			// Use local variables to hold the index values
			startEventIndex := indices["startEventIndex"]
			endEventIndex := indices["endEventIndex"]
			intermediateCatchEventIndex := indices["intermediateCatchEventIndex"]
			intermediateThrowEventIndex := indices["intermediateThrowEventIndex"]
			taskIndex := indices["taskIndex"]
			userTaskIndex := indices["userTaskIndex"]
			scriptTaskIndex := indices["scriptTaskIndex"]
			flowIndex := indices["flowIndex"]
			inclusiveGatewayIndex := indices["inclusiveGatewayIndex"]
			parallelGatewayIndex := indices["parallelGatewayIndex"]

			switch typ {
			case "startevent":
				handleStartEvent(v, i, name, extName, typ, hash, nextHash, numStartEvent, &startEventIndex)
				indices["startEventIndex"] = startEventIndex
			case "event":
				handleEvent(v, i, name, extName, typ, hash, numIntermediateCatchEvent, numIntermediateThrowEvent, numEndEvent, &intermediateCatchEventIndex, &intermediateThrowEventIndex, &endEventIndex)
				indices["intermediateCatchEventIndex"] = intermediateCatchEventIndex
				indices["intermediateThrowEventIndex"] = intermediateThrowEventIndex
				indices["endEventIndex"] = endEventIndex
			case "flow":
				handleFlow(v, i, typ, hash, numFlows, &flowIndex)
				indices["flowIndex"] = flowIndex
			case "gateway":
				handleGateway(v, i, name, extName, typ, hash, numInclusiveGateway, numParallelGateway, &inclusiveGatewayIndex, &parallelGatewayIndex)
			case "activity":
				handleActivity(v, i, name, extName, typ, hash, numTask, numUserTask, numScriptTask, &taskIndex, &userTaskIndex, &scriptTaskIndex)
				indices["taskIndex"] = taskIndex
				indices["userTaskIndex"] = userTaskIndex
				indices["scriptTaskIndex"] = scriptTaskIndex
			}
		}
		log.Printf("indices: %v", indices)
		log.Print("-------------------------")
	}
}
