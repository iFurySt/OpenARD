package tracecontext

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

const Header = "traceparent"

type contextKey struct{}

type TraceContext struct {
	TraceID string
	SpanID  string
	Flags   string
}

func (trace TraceContext) String() string {
	if trace.Flags == "" {
		trace.Flags = "00"
	}
	return "00-" + trace.TraceID + "-" + trace.SpanID + "-" + trace.Flags
}

func Start(ctx context.Context, incoming string) (context.Context, TraceContext) {
	if trace, ok := From(ctx); ok {
		return ctx, trace
	}
	if parent, ok := Parse(incoming); ok {
		trace := TraceContext{
			TraceID: parent.TraceID,
			SpanID:  randomHex(8),
			Flags:   parent.Flags,
		}
		return With(ctx, trace), trace
	}
	return Ensure(ctx)
}

func Ensure(ctx context.Context) (context.Context, TraceContext) {
	if trace, ok := From(ctx); ok {
		return ctx, trace
	}
	trace := TraceContext{
		TraceID: randomHex(16),
		SpanID:  randomHex(8),
		Flags:   "00",
	}
	return With(ctx, trace), trace
}

func With(ctx context.Context, trace TraceContext) context.Context {
	trace.TraceID = strings.ToLower(strings.TrimSpace(trace.TraceID))
	trace.SpanID = strings.ToLower(strings.TrimSpace(trace.SpanID))
	trace.Flags = strings.ToLower(strings.TrimSpace(trace.Flags))
	if trace.Flags == "" {
		trace.Flags = "00"
	}
	if !validTrace(trace) {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, trace)
}

func From(ctx context.Context) (TraceContext, bool) {
	trace, ok := ctx.Value(contextKey{}).(TraceContext)
	if !ok || !validTrace(trace) {
		return TraceContext{}, false
	}
	return trace, true
}

func Parse(value string) (TraceContext, bool) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(value)), "-")
	if len(parts) != 4 {
		return TraceContext{}, false
	}
	if parts[0] != "00" {
		return TraceContext{}, false
	}
	trace := TraceContext{
		TraceID: parts[1],
		SpanID:  parts[2],
		Flags:   parts[3],
	}
	if !validTrace(trace) {
		return TraceContext{}, false
	}
	return trace, true
}

func SetHeader(header http.Header, ctx context.Context) {
	if trace, ok := From(ctx); ok {
		header.Set(Header, trace.String())
	}
}

func validTrace(trace TraceContext) bool {
	return validHex(trace.TraceID, 32, true) &&
		validHex(trace.SpanID, 16, true) &&
		validHex(trace.Flags, 2, false)
}

func validHex(value string, length int, rejectAllZero bool) bool {
	if len(value) != length {
		return false
	}
	if _, err := hex.DecodeString(value); err != nil {
		return false
	}
	if rejectAllZero && strings.Trim(value, "0") == "" {
		return false
	}
	return true
}

func randomHex(bytesLength int) string {
	for {
		buffer := make([]byte, bytesLength)
		if _, err := rand.Read(buffer); err != nil {
			panic(err)
		}
		value := hex.EncodeToString(buffer)
		if strings.Trim(value, "0") != "" {
			return value
		}
	}
}
