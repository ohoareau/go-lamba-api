package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/ohoareau/gola/utils"
)

func JwtAuth(r *chi.Mux) {
	r.Use(
		jwtauth.Verifier(
			jwtauth.New(
				utils.GetEnvVar("JWT_AUTH_ALG", "HS256"),
				[]byte(utils.GetEnvVar("JWT_AUTH_SECRET", "MySEcret2021!")),
			 nil)))
	r.Use(jwtauth.Authenticator)
}
