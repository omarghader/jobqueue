package system

import (
	"strings"

	"github.com/omarghader/jobqueue/dao/mongodb"
	"github.com/omarghader/jobqueue/services/publisher"
	"github.com/omarghader/jobqueue/services/scheduler"
	"github.com/omarghader/jobqueue/shared"

	"github.com/Sirupsen/logrus"
	"zenithar.org/go/common/registrar"
)

// Setup the application
func Setup(flags *shared.Flags) Application {
	// Initialize the registrar
	shared.Registrar = registrar.Registry()

	logrus.Infoln("**********************************************************")

	// Create a MySQL-based persistent backend.
	store, err := mongodb.NewStore("mongodb://localhost:27017/joqueue")
	if err != nil {
		panic(err)
	}

	// Create a scheduler Service with 10 concurrent workers.
	schedulerService := scheduler.NewService()
	schedulerService.SetStore(store)
	schedulerService.SetConcurrency(1, 10)
	shared.Registrar.Register("schedulerService", schedulerService)
	logrus.Debugln("[SRV] scheduler service")

	// Create a manager with the MySQL store and 10 concurrent workers.
	publisherService, err := publisher.NewService(flags.AmqpURI, strings.Split(flags.WebhookURLs, ","))
	if err != nil {
		logrus.Errorln("Cannot Register Publisher Service", err)
	}
	shared.Registrar.Register("publisherService", publisherService)
	logrus.Debugln("[SRV] Publisher service")

	logrus.Infoln("**********************************************************")

	return &baseApplication{
		Scheduler: schedulerService,
		Publisher: publisherService,
	}
}
