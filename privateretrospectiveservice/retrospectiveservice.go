package retrospectiveservice

import (
	"fmt"
	"net/http"

	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"
)

// we will control this with mutex
var retrospective string

func handleNotInHeader(rw http.ResponseWriter, r *http.Request, param string) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(fmt.Sprintf("%s Missing", param)))
}

// RetrospectiveController is the Real Time Socket route handler
type RetrospectiveController struct {
	logger                      *zap.Logger
	melodyInstance              *melody.Melody
	access                      sync.Mutex
	currentUser                 string
	status                      bool
	retrospectiveAvailSuccess   prometheus.Counter
	retrospectiveAvailFail      prometheus.Counter
	retrospectiveUpdateSuccess  prometheus.Counter
	retrospectiveUpdateFail     prometheus.Counter
	retrospectiveReleaseSuccess prometheus.Counter
	retrospectiveReleaseFail    prometheus.Counter
}

// NewRetrospectiveController returns a fresh Retrospective controller
func NewRetrospectiveController(logger *zap.Logger, melodyInstance *melody.Melody, retrospectiveAvailSuccess prometheus.Counter, retrospectiveAvailFail prometheus.Counter, retrospectiveUpdateSuccess prometheus.Counter,
	retrospectiveUpdateFail prometheus.Counter, retrospectiveReleaseSuccess prometheus.Counter, retrospectiveReleaseFail prometheus.Counter) *RetrospectiveController {
	return &RetrospectiveController{
		logger:                      logger,
		melodyInstance:              melodyInstance,
		status:                      false,
		retrospectiveAvailSuccess:   retrospectiveAvailSuccess,
		retrospectiveAvailFail:      retrospectiveAvailFail,
		retrospectiveUpdateSuccess:  retrospectiveUpdateSuccess,
		retrospectiveUpdateFail:     retrospectiveUpdateFail,
		retrospectiveReleaseSuccess: retrospectiveReleaseSuccess,
		retrospectiveReleaseFail:    retrospectiveReleaseFail,
	}
}

// Assigns the Websocket to the request.
func (ctrl *RetrospectiveController) AssignSocket(rw http.ResponseWriter, r *http.Request) {
	ctrl.melodyInstance.HandleRequest(rw, r)
}

func (ctrl *RetrospectiveController) BroadcastMessage(rw http.ResponseWriter, r *http.Request) {
	// lets broadcast to everyone
	ctrl.melodyInstance.Broadcast([]byte(retrospective))
}

func (ctrl *RetrospectiveController) CheckAccess(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	if !ctrl.status {
		rw.Write([]byte("false"))
	} else {
		rw.Write([]byte("true"))
	}
}

func (ctrl *RetrospectiveController) AvailRetrospective(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Username"]; !ok {
		ctrl.logger.Warn("User Name was not found in the header")

		// increment retrospectiveAvailFail prometheus metric here
		ctrl.retrospectiveAvailFail.Inc()

		handleNotInHeader(rw, r, "user")
		return
	}
	if !ctrl.status {
		ctrl.logger.Warn("Tried to avail the locked retrospective")

		// increment retrospectiveAvailFail prometheus metric here
		ctrl.retrospectiveAvailFail.Inc()

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Access to Retrospective Denied"))
		return
	}
	ctrl.status = true
	ctrl.access.Lock()
	ctrl.currentUser = r.Header["Username"][0]

	// increment retrospectiveAvailSuccess prometheus metric here
	ctrl.retrospectiveAvailSuccess.Inc()

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Access to Retrospective Granted"))
}

func (ctrl *RetrospectiveController) ChangeRetrospective(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Username"]; !ok {
		ctrl.logger.Warn("User Name was not found in the header")

		// increment retrospectiveUpdateFail prometheus metric here
		ctrl.retrospectiveUpdateFail.Inc()

		handleNotInHeader(rw, r, "user")
		return
	}
	if _, ok := r.Header["Retrospective"]; !ok {
		ctrl.logger.Warn("Retrospective was not found in the header")

		// increment retrospectiveUpdateFail prometheus metric here
		ctrl.retrospectiveUpdateFail.Inc()

		handleNotInHeader(rw, r, "retrospective")
		return
	}

	if ctrl.currentUser == r.Header["Username"][0] && ctrl.status {
		retrospective = r.Header["Retrospective"][0]
		ctrl.logger.Info("Restrospective updated by user", zap.String("user", ctrl.currentUser))

		// increment retrospectiveUpdateSuccess prometheus metric here
		ctrl.retrospectiveUpdateSuccess.Inc()

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Retrospective Updated"))
	} else if !ctrl.status {
		ctrl.logger.Warn("User tried to access resource without locking the resource")

		// increment retrospectiveUpdateFail prometheus metric here
		ctrl.retrospectiveUpdateFail.Inc()

		rw.WriteHeader(http.StatusBadRequest)
	} else {
		// increment retrospectiveUpdateFail prometheus metric here
		ctrl.retrospectiveUpdateFail.Inc()

		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Access Denied as resource is in use"))
		ctrl.logger.Warn("User tried to access while resource was locked")

	}
}

func (ctrl *RetrospectiveController) ReleaseRetrospective(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Username"]; !ok {
		ctrl.logger.Warn("User Name was not found in the header")

		// increment retrospectiveReleaseFail prometheus metric here
		ctrl.retrospectiveReleaseFail.Inc()

		handleNotInHeader(rw, r, "user")
		return
	}

	// only the current user should be able to call this.
	if ctrl.currentUser == r.Header["Username"][0] && ctrl.status {
		ctrl.access.Unlock()
		ctrl.status = false

		// increment retrospectiveReleaseSuccess prometheus metric here
		ctrl.retrospectiveReleaseSuccess.Inc()

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Retrospective Released"))
	} else if !ctrl.status {
		ctrl.logger.Warn("User tried to release resource without locking the resource")

		// increment retrospectiveReleaseFail prometheus metric here
		ctrl.retrospectiveReleaseFail.Inc()

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Retrospective could not be released as not locked."))
	} else {
		ctrl.logger.Warn("Unauthorized User tried to release")

		// increment retrospectiveReleaseFail prometheus metric here
		ctrl.retrospectiveReleaseFail.Inc()

		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Access Denied as no auth"))
	}
}
