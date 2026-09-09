package httpx

import "github.com/labstack/echo/v5"

// Handler adapts controllers that write responses directly to Echo's handler contract.
func Handler(handler func(*echo.Context)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		handler(c)
		return nil
	}
}
