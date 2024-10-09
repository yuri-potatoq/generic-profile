package enrollment

import (
	"github.com/yuri-potatoq/generic-profile/infra/db"
	"time"
)

// Gender enum
type Gender string

const (
	GenderNone   Gender = "None"
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

func (m *Gender) Scan(src interface{}) error {
	return db.ScanEnum[Gender](m, "Gender")(src)
}

// EnrollmentShift enum
type EnrollmentShift string

const (
	ShiftNone      EnrollmentShift = "None"
	ShiftMorning   EnrollmentShift = "Morning"
	ShiftAfternoon EnrollmentShift = "Afternoon"
)

func (m *EnrollmentShift) Scan(src interface{}) error {
	return db.ScanEnum[EnrollmentShift](m, "EnrollmentShift")(src)
}

// Modalities enum
type Modalities string

const (
	ModalityFootball   Modalities = "Football"
	ModalityBasketball Modalities = "Basketball"
	ModalitySwimming   Modalities = "Swimming"
	ModalityYoga       Modalities = "Yoga"
	ModalityVolleyball Modalities = "Volleyball"
)

func (m *Modalities) Scan(src interface{}) error {
	return db.ScanEnum[Modalities](m, "Modalities")(src)
}

// Address struct
type Address struct {
	ZipCode     string `db:"zipcode"`
	Street      string `db:"street"`
	City        string `db:"city"`
	State       string `db:"state"`
	HouseNumber int    `db:"house_number"`
}

// ChildProfile struct
type ChildProfile struct {
	FullName    string    `db:"full_name"`
	Birthdate   time.Time `db:"birthdate"`
	Gender      Gender    `db:"gender"`
	MedicalInfo string    `db:"medical_info"`
}

// ChildParent struct
type ChildParent struct {
	FullName    string `db:"full_name"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
}

// EnrollmentState struct
type EnrollmentState struct {
	ID              int
	ChildParent     ChildParent
	ChildProfile    ChildProfile
	Address         Address
	Modalities      []Modalities
	EnrollmentShift EnrollmentShift
	Terms           bool
}
