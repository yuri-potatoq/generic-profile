package enrollment

// Gender enum
type Gender string

const (
	GenderNone   Gender = "None"
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

// EnrollmentShift enum
type EnrollmentShift string

const (
	ShiftNone      EnrollmentShift = "None"
	ShiftMorning   EnrollmentShift = "Morning"
	ShiftAfternoon EnrollmentShift = "Afternoon"
)

// Modalities enum
type Modalities string

const (
	ModalityFootball   Modalities = "Football"
	ModalityBasketball Modalities = "Basketball"
	ModalitySwimming   Modalities = "Swimming"
	ModalityYoga       Modalities = "Yoga"
	ModalityVolleyball Modalities = "Volleyball"
)

// Address struct
type Address struct {
	ZipCode     string `db:"zipCode"`
	Street      string `db:"street"`
	City        string `db:"city"`
	State       string `db:"state"`
	HouseNumber int    `db:"houseNumber"`
}

// ChildProfile struct
type ChildProfile struct {
	FullName    string `db:"fullName"`
	Birthdate   string `db:"birthdate"`
	Gender      Gender `db:"gender"`
	MedicalInfo string `db:"medicalInfo"`
}

// ChildParent struct
type ChildParent struct {
	FullName    string `db:"fullName"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phoneNumber"`
}

// EnrollmentState struct
type EnrollmentState struct {
	ID              int             `db:"id"`
	ChildParent     ChildParent     `db:"childParent"`
	ChildProfile    ChildProfile    `db:"childProfile"`
	Address         Address         `db:"address"`
	Modalities      []Modalities    `db:"modalities"`
	EnrollmentShift EnrollmentShift `db:"enrollmentShift"`
	Terms           bool            `db:"terms"`
}
