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

	// Create a new reflectValue
	v := newReflectValue(p)

	// Get the number of fields of the reflectValue.
	v.TargetNumField = v.Target.NumField()

	// Set the default field of the reflectValue.
	v.Def = v.Target.FieldByName("Def")
	if !v.Def.IsValid() {
		fmt.Println("Def field is not valid")
		return nil
	}

	// Set the default attributes if the field is settable.
	// Note: I use this condition for a idea, if the v.Def can't set and needs to be a) set or b) initialized automatically right now.
	// Othwise, a better check is above, which fullfill a validation check. It works contrary to this condition.
	if v.Def.CanSet() {
		definitions := NewDefinitions() // Note: I'm not sure if I should use the NewDefinitions function here. Maybe reflecting it?
		definitions.SetDefaultAttributes()
		v.Def.Set(reflect.ValueOf(definitions))
	}

	// Create a new mapping and quantity
	m := new(mapping) // Note: Mapping and Quantity is a helper structure and can be put together in a single structure?
	q := new(quantity)

	// Assign the fields of the reflectValue to the corresponding maps.
	m.Assign(v)

	if len(m.Anonym) > 0 {
		v.handlePool(q, m)
	} else {
		v.handleSingle(q, m)
	}

	v.reflectProcess(q)

	di := v.target(q, m) // Note: any solutions with generics here?
	log.Printf("quantity: %+v", q)
	log.Print("-------------------------")

	v.process(q) // Note: with the maps it is possible to build the collaboration after the processes
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
			v.anonym(anonymField, q)
		}
	} else {
		v.nonAnonym(q, m)
	}
	return v.Target.Interface()
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
func (v *reflectValue) process(q *quantity) {
	v.ProcessExec = make([]bool, q.Process)
	v.ProcessType = make([]string, q.Process)
	v.ProcessHash = make([]string, q.Process)
	v.ParticipantHash = make([]string, q.Participant)
	processor := NewElementProcessor(v, q)
	if q.Pool > 0 {
		v.multipleProcess(q)
		processor.ProcessElements()
	} else {
		v.singleProcess(q)
		processor.ProcessSingleElement()
	}
}

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
		f := v.Target.FieldByName(field)
		fName, _ := v.Target.Type().FieldByName(field)
		typ := typ(fName.Name)
		hash, _ := hash(typ)
		f.Set(reflect.ValueOf(hash))
	}
}

// singleProcess ...
// (Note: this method is used in a single process and not in usage. Do I need it?)
func (v *reflectValue) singleProcess(q *quantity) {
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

	// get the methods
	methodCalls := methods() // Note: look at internal.go

	for _, call := range methodCalls {
		method := v.Process[0].MethodByName(call.name)
		if method.IsValid() && call.arg > 0 {
			method.Call([]reflect.Value{reflect.ValueOf(call.arg)})
		}
	}

}

// anonym sets the fields in the BPMN model.
func (v *reflectValue) anonym(field string, q *quantity) {
	targetFieldByName := v.Target.FieldByName(field)
	targetNumField := targetFieldByName.NumField()
	for i := 0; i < targetNumField; i++ {
		name := targetFieldByName.Type().Field(i).Name
		switch targetFieldByName.Field(i).Kind() {
		case reflect.Bool:
			v.config(name, i, targetFieldByName)
		case reflect.Struct:
			//q.countFlow(name)    // Note: should I count here?
			//q.countElement(name) // Note: should I count here?
			v.current(i, targetFieldByName)
			v.next(i, targetFieldByName)
		}
	}
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

// multipleProcess processes multiple processes in the BPMN model.
func (v *reflectValue) multipleProcess(q *quantity) {
	l, j, n := 0, 0, 0
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

}

// handlePool assigns the field Pool to a reflected value.
// (Note: walk through m.Anonym  is maybe redudant)
func (v *reflectValue) handlePool(q *quantity, m *mapping) {
	for _, anonymField := range m.Anonym {
		if strings.Contains(anonymField, "Pool") {
			v.Pool = v.Target.FieldByName(anonymField)
			q.countFieldsInPool(v)
			q.countFieldsInProcess(v)
			q.Pool++
			break
		}
	}
}

// handleSingle ...
func (v *reflectValue) handleSingle(q *quantity, m *mapping) {
	for _, bpmnType := range m.BPMNType {
		if strings.Contains(bpmnType, "Process") {
			v.ProcessName = append(v.ProcessName, v.Name)
			q.countFieldsInProcess(v)
			q.Process++
			break
		}
	}
}
