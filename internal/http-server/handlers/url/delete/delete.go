package delete

import (
	"errors"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"log/slog"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5/middleware"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

//go:generate go run github.com/vektra/mockery/v2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to decode request body"))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		err = urlDeleter.DeleteURL(req.Alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", req.Alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to delete URL"))

			return
		}

		log.Info("url deleted", "alias", req.Alias)

		render.JSON(w, r, resp.OK())
	}
}
