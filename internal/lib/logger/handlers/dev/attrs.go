package dev

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v2"
)

const (
	json_marshaller uint8 = iota
	yaml_marshaller
)

type nextFunc func([]string, slog.Attr) slog.Attr

type attrsMap map[string]any

func WithYamlMarshaller() optsFunc {
	return func(h *Handler) {
		h.marshalType = yaml_marshaller
	}
}

func WithJsonMarshaller() optsFunc {
	return func(h *Handler) {
		h.marshalType = json_marshaller
	}
}

func suppressDefaultAttrs(
	next nextFunc,
) nextFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey, slog.LevelKey, slog.MessageKey:
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func (h *Handler) computeAttrs(
	ctx context.Context,
	rec slog.Record,
) (attrsMap, error) {

	h.mutex.Lock()
	defer func() {
		h.buf.Reset()
		h.mutex.Unlock()
	}()

	if err := h.handler.Handle(ctx, rec); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs attrsMap
	err := json.Unmarshal(h.buf.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func (h *Handler) marshal(attrs attrsMap) (string, error) {
	var (
		data []byte
		err  error
	)
	if len(attrs) > 0 {

		switch h.marshalType {
		case json_marshaller:
			data, err = json.MarshalIndent(attrs, "", "  ")

		case yaml_marshaller:
			data, err = yaml.Marshal(attrs)
		}

		if err != nil {
			return "", err
		}
	}
	return string(data), nil
}
