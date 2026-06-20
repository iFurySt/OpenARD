package httpapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ifuryst/ard/internal/requestid"
	"github.com/ifuryst/ard/internal/tracecontext"
)

const requestIDKey = "request_id"
const traceContextKey = "trace_context"

func requestIDMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestID := strings.TrimSpace(context.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.NewString()
		}
		context.Set(requestIDKey, requestID)
		context.Header("X-Request-ID", requestID)
		context.Request = context.Request.WithContext(requestid.With(context.Request.Context(), requestID))
		context.Next()
	}
}

func traceContextMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestContext, trace := tracecontext.Start(context.Request.Context(), context.GetHeader(tracecontext.Header))
		context.Set(traceContextKey, trace)
		context.Header(tracecontext.Header, trace.String())
		context.Request = context.Request.WithContext(requestContext)
		context.Next()
	}
}

func jsonAccessLogMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		startedAt := time.Now()
		context.Next()
		trace := traceContextFromContext(context)
		event := map[string]any{
			"ts":        time.Now().UTC().Format(time.RFC3339Nano),
			"level":     "info",
			"event":     "http_request",
			"requestId": requestIDFromContext(context),
			"traceId":   trace.TraceID,
			"spanId":    trace.SpanID,
			"method":    context.Request.Method,
			"path":      context.Request.URL.Path,
			"status":    context.Writer.Status(),
			"latencyMs": time.Since(startedAt).Milliseconds(),
			"clientIp":  context.ClientIP(),
		}
		if len(context.Errors) > 0 {
			event["level"] = "error"
			event["errors"] = context.Errors.String()
		}
		data, err := json.Marshal(event)
		if err != nil {
			return
		}
		fmt.Fprintln(gin.DefaultWriter, string(data))
	}
}

func requestIDFromContext(context *gin.Context) string {
	value, ok := context.Get(requestIDKey)
	if !ok {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}

func traceContextFromContext(context *gin.Context) tracecontext.TraceContext {
	value, ok := context.Get(traceContextKey)
	if !ok {
		return tracecontext.TraceContext{}
	}
	trace, _ := value.(tracecontext.TraceContext)
	return trace
}
