package v1

import (
	"fmt"
	"github.com/yuri-potatoq/generic-profile/enrollment"
	"net/http"
)

func GetEnrollmentHandler(service enrollment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, rq *http.Request) {
		stt, err := service.GetEnrollmentState(rq.Context(), 0)
		fmt.Println(stt)
		fmt.Println(err)
	}
}
