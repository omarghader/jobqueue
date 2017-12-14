package shared

import (
	"zenithar.org/go/common/registrar"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

var (
	// Config contains flags passed to the API
	Config *Flags
	// Log is the API's logrus instance
	Log *logrus.Logger
	// Raven is sentry client
	Raven *raven.Client
	// Registrar is the bean registry
	Registrar registrar.Registrar
	// ApplicationID is the application identifier
	ApplicationID string
	// TenantID is the tenant identifier
	TenantID string

	// // Tokens is the token cache service
	// Tokens services.JWTokenService
)

const (
	// CustomerRepository is the name of the CustomerRepository bean
	InstagramRepository = "application"
)
