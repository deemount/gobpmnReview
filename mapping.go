package main

import "reflect"

// mapping ...
type mapping struct {
	Anonym   map[int]string
	Config   map[int]string
	BPMNType map[int]string
}

// Assign maps the fields of a reflectValue to the corresponding maps.
func (m *mapping) Assign(v *reflectValue) {
	// Initialize the maps.
	m.Anonym = make(map[int]string)
	m.Config = make(map[int]string)
	m.BPMNType = make(map[int]string)

	// Define a mapping of field types to their corresponding maps.
	fieldMap := map[reflect.Kind]map[int]string{
		reflect.Bool:   m.Config,
		reflect.Struct: m.BPMNType,
	}

	anonymIndex := 0
	configIndex := 0
	bpmnTypeIndex := 0

	// Iterate over the fields and assign them to the corresponding maps.
	for _, field := range v.Fields {
		if field.Anonymous {
			m.Anonym[anonymIndex] = field.Name
			anonymIndex++
		} else {
			// Get the map for the field type.
			fieldMapForType, ok := fieldMap[field.Type.Kind()]
			if !ok {
				// If the field type is not recognized, check if it's a BPMN struct.
				if field.Type.Name() == "BPMN" {
					m.BPMNType[bpmnTypeIndex] = field.Name
					bpmnTypeIndex++
				}
			} else {
				// Assign the field to the corresponding map.
				fieldMapForType[configIndex] = field.Name
				configIndex++
			}
		}
	}
}
