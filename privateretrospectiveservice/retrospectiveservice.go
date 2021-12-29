package couponservice

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"
)

// RetrospectiveController is the Real Time Socket route handler
type RetrospectiveController struct {
	logger *zap.Logger
}

// NewRetrospectiveController returns a fresh Retrospective controller
func NewRetrospectiveController(logger *zap.Logger, instance *redis.Client) *RetrospectiveController {
	return &RetrospectiveController{
		logger: logger,
	}
}

func (ctrl *RetrospectiveController) AssignSocket(rw http.ResponseWriter, r *http.Request) {
	m := melody.New()
	fmt.Print(m)
}
