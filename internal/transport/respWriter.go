package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type RespWriter struct {
	Text string `json:"message"`
}

func NewRespWriter(w http.ResponseWriter, text string, statusCode int, logger *slog.Logger) {
	w.WriteHeader(statusCode)
	res, err := json.Marshal(RespWriter{
		Text: text,
	})
	if err != nil {
		logger.Error("json marshaling failed", slog.Any("error", err))
		return
	}
	if _, err := w.Write(res); err != nil {
		logger.Error("writing in respone failed", slog.Any("error", err))
		return
	}
	logger.Info("repsonse write", slog.String("message", text), slog.Int("status code", statusCode))
}
