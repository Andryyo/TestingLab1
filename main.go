// TestingLab1 project main.go
package main

import (
	"log"
	"fmt"
	"math/rand"
	"time"
)

type System struct {
	components map[string]Component
	function AndCondition
}

func (s *System) NewFailVector() FailVector {
	result := FailVector{}
	result.failed = make(map[string]bool)
	result.probability = 1
	for key := range s.components {
		result.failed[key] = false
		result.probability *= 1 - s.components[key].failProbability
	}
	return result
}

func (s *System) corrupt(failVector FailVector) {
	for key, val := range failVector.failed {
		component := s.components[key]
		component.working = !val
		s.components[key] = component
	}
}

func (s *System) repair() {
	for _, component := range s.components {
		component.working = true
	}
}

func (s *System) calcFailProbability(errorsCount int) {
	coverage := 100
	switch errorsCount {
		case 0, 1, 2: coverage = 100
		case 3: coverage = 50
		case 4: coverage = 25
		default: coverage = 100
	}
	failVectors := s.generateRandomFailVectors(errorsCount, coverage)
	var result float64 = 0
	failsCount := 0
	totalCount := 0
	for _, failVector := range failVectors {
		s.corrupt(failVector)
		resultString := fmt.Sprintf("%v: ", failVector)
		for _, function := range s.function.conditions {
			resultString += fmt.Sprintf("%v, ", function.check(s.components))
		}
		//log.Println(resultString)
		if !s.function.check(s.components) {
			result += failVector.probability
			failsCount ++
		}
		totalCount ++
		s.repair()
	} 
	result = result * float64(coverage) / 100
	log.Printf("Fail probability at %v errors: %v, %v\\%v\n", errorsCount, result, failsCount, totalCount)
}

func (s *System) generateRandomFailVectors(requiredErrorsCount int, coveragePercentage int) []FailVector {
	//if coveragePercentage == 100 {
	//	return s.generateAllFailVectors(requiredErrorsCount)
	//}
	componentsCount := len(s.components)
	resultsCount := 1
	for i := 0; i<requiredErrorsCount; i++ {
		resultsCount *= componentsCount - i
	}
	for i := 1; i<=requiredErrorsCount; i++ {
		resultsCount /= i
	}
	resultsCount = resultsCount * coveragePercentage / 100
	results := make([]FailVector, 0, resultsCount)
	keys := make([]string, 0, len(s.components))
	for key := range s.components {
		keys = append(keys, key)
	}
	for len(results) < resultsCount {
		failVector := s.NewFailVector()
		errorsKeys := make([]string, 0, requiredErrorsCount)
		for len(errorsKeys) < requiredErrorsCount {
			key := keys[rand.Intn(componentsCount)]
			for _, errorKey := range errorsKeys {
				if errorKey == key {
					continue
				}
			}
			errorsKeys = append(errorsKeys, key)
		}
		for _, errorKey := range errorsKeys {
			failVector.failed[errorKey] = true
			failVector.probability /= 1-s.components[errorKey].failProbability
			failVector.probability *= s.components[errorKey].failProbability
		}
		results = append(results, failVector)
	}
	return results
}

// !! Function not working
func (s *System) generateAllFailVectors(requiredErrorsCount int) []FailVector {
	if requiredErrorsCount == 0 {
		return []FailVector{s.NewFailVector()}
	}
	results := make([]FailVector, 0)
	tmp := s.generateAllFailVectors(requiredErrorsCount - 1)
	for key := range s.components {
		for _, failVector := range tmp {
			if failVector.failed[key] == false {
				result := s.NewFailVector()
				for k := range s.components {
					result.failed[k] = failVector.failed[k]
				}
				result.probability = failVector.probability
				result.failed[key] = true
				result.probability /= 1-s.components[key].failProbability
				result.probability *= s.components[key].failProbability
				results = append(results, result)
			}
		}
	}
	return results
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
	probability float64
}

func (f FailVector) String() string {
	result := ""
	for key, value := range f.failed {
		if value {
			result += fmt.Sprintf("%v failed, ", key) 
		}
	}
	result += fmt.Sprintf("%v", f.probability)
	return result
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

func NewAndCondition(conditions ...Conditioner) AndCondition {
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
			NewCondition("B4"),
			NewOrCondition(NewCondition("C5"),NewCondition("C6")),
			NewCondition("D8"),
			),
		NewAndCondition(
			NewCondition("Pr3"),
			NewOrCondition(NewCondition("B1"),NewCondition("B2")),
			NewCondition("A1"),
			NewCondition("M1"),
			NewCondition("C4"),
			NewCondition("D6"),
			),
		NewAndCondition(
			NewCondition("Pr6"),
			NewOrCondition(NewCondition("B4"),NewCondition("B5")),
			NewCondition("A2"),
			NewOrCondition(NewCondition("M1"),NewCondition("M2")),
			NewCondition("A1"),
			NewOrCondition(NewCondition("B1"),NewCondition("B2")),
			NewCondition("C1"),
			NewOrCondition(NewCondition("D1"),NewCondition("D2")),
			),
		NewAndCondition(
			NewCondition("Pr5"),
			NewCondition("B4"),
			NewCondition("A2"),
			NewOrCondition(NewCondition("M1"),NewCondition("M2")),
			NewCondition("A1"),
			NewOrCondition(NewCondition("B1"),NewCondition("B2")),
			NewOrCondition(NewCondition("C1"),NewCondition("C2")),
			NewCondition("D2"),
			),
	)
	return result
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewSystem()
	s.calcFailProbability(0)
	s.calcFailProbability(1)
	s.calcFailProbability(2)
	s.calcFailProbability(3)
	s.calcFailProbability(4)
}
