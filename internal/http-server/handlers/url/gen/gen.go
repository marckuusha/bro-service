package gen

import (
	"context"
	"net/http"
	"url-shortener/internal/client/grpcserv"
	"url-shortener/internal/model"

	"github.com/go-chi/render"
)

func New(grpcClient *grpcserv.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()
		name := queryParams.Get("name")

		if name == "" {
			render.JSON(w, r, model.Response{
				Status: http.StatusBadRequest,
				Error:  "name is empty",
			})
			return
		}

		str, err := grpcClient.SendMsg(context.Background(), name)
		if err != nil {
			render.JSON(w, r, model.Response{
				Status: http.StatusBadRequest,
				Error:  "failed send to grpc req",
			})
			return
		}

		render.JSON(w, r, model.Response{
			Status:   http.StatusOK,
			Template: str,
		})
	}
}
