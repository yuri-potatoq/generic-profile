package enrollment

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(NewEnrollmentRepository, NewEnrollmentService)
