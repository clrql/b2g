package routes

import (
	"b2g/middlewares"
	"b2g/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func MainRouter(r fiber.Router) {
	r.Post("/auth/register", auth.Register)
	r.Get("/auth/verify", auth.Verify)
	r.Post("/auth/login", auth.Login)
	r.Get("/auth/authed", middlewares.Authed, auth.Authed)
	r.Get("/auth/logout", auth.Logout)
	r.Get("/auth/get-reset-password-token", auth.GetResetPasswordToken)
	r.Post("/auth/reset-password", auth.ResetPassword)
}
