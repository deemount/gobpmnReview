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

	// Create a new mapping
	m := new(mapping) // Note: Mapping and Quantity is a helper structure and can be put together in a single structure?
	// Assign the fields of the reflectValue to the corresponding maps.
	m.Assign(v)

	// Create a new quantity
	q := new(quantity)

	// Count the anonymous fields in the BPMN model.
	if len(m.Anonym) > 0 {
		v.handlePool(q, m)
	} else {
		v.handleSingle(q, m)
	}

	// Reflect the processes in the BPMN model by given quantity.
	// v.Process[] is a slice of reflect.Value, which represents the processes in the BPMN model.
	v.reflectProcess(q)

	// Get the target of the reflectValue.
	// target holds the data structure of the BPMN model in the reflectValue.
	target := v.target(q, m)

	//
	v.process(q) // Note: with the maps it is possible to build the collaboration after the processes
	v.collaboration(q)

	return target
}

// reflectValue ...
type reflectValue struct {

	// General Configuration
	Name   string
	Fields []reflect.StructField

	// BPMN Configuration
	Target         reflect.Value
	TargetNumField int
	Def            reflect.Value
	Pool           reflect.Value

	// Process Configuration
	Process     []reflect.Value
	ProcessName []string
	ProcessType []string
	ProcessHash []string
	ProcessExec []bool

	// Participant Configuration
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

// handlePool assigns the field Pool to a reflected value.
// Note: walk through m.Anonym is maybe redudant
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

// target ...
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
	for i := 0; i < q.Process; i++ {
		v.Process[i] = v.Def.MethodByName("GetProcess").Call([]reflect.Value{reflect.ValueOf(i)})[0]
	}
}

// process sets the maps and process parameters in the BPMN model.
func (v *reflectValue) process(q *quantity) {

	// initialize the slices
	v.ProcessExec = make([]bool, q.Process)
	v.ProcessType = make([]string, q.Process)
	v.ProcessHash = make([]string, q.Process)
	v.ParticipantHash = make([]string, q.Participant)

	// process each process name
	processor := NewElementProcessor(v, q)

	// if the pool is greater than 0, then process multiple processes
	if q.Pool > 0 {
		// process multiple processes
		err := v.configurePool()
		if err != nil {
			log.Printf("Error configuring multiple process: %v", err)
		}
		v.multipleProcess(q)
		processor.ProcessMultipleElement()
	} else {

		// process a single process
		err := v.configureSingleProcess()
		if err != nil {
			log.Printf("Error processing single process: %v", err)
		}
		v.singleProcess(0, q)
		processor.ProcessSingleElement()
	}
}

// nonAnonym sets the IsExecutable field to true and populates the reflection fields with hash values.
// (Note: the method is used in a single process).
func (v *reflectValue) nonAnonym(q *quantity, m *mapping) {
	v.setIsExecutableByField(m.Config)
	v.populateReflectionFields(q, m.BPMNType)
}

// setIsExecutableByField configures the process executable status.
// It sets the IsExecutable field to true.
// Note: maybe redudant; look at method config.
func (v *reflectValue) setIsExecutableByField(fields map[int]string) {
	for _, field := range fields {
		f := v.Target.FieldByName(field)
		f.SetBool(true)
		break
	}
}

// setIsExecutableByMethod configures the process executable status.
// It calls SetIsExecutable method on the process.
func (v *reflectValue) setIsExecutableByMethod(process reflect.Value, isExecutable bool) error {
	method := process.MethodByName("SetIsExecutable")
	if !method.IsValid() {
		return fmt.Errorf("SetIsExecutable method not found")
	}
	method.Call([]reflect.Value{reflect.ValueOf(isExecutable)})
	return nil
}

// populateReflectionFields populates the reflection fields with hash values.
// Note: this method is used in a single process and in usage in the method "nonAnonym".
func (v *reflectValue) populateReflectionFields(q *quantity, reflectionFields map[int]string) {
	for _, field := range reflectionFields {
		f := v.Target.FieldByName(field)
		fName, _ := v.Target.Type().FieldByName(field)
		typ := typ(fName.Name)
		hash, _ := hash(typ)
		f.Set(reflect.ValueOf(hash))
	}
}

