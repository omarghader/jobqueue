package main

import (
	"net"
	"net/http"
	"time"

	"github.com/omarghader/jobqueue/shared"
	"github.com/omarghader/jobqueue/system"

	"github.com/Sirupsen/logrus"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zenazn/goji/graceful"
)

var (
	// General flags
	// definitionFilePath = flag.String("conf", "config/config.yml", "config definition file to load")

	bindAddress      = flag.String("bind", ":5000", "Network address used to bind")
	logFormatterType = flag.String("log", "text", "Log formatter type. Either \"json\" or \"text\"")
	logLevel         = flag.String("log_level", "info", "Defines the log level (panic, fatal, error, warn, info, debug)")
	forceColors      = flag.Bool("force_colors", false, "Force colored prompt?")
	ravenDSN         = flag.String("raven_dsn", "", "Defines the sentry endpoint dsn")

	databaseDriver    = flag.String("db_driver", "mongodb", "Specify the database to use (mongodb, rethinkdb)")
	databaseHost      = flag.String("db_host", "localhost:27017", "Database hosts, split by ',' to add several hosts")
	databaseNamespace = flag.String("db_namespace", "markos", "Select the database")
	databaseUser      = flag.String("db_user", "", "Database user")
	databasePassword  = flag.String("db_password", "", "Database user password")

	amqpURI     = flag.String("amqp_uri", "amqp://guest:guest@localhost:5672", "Defines where the RabbitMQ server is")
	webhookURLS = flag.String("webhook_urls", "http://localhost:8080", "Defines where the webhook urls is")

	memcachedHosts = flag.String("memcached_hosts", "", "Memcached servers for cache (ex: 127.0.0.1:11211)")
	redisHost      = flag.String("redis_host", "", "Redis server for cache")
)

func init() {
	flag.Parse()

	// Set localtime to UTC
	time.Local = time.UTC

	logrus.Infoln("**********************************************************")
	logrus.Infoln("Scheduler starting ...")
	// logrus.Infof("Version : %s (%s-%s)", version.Version, version.Revision, version.Branch)

	// Set the formatter depending on the passed flag's value
	if *logFormatterType == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: *forceColors,
		})
	} else if *logFormatterType == "json" {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: *forceColors,
		})
	}

	// Defines the log level
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatalln("Invalid log level ! (panic, fatal, error, warn, info, debug) ")
	}
	logrus.SetLevel(level)

}

func main() {

	// Put config into the environment package
	shared.Config = &shared.Flags{
		BindAddress:      *bindAddress,
		LogFormatterType: *logFormatterType,
		ForceColors:      *forceColors,
		RavenDSN:         *ravenDSN,

		DatabaseDriver:    *databaseDriver,
		DatabaseHost:      *databaseHost,
		DatabaseNamespace: *databaseNamespace,
		DatabaseUser:      *databaseUser,
		DatabasePassword:  *databasePassword,

		AmqpURI:     *amqpURI,
		WebhookURLs: *webhookURLS,
		// DefinitionFilePath: *definitionFilePath,
		MemcachedHosts: *memcachedHosts,
		RedisHost:      *redisHost,
	}

	// Initialize the application
	app := system.Setup(shared.Config)

	// Start application
	go app.Start()

	logrus.Infoln("**********************************************************")

	// Make the mux handle every request

	logrus.Infoln("[PROM] Metrics endpoint : '/metrics'")
	http.Handle("/metrics", prometheus.Handler())
	http.Handle("/", app.Router())

	// Log that we're starting the server
	logrus.WithFields(logrus.Fields{
		"address": shared.Config.BindAddress,
	}).Info("Starting the HTTP server")

	// Initialize the goroutine listening to signals passed to the app
	graceful.HandleSignals()

	// Pre-graceful shutdown event
	graceful.PreHook(func() {
		logrus.Info("Received a signal, stopping the application")
		app.Stop()
	})

	// Post-shutdown event
	graceful.PostHook(func() {
		logrus.Info("Stopped the application")
	})

	// Listen to the passed address
	listener, err := net.Listen("tcp", shared.Config.BindAddress)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"address": *bindAddress,
		}).Fatal("Cannot set up a TCP listener")
	}

	// Start the listening
	err = graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		// Don't use .Fatal! We need the code to shut down properly.
		logrus.Error(err)
	}
	// If code reaches this place, it means that it was forcefully closed.

	// Wait until open connections close.
	graceful.Wait()
}
