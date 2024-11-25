package main

import (
	"fmt"
	"log"
	"unsafe"
)

var bpmnModeller BPMNModeller

// ExampleProcess is a struct that represents a process model
// with BPMN elements, with named fields.
type ExampleProcess struct {

	// Def MUST be set to a DefinitionsRepository,
	// otherwise the model will not be valid. It MUST be set at the first place.
	Def DefinitionsRepository // Refers to the DefinitionsRepository

	// If IsExecutable is set, the process is executable and set to true.
	// Otherwise, the process is not executable and set to false.
	// If more than one process in a model is given, only one process can be executable.
	// It will be then the first process given in the model, which will be executable.
	IsExecutable bool // Process Configuration

	// All elements of the BPMN model
	Process        BPMN // BPMN Element
	StartEvent     BPMN // BPMN Element
	FromStartEvent BPMN // BPMN Element
	Task           BPMN // BPMN Element
	FromTask       BPMN // BPMN Element
	EndEvent       BPMN // BPMN Element
}

// RentingProcess is a struct that represents a collaborative process model
// with anonymous fields. It is a composition of two processes.
type (

	// RentingProcess structure
	RentingProcess struct {

		// The first field must be set to a DefinitionsRepository.
		Def DefinitionsRepository // Refers to the DefinitionsRepository

		// The second field must be set to a pool.
		// A pool must have the name Pool in itself to become identified as such and should/can
		// have the name as the process. Only the first process in the model is executable.
		RentingPool // RentingPool represents a collaboration

		// All anonymous fields, which are set after a pool, are considered as processes.
		Tenant   // Tenant represents a process
		Landlord // Landlord represents a process
	}

	// RentingPool ...
	RentingPool struct {
		TenantIsExecutable   bool // Process Configuration
		LandlordIsExecutable bool // Process Configuration
		Collaboration        BPMN // BPMN Element
		TenantProcess        BPMN // BPMN Element
		TenantParticipant    BPMN // BPMN Element
		LandlordProcess      BPMN // BPMN Element
		LandlordParticipant  BPMN // BPMN Element
	}

	// Tenant
	Tenant struct {
		StartEvent     BPMN // BPMN Element
		FromStartEvent BPMN // BPMN Element
		Task           BPMN // BPMN Element
		FromTask       BPMN // BPMN Element
		EndEvent       BPMN // BPMN Element
	}

	// Landlord
	Landlord struct {
		StartEvent     BPMN // BPMN Element
		FromStartEvent BPMN // BPMN Element
		FirstTask      BPMN // BPMN Element
		FromFirstTask  BPMN // BPMN Element
		SecondTask     BPMN // BPMN Element
		FromSecondTask BPMN // BPMN Element
		ScriptTask     BPMN // BPMN Element
		FromScriptTask BPMN // BPMN Element
		EndEvent       BPMN // BPMN Element
	}
)

func main() {

	/*
	 * ExampleProcess
	 */

	exampleProcess := NewReflectDI(ExampleProcess{}).(ExampleProcess)

	// manipulate the incoming and outgoing in a task in the process model

	// set incoming for task
	exampleProcess.Def.GetProcess(0).GetTask(0).SetIncoming(1)

	// dereferencing the pointer
	var incomingHash string = *exampleProcess.Def.GetProcess(0).GetStartEvent(0).GetOutgoing(0).GetFlow()
	exampleProcess.Def.GetProcess(0).GetTask(0).GetIncoming(0).SetFlow(incomingHash)

	// set outgoing for task
	exampleProcess.Def.GetProcess(0).GetTask(0).SetOutgoing(1)

	// dereferencing the pointer
	var outgoingHash string = *exampleProcess.Def.GetProcess(0).GetEndEvent(0).GetIncoming(0).GetFlow()
	exampleProcess.Def.GetProcess(0).GetTask(0).GetOutgoing(0).SetFlow(outgoingHash)

	// some logging messages for the terminal
	log.Printf("exampleProcess.Target: %+#v", exampleProcess) // represents the Target
	log.Print("-------------------------")
	log.Printf("exampleProcess.Def: %+#v", exampleProcess.Def) // represents the model to create
	log.Print("-------------------------")
	fmt.Printf("Size: %d\n", unsafe.Sizeof(exampleProcess))

	// create the process model
	example, err := NewBPMNModeller(WithCounter(), WithDefinitions(exampleProcess.Def))
	if err != nil {
		panic(err)
	}
	example.Marshal()

	/*
	 * RentalProcess
	 */

	rentalProcess := NewReflectDI(RentingProcess{}).(RentingProcess)

	// some logging messages for the terminal
	log.Printf("rentalProcess.Target: %+#v", rentalProcess) // represents the Target
	log.Print("-------------------------")
	log.Printf("rentalProcess.Def: %+#v", rentalProcess.Def) // represents the model to create
	log.Print("-------------------------")
	fmt.Printf("Size: %d\n", unsafe.Sizeof(rentalProcess))

	rental, err := NewBPMNModeller(WithCounter(), WithDefinitions(rentalProcess.Def))
	if err != nil {
		panic(err)
	}
	rental.Marshal()

}