// setProcessIdentifiers sets the process ID and name
func (v *reflectValue) setProcessIdentifiers(process reflect.Value, name string, field reflect.Value) error {
	hash := field.FieldByName("Hash").String()

	// Set Process ID
	if err := callProcessMethod(process, "SetID",
		[]reflect.Value{
			reflect.ValueOf(strings.ToLower(name)),
			reflect.ValueOf(hash),
		}); err != nil {
		return fmt.Errorf("failed to set process ID: %w", err)
	}

	// Set Process Name
	if err := callProcessMethod(process, "SetName",
		[]reflect.Value{reflect.ValueOf(name)}); err != nil {
		return fmt.Errorf("failed to set process name: %w", err)
	}

	return nil
}

// singleProcess configures and initializes a single BPMN process.
// It sets up the process properties and applies the required BPMN elements.
// Note: this method is used for a single process and in usage in the method "process".
func (v *reflectValue) singleProcess(processIdx int, q *quantity) error {
	if err := v.configureSingleProcess(); err != nil {
		return fmt.Errorf("failed to configure process: %w", err)
	}

	if err := v.applyProcessMethods(processIdx, q); err != nil {
		return fmt.Errorf("failed to apply process methods: %w", err)
	}

	return nil
}

// configureSingleProcess configures a single BPMN process.
func (v *reflectValue) configureSingleProcess() error {

	if v.Process == nil || len(v.Process) == 0 {
		return fmt.Errorf("invalid process: Process slice is nil or empty")
	}

	process := v.Process[0]

	for i := 0; i < v.TargetNumField; i++ {
		field := v.Target.Field(i)
		fieldType := v.Target.Type().Field(i)

		switch {
		case strings.Contains(fieldType.Name, "IsExecutable"):
			if err := v.setIsExecutableByMethod(process, field.Bool()); err != nil {
				return err
			}
		case fieldType.Name == "Process":
			if err := v.setProcessIdentifiers(process, fieldType.Name, field); err != nil {
				return err
			}
		}

	}

	return nil

}

// ProcessMethod represents a BPMN process method with its argument
type ProcessMethod struct {
	Name string
	Arg  int
}

// GetProcessMethods returns the standard BPMN process methods with their quantities.
// Each method represents a BPMN element type and its required quantity.
// Note: q quantity for a single process should be used in the Arg field.
func GetProcessMethods(processIdx int, q *quantity) []ProcessMethod {
	process := q.ProcessElements[processIdx]
	return []ProcessMethod{
		{Name: "SetStartEvent", Arg: process["StartEvent"]},
		{Name: "SetEndEvent", Arg: process["EndEvent"]},
		{Name: "SetIntermediateCatchEvent", Arg: process["IntermediateCatchEvent"]},
		{Name: "SetIntermediateThrowEvent", Arg: process["IntermediateThrowEvent"]},
		{Name: "SetInclusiveGateway", Arg: process["InclusiveGateway"]},
		{Name: "SetExclusiveGateway", Arg: process["ExclusiveGateway"]},
		{Name: "SetParallelGateway", Arg: process["ParallelGateway"]},
		{Name: "SetUserTask", Arg: process["UserTask"]},
		{Name: "SetScriptTask", Arg: process["ScriptTask"]},
		{Name: "SetTask", Arg: process["Task"]},
		{Name: "SetSequenceFlow", Arg: process["SequenceFlow"]},
	}
}

// callProcessMethod is a helper function to call process methods safely
func callProcessMethod(process reflect.Value, methodName string, args []reflect.Value) error {
	method := process.MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("method %s not found", methodName)
	}
	method.Call(args)
	return nil
}

