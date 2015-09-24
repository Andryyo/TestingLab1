// TestingLab1 project main.go
package main

import (
	"fmt"
	"math/rand"
)

type System struct {
	components map[string]Component
	function Conditioner
}

func (s *System) generateFailVector(requiredErrorsCount int) FailVector {
	result := FailVector{}
	keys := make([]string, 0, len(s.components))
	for key := range s.components {
		keys = append(keys, key)
		result.failed[key] = false
	}
	errorsCount := 0;
	for errorsCount < requiredErrorsCount {
		index := rand.Intn(len(keys))
		if !result.failed[keys[index]] {
			result.failed[keys[index]] = false
			errorsCount++
		}
	}
	return result
}

type Component struct {
	name string
	failProbability float64
	working bool
}

type Processor struct {
	Component
	load int
	maxLoad int
}

type FailVector struct {
	failed map[string]bool
}

type Conditioner interface {
	check(map[string]Component) bool 
}

type Condition struct {
	name string
}

func NewCondition(name string) Conditioner {
	return Condition{name}
}

func (c Condition) check(components map[string]Component) bool {
	return components[c.name].working
}

type AndCondition struct {
	conditions []Conditioner
}

func (c AndCondition) check(components map[string]Component) bool {
	for _, condition := range c.conditions {
		if !condition.check(components) {
			return false
		}
	}
	return true
}

func NewAndCondition(conditions ...Conditioner) Conditioner {
	result := AndCondition{}
	for _, condition := range conditions {
		result.conditions = append(result.conditions, condition)
	}
	return result
}

type OrCondition struct {
	conditions []Conditioner
}

func (c OrCondition) check(components map[string]Component) bool {
	for _, condition := range c.conditions {
		if condition.check(components) {
			return true
		}
	}
	return false
}

func NewOrCondition(conditions ...Conditioner) Conditioner {
	result := OrCondition{}
	for _, condition := range conditions {
		result.conditions = append(result.conditions, condition)
	}
	return result
}

func NewComponent(name string, failProbability float64) Component {
	return Component{name, failProbability, true}
}

func NewSystem() *System {
	result := &System{}
	result.components = map[string]Component {
		"Pr1" : NewComponent("Pr1", 1.2*0.0001),
		"Pr2" : NewComponent("Pr2", 1.2*0.0001),
		"Pr3" : NewComponent("Pr3", 1.2*0.0001),
		"Pr5" : NewComponent("Pr5", 1.2*0.0001),
		"Pr6" : NewComponent("Pr6", 1.2*0.0001),
		"A1" : NewComponent("A1", 1.2*0.0001),
		"A2" : NewComponent("A2", 1.2*0.0001),
		"B1" : NewComponent("B1", 1.5*0.00001),
		"B2" : NewComponent("B2", 1.5*0.00001),
		"B4" : NewComponent("B4", 1.5*0.00001),
		"B5" : NewComponent("B5", 1.5*0.00001),
		"C1" : NewComponent("C1", 4.1*0.0001),
		"C2" : NewComponent("C2", 4.1*0.0001),
		"C4" : NewComponent("C4", 4.1*0.0001),
		"C5" : NewComponent("C5", 4.1*0.0001),
		"C6" : NewComponent("C6", 4.1*0.0001),
		"D1" : NewComponent("D1", 2.2*0.00001),
		"D2" : NewComponent("D2", 2.2*0.00001),
		"D4" : NewComponent("D4", 2.2*0.00001),
		"D6" : NewComponent("D6", 2.2*0.00001),
		"D8" : NewComponent("D8", 2.2*0.00001),
		"M1" : NewComponent("M1", 3.6*0.0001),
		"M2" : NewComponent("M2", 3.6*0.0001),
	}
	result.function = NewAndCondition(
		NewAndCondition(
			NewOrCondition(NewCondition("Pr1"),NewCondition("Pr2")),
			NewOrCondition(NewCondition("B1"),NewCondition("B2")),
			NewCondition("A1"),
			NewOrCondition(NewCondition("M1"),NewCondition("M2")),
			NewCondition("A2"),
			NewOrCondition(NewCondition("B4"),NewCondition("B5")),
			NewOrCondition(NewCondition("C5"),NewCondition("C6")),
			NewCondition("D8"),
			),
		NewAndCondition(),
		NewAndCondition(),
		NewAndCondition(),
	)
	return result
}

func main() {
	fmt.Println("Hello World!")
}
