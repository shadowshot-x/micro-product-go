package authservice

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shadowshot-x/micro-product-go/authservice/data"
	"go.uber.org/zap"
)

var (
	singupRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signup_total",
		Help: "Total number of signup requests",
	})
	signupSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signup_success",
		Help: "Successful signup requests",
	})
	signupFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signup_fail",
		Help: "Failed signup requests",
	})
)

// SignupController is the Signup route handler
type SignupController struct {
	logger            *zap.Logger
	promSignupTotal   prometheus.Counter
	promSignupSuccess prometheus.Counter
	promSignupFail    prometheus.Counter
}

// NewSignupController returns a frsh Signup controller
func NewSignupController(logger *zap.Logger) *SignupController {
	return &SignupController{
		logger:            logger,
		promSignupTotal:   singupRequests,
		promSignupSuccess: signupSuccess,
		promSignupFail:    signupFail,
	}
}

// adds the user to the database of users
func (ctrl *SignupController) SignupHandler(rw http.ResponseWriter, r *http.Request) {
	// we increment the signup request counter
	ctrl.promSignupTotal.Inc()

	// extra error handling should be done at server side to prevent malicious attacks
	if _, ok := r.Header["Email"]; !ok {
		ctrl.logger.Warn("Email was not found in the header")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		ctrl.promSignupFail.Inc()
		return
	}
	if _, ok := r.Header["Username"]; !ok {
		ctrl.logger.Warn("Username was not found in the header")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Username Missing"))
		ctrl.promSignupFail.Inc()
		return
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		ctrl.logger.Warn("Passwordhash was not found in the header")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Passwordhash Missing"))
		ctrl.promSignupFail.Inc()
		return
	}
	if _, ok := r.Header["Fullname"]; !ok {
		ctrl.logger.Warn("Fullname was not found in the header")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Fullname Missing"))
		ctrl.promSignupFail.Inc()
		return
	}

	// validate and then add the user
	check := data.AddUserObject(r.Header["Email"][0], r.Header["Username"][0], r.Header["Passwordhash"][0],
		r.Header["Fullname"][0], 0)
	// if false means username already exists
	if !check {
		ctrl.logger.Warn("User already exists", zap.String("email", r.Header["Email"][0]), zap.String("username", r.Header["Username"][0]))
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Email or Username already exists"))
		ctrl.promSignupFail.Inc()
		return
	}
	ctrl.logger.Info("User created", zap.String("email", r.Header["Email"][0]), zap.String("username", r.Header["Username"][0]))
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
	// this will mean the request was successfully added
	ctrl.promSignupSuccess.Inc()
}
