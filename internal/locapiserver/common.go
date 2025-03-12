package locapiserver

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

/* Common */
type HttpErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *HttpErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (s *LocApiServer) httpErrUnauthorized(err error) render.Renderer {
	return &HttpErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorText:      "Unauthorized",
	}
}

func (s *LocApiServer) httpErrUnexpected(err error) render.Renderer {
	return &HttpErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorText:      "Internal Server Error",
	}
}

func (s *LocApiServer) httpErrInvalidRequest(err error) render.Renderer {
	return &HttpErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      "Invalid Request",
	}
}

func getCtxValueString(ctx context.Context, key string) string {
	ret := ctx.Value(key)
	if ret == nil {
		return ""
	}

	return ret.(string)
}