// applyProcessMethods applies all BPMN element methods to the process
func (v *reflectValue) applyProcessMethods(processIdx int, q *quantity) error {
	methods := GetProcessMethods(processIdx, q)

	for _, methodCall := range methods {
		if methodCall.Arg <= 0 {
			continue
		}

		if err := callProcessMethod(v.Process[processIdx], methodCall.Name,
			[]reflect.Value{reflect.ValueOf(methodCall.Arg)}); err != nil {
			return fmt.Errorf("failed to apply method %s: %w", methodCall.Name, err)
		}
	}

	return nil
}

/*
 * @ multiple processes
 */
// anonym sets the fields in the BPMN model.
func (v *reflectValue) anonym(f string, q *quantity) {
	targetField := v.Target.FieldByName(f) // must be a struct, which represents a process
	targetNum := targetField.NumField()    // get the number of fields in the struct
	for i := 0; i < targetNum; i++ {
		name := targetField.Type().Field(i).Name
		switch targetField.Field(i).Kind() {
		case reflect.Bool:
			v.config(name, i, targetField)
		case reflect.Struct:
			v.currentHash(i, targetField)
			v.nextHash(i, targetField)
		}
	}
}

// config sets the IsExecutable field to true if the name contains "IsExecutable" and index is 0.
// I called it config, because it is a configuration of the field.
// Note: the method is used in a multiple process and it's a third option call for the field IsExecutable.
func (v *reflectValue) config(name string, index int, target reflect.Value) {
	if strings.Contains(name, "IsExecutable") && index == 0 {
		target.Field(0).SetBool(true)
	}
}

// currentHash sets the hash value of the current field if it is empty.
func (v *reflectValue) currentHash(index int, target reflect.Value) {
	h := fmt.Sprintf("%s", target.Field(index).FieldByName("Hash"))
	if h == "" {
		typ := typ(target.Type().Field(index).Name)
		hash, _ := hash(typ)
		target.Field(index).Set(reflect.ValueOf(hash))
	}
}

// nextHash sets the hash value of the next field.
func (v *reflectValue) nextHash(index int, target reflect.Value) {
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
func (v *reflectValue) multipleProcess(q *quantity) error {

	if err := v.configurePool(); err != nil {
		return fmt.Errorf("failed to configure multiple processes: %w", err)
	}

	for processIdx := 0; processIdx < q.Process; processIdx++ {

		process := v.Process[processIdx]
		isExecutable := v.ProcessExec[processIdx]
		elements := q.ProcessElements[processIdx]
		typ := v.ProcessType[processIdx]
		hash := v.ProcessHash[processIdx]
		name := v.ProcessName[processIdx]

		if processIdx == 0 {
			if err := v.setIsExecutableByMethod(process, isExecutable); err != nil {
				return fmt.Errorf("failed to set process executable: %w", err)
			}
		}

		// Set the process ID and name
		process.MethodByName("SetID").Call([]reflect.Value{reflect.ValueOf(typ), reflect.ValueOf(hash)})
		process.MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf(name)})

		// Set the process elements
		for method, arg := range elements {
			methodName := fmt.Sprintf("Set%s", method)
			process.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(arg)})
		}
	}

	return nil

}

// configurePool configures multiple BPMN processes.
func (v *reflectValue) configurePool() error {

	if v.Pool.NumField() == 0 {
		return fmt.Errorf("invalid pool: Pool is empty")
	}

	l, j, n := 0, 0, 0
	for i := 0; i < v.Pool.NumField(); i++ {
		name := v.Pool.Type().Field(i).Name
		field := v.Pool.Field(i)

		switch {
		case strings.Contains(name, "IsExecutable"):
			if field.IsValid() {
				v.ProcessExec[j] = field.Bool()
				j++
			}
		case strings.Contains(name, "Process"):
			if field.IsValid() {
				v.ProcessHash[l] = field.FieldByName("Hash").String()
				v.ProcessType[l] = field.FieldByName("Type").String()
				l++
			}
		case strings.Contains(name, "Participant"):
			if field.IsValid() {
				v.ParticipantHash[n] = field.FieldByName("Hash").String()
				n++
			}
		}

	}

	return nil

}
