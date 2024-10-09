package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/yuri-potatoq/generic-profile/enrollment"
	"net/http"
	"strconv"
	"time"
)

type AddressDTO struct {
	ZipCode     string `json:"zipcode"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	HouseNumber int    `json:"number"`
}

type ChildProfileDTO struct {
	FullName    string            `json:"fullName"`
	Birthdate   time.Time         `json:"birthdate"`
	Gender      enrollment.Gender `json:"gender"`
	MedicalInfo string            `json:"medicalInfo"`
}

type ChildParentDTO struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type EnrollmentStateDTO struct {
	ID              *int                        `json:"id,omitempty"`
	ChildParent     *ChildParentDTO             `json:"childParent,omitempty"`
	ChildProfile    *ChildProfileDTO            `json:"childProfile,omitempty"`
	Address         *AddressDTO                 `json:"address,omitempty"`
	Modalities      []enrollment.Modalities     `json:"modalities,omitempty"`
	EnrollmentShift *enrollment.EnrollmentShift `json:"enrollmentShift,omitempty"`
	Terms           *bool                       `json:"terms,omitempty"`
}

func MapEnrollmentStateResponse(stt enrollment.EnrollmentState) *EnrollmentStateDTO {
	return &EnrollmentStateDTO{
		ID: &stt.ID,
		Address: &AddressDTO{
			ZipCode:     stt.Address.ZipCode,
			State:       stt.Address.State,
			City:        stt.Address.City,
			HouseNumber: stt.Address.HouseNumber,
			Street:      stt.Address.Street,
		},
		ChildParent: &ChildParentDTO{
			FullName:    stt.ChildParent.FullName,
			Email:       stt.ChildParent.Email,
			PhoneNumber: stt.ChildParent.PhoneNumber,
		},
		ChildProfile: &ChildProfileDTO{
			FullName:    stt.ChildProfile.FullName,
			MedicalInfo: stt.ChildProfile.MedicalInfo,
			Gender:      stt.ChildProfile.Gender,
			Birthdate:   stt.ChildProfile.Birthdate,
		},
		Modalities:      stt.Modalities,
		EnrollmentShift: &stt.EnrollmentShift,
	}
}

type GetEnrollmentHandler struct {
	s enrollment.Service
}

func NewGetEnrollmentHandler(s enrollment.Service) *GetEnrollmentHandler {
	return &GetEnrollmentHandler{s}
}

func (h *GetEnrollmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// TODO: handle error properly

	stt, err := h.s.GetEnrollmentState(r.Context(), id)
	if err != nil {
		//TODO: write better errors
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(MapEnrollmentStateResponse(stt))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

type PatchEnrollmentHandler struct {
	s enrollment.Service
}

func NewPatchEnrollmentHandler(s enrollment.Service) *PatchEnrollmentHandler {
	return &PatchEnrollmentHandler{s}
}

func (h *PatchEnrollmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body EnrollmentStateDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}

	stt, err := h.s.BulkUpdate(r.Context(), ToPartialUpdate(body))
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(MapEnrollmentStateResponse(stt))
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
}

func ToPartialUpdate(body EnrollmentStateDTO) enrollment.PartialUpdate {
	var opts []enrollment.PartialUpdateOpt
	if body.ChildProfile != nil {
		opts = append(opts, enrollment.UpdateWithChildProfile(enrollment.ChildProfile{
			FullName:    body.ChildProfile.FullName,
			MedicalInfo: body.ChildProfile.MedicalInfo,
			Gender:      body.ChildProfile.Gender,
			Birthdate:   body.ChildProfile.Birthdate,
		}))
	}
	if body.Address != nil {
		opts = append(opts, enrollment.UpdateWithAddress(enrollment.Address{
			Street:      body.Address.Street,
			ZipCode:     body.Address.ZipCode,
			City:        body.Address.City,
			HouseNumber: body.Address.HouseNumber,
			State:       body.Address.State,
		}))
	}
	if body.ChildParent != nil {
		opts = append(opts, enrollment.UpdateWithChildParent(enrollment.ChildParent{
			Email:       body.ChildParent.Email,
			FullName:    body.ChildParent.FullName,
			PhoneNumber: body.ChildParent.PhoneNumber,
		}))
	}
	if body.EnrollmentShift != nil {
		opts = append(opts, enrollment.UpdateWithShift(*body.EnrollmentShift))
	}
	if len(body.Modalities) > 0 {
		opts = append(opts, enrollment.UpdateWithModalities(body.Modalities))
	}
	if body.Terms != nil {
		opts = append(opts, enrollment.UpdateWithTerm(*body.Terms))
	}

	partial := enrollment.NewPartialUpdate(opts...)
	partial.ID = body.ID

	return partial
}
