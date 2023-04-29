package logging

import (
	"app/utils"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	timestampFieldName      = "timestamp"
	levelFieldName          = "levelname"
	requestStartedMessage   = "Request started"
	requestEndedMessage     = "Request ended"
	correlationIdContextKey = "correlationId"
	logKey                  = "log"
)

const contentType = "Content-Type"

func init() {
	zerolog.LevelFieldName = levelFieldName
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = timestampFieldName
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().
		Str("environment", os.Getenv("ENV")).
		Str("service", utils.APP_NAME).
		Str("version", utils.GetEnvOrDefault("COMMIT_TAG", "main")).
		Timestamp().Logger()

	log.Logger.Info().Msg("Finished logger setup")
}

func WithContext(c context.Context) *zerolog.Logger {
	if ctxLogger, found := c.Value(logKey).(zerolog.Logger); found {
		return &ctxLogger
	}
	return &log.Logger
}

func Middleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		startTime := time.Now().UTC()
		correlationId := uuid.NewString()
		gctx.Set(correlationIdContextKey, correlationId)
		innitialLog := log.Logger.With().
			Str("correlation_id", correlationId).
			Dict("http", zerolog.Dict().
				Str("method", gctx.Request.Method).
				Str("url", gctx.Request.URL.String()).
				Str("host", gctx.Request.Host).
				Str("query", gctx.Request.URL.RawQuery).
				Str("user_agent", gctx.Request.UserAgent()).
				Str("remote_address", gctx.Request.RemoteAddr)).
			Str("rq_started", startTime.Format(time.RFC3339)).
			Str("rq_content_type", gctx.Request.Header.Get("Content-Type")).
			Timestamp().Logger()

		innitialLog.Info().Msg(requestStartedMessage)

		gctx.Next()

		exitLog := log.Logger.With().
			Str("correlation_id", correlationId).
			Dur("rs_time_ms", time.Since(startTime)).
			Timestamp().Logger()

		switch {
		case gctx.Writer.Status() >= http.StatusBadRequest &&
			gctx.Writer.Status() < http.StatusInternalServerError:
			{
				exitLog.WithLevel(zerolog.WarnLevel).
					Msg(requestEndedMessage)
			}
		case gctx.Writer.Status() >= http.StatusInternalServerError:
			{
				exitLog.WithLevel(zerolog.ErrorLevel).
					Msg(requestEndedMessage)
			}
		default:
			exitLog.WithLevel(zerolog.InfoLevel).
				Msg(requestEndedMessage)
		}
	}
}
