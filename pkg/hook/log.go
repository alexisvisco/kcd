package hook

import (
	"fmt"
	"net/http"

	"github.com/expectedsh/errors"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func Log(_ http.ResponseWriter, r *http.Request, err error) {
	var logger *logrus.Entry

	e, ok := err.(*errors.Error)
	if ok {
		if logrus.GetLevel() <= logrus.DebugLevel {
			fmt.Println()
			fmt.Println(e.FormatStacktrace())
		}
		logger = e.Log()
	} else {
		logger = logrus.WithError(err)
	}

	reqID := middleware.GetReqID(r.Context())
	if reqID != "" {
		logger = logger.WithField("request-id", reqID)
	}

	logger.WithFields(map[string]interface{}{
		"remote":  r.RemoteAddr,
		"request": r.URL.Path,
		"params":  r.URL.RawQuery,
		"method":  r.Method,
	}).Error()
}
