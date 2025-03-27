package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"main/internal/model"
	"net/http"
	"strings"
)

type Middleware struct {
	cfg *model.Config
}

func NewMiddleware(cfg *model.Config) *Middleware {
	return &Middleware{cfg}
}

func (m *Middleware) BodyLogger(c echo.Context, requestBody []byte, responseBody []byte) {

	method := c.Request().Method
	statusCode := c.Response().Status
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	requestedUrl := c.Request().URL.String()

	logMsg := fmt.Sprintf("| %s %s | STATUSCODE %d | METHOD %s | URL %s | RESPONSE %s | REQUEST %s |",
		echo.HeaderXRequestID, requestId,
		statusCode,
		method,
		requestedUrl,
		responseBody,
		requestBody)

	log.Println(strings.Join(strings.Fields(logMsg), " "))
}

func (m *Middleware) Authorization(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		token := c.Get("Authorization")
		if token == nil || token.(string) == "" {
			return c.String(http.StatusUnauthorized, "Authorization token is empty")
		}

		return next(c)
	}
}
