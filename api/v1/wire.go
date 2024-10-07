package v1

import "github.com/google/wire"

var WireSet = wire.NewSet(NewPatchEnrollmentHandler, NewGetEnrollmentHandler)
