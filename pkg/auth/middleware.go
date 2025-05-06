package auth

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type MiddleWare struct {
	tm TokenManager
}

func NewMiddleWare(tm TokenManager) MiddleWare {
	return MiddleWare{tm: tm}
}

func (w *MiddleWare) IsAuthed(handle httprouter.Handle) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		var id int
		header := request.Header.Get("Authorization")
		headerArray := strings.Split(header, " ")

		id, role, err := w.tm.ValidateToken(headerArray[1])
		if err != nil {
			return
		}
		ctx := context.WithValue(request.Context(), "user_id", id)
		ctx = context.WithValue(request.Context(), "user_role", role)

		handle(writer, request.WithContext(ctx), params)
	}
}
