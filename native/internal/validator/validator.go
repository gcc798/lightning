package validator

import (
	"regexp"
	"sync"

	playground "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var (
	macRegex   = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	phoneRegex = regexp.MustCompile(`^1\d{10}$`)
	snRegex    = regexp.MustCompile(`^[A-Za-z0-9_-]{6,40}$`)
	once       sync.Once
	engine     = playground.New()
)

func Register() {
	once.Do(func() {
		engine.SetTagName("binding")
		_ = engine.RegisterValidation("mac", validateMAC)
		_ = engine.RegisterValidation("cnphone", validatePhone)
		_ = engine.RegisterValidation("sn", validateSN)
	})
}

// EchoValidator bridges the existing binding tags to Echo validation.
type EchoValidator struct{}

func (EchoValidator) Validate(value any) error {
	Register()
	return engine.Struct(value)
}

// EchoBinder preserves Echo's bind-and-validate behavior on Echo.
type EchoBinder struct {
	defaultBinder echo.DefaultBinder
}

func (b *EchoBinder) Bind(c *echo.Context, target any) error {
	if err := b.defaultBinder.Bind(c, target); err != nil {
		return err
	}
	return c.Validate(target)
}

func validateMAC(fl playground.FieldLevel) bool   { return macRegex.MatchString(fl.Field().String()) }
func validatePhone(fl playground.FieldLevel) bool { return phoneRegex.MatchString(fl.Field().String()) }
func validateSN(fl playground.FieldLevel) bool    { return snRegex.MatchString(fl.Field().String()) }
