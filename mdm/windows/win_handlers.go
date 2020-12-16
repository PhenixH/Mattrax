package windows

import (
	"github.com/gorilla/mux"
)

func (p *Protocol) Mount(r *mux.Router) error {
	r.HandleFunc("/ManagementServer/Manage.svc", Manage(p)).Name("manage").Methods("POST")
	r.HandleFunc("/EnrollmentServer/Enrollment.svc", Enrollment(p)).Name("enrollment").Methods("POST")
	r.HandleFunc("/EnrollmentServer/Policy.svc", Policy(p)).Name("policy").Methods("POST")
	r.HandleFunc("/EnrollmentServer/Discovery.svc", Discover(p)).Name("discovery").Methods("GET", "POST")

	return nil
}
