package logger

import (
	"context"
	"io"
	stdLog "log"
	"log/slog"
)

type NullHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type NullHandler struct {
	slog.Handler
	l *stdLog.Logger
}

func (opts NullHandlerOptions) NewNullHandler(
	out io.Writer,
) *NullHandler {
	h := &NullHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

func (h *NullHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

func (h *NullHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *NullHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
