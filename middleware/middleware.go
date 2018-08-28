package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// HandleWithError is an httprouter.Handle that returns an error.
type HandleWithError func(http.ResponseWriter, *http.Request, httprouter.Params) error

// HTTP runs HandleWithError and converts it to httprouter.Handle.
// The conversion is needed because httprouter.Router needs httprouter.Handle
// in its signature.
func HTTP(handle HandleWithError) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		reqID := r.Header.Get("X-Request-ID")
		start := time.Now()

		err := handle(w, r, params)

		// elapsed time in milliseconds
		elapsed := time.Since(start).Seconds() * 1000
		elapsedStr := strconv.FormatFloat(elapsed, 'f', -1, 64)

		logger, _ := zap.NewProduction()

		if err != nil {
			logger.Error(err.Error(),
				zap.String("request_id", reqID),
				zap.String("duration", elapsedStr),
				zap.Strings("tags", []string{r.URL.Path}),
			)
		} else {
			logger.Info("everything is fine",
				zap.String("request_id", reqID),
				zap.String("duration", elapsedStr),
				zap.Strings("tags", []string{r.URL.Path}),
			)
		}
		if err != nil && err.Error() == "Close Explisitly" {
			os.Exit(1)
		}
	}
}
