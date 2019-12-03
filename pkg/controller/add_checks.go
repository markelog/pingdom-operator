package controller

import (
	"github.com/markelog/pingdom-operator/pkg/controller/checks"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, checks.Add)
}
