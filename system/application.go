package system

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/omarghader/jobqueue/services"
	"github.com/omarghader/jobqueue/shared"

	"github.com/omarghader/jobqueue/protocol"
	"github.com/omarghader/jobqueue/routes"

	//"github.com/sebest/xff"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	instabotProtocol "omarghader.com/instagram/instabot/protocol"
	mdlwr "zenithar.org/go/common/web/middleware"
)

// Application is the application instance
type Application interface {
	Router() http.Handler
	Start() error
	Stop() error
}

type baseApplication struct {
	Scheduler services.SchedulerService
	Publisher services.PublisherService
}

// -----------------------------------------------------------------------------

func (a *baseApplication) Router() http.Handler {
	// Create a new goji mux
	mux := web.New()

	// Initialize middlewares
	// mux.Use(xff.Handler)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	//mux.Use(mdlwr.XRequestID)
	mux.Use(mdlwr.CORS)
	mux.Use(middleware.AutomaticOptions)
	mux.Use(mdlwr.GzipHandler)
	if shared.Raven != nil {
		mux.Use(mdlwr.BuildErrorCatcher(shared.Raven))
	} else {
		mux.Use(middleware.Recoverer)
	}

	// Defines the routes
	mux.Get("/", routes.Hello)

	schedulerService := shared.Registrar.Lookup("schedulerService")
	statsController := routes.NewStatsController(schedulerService.(services.SchedulerService))
	mux.Get("/api/v1/stats", statsController.Get)
	mux.Get("/api/v1/jobs", statsController.List)
	mux.Post("/api/v1/job", statsController.Add)

	// Compile the routes
	mux.Compile()

	return mux
}

// Start asynchronous tasks
func (a *baseApplication) Start() error {
	m := a.Scheduler
	// Register one or more topics and their processor
	m.Register("clicks", func(args ...interface{}) error {
		// Handle "clicks" topic
		fmt.Println("[CLICKS] Received Event at : ", time.Now())
		return nil
	})

	m.Register("publish:rabbitmq", func(args ...interface{}) error {
		// Handle "clicks" topic
		notification := args[0].(protocol.Notification)
		from := ""
		to := ""
		metadata := map[string]string{}
		if len(args) > 1 {
			for i := 0; i < len(args[1].([]interface{})); i++ {
				v := args[1].([]interface{})[i].(string)
				if strings.Contains(v, "from:") {
					from = strings.Split(v, "from:")[1]
				} else if strings.Contains(v, "to:") {
					to = strings.Split(v, "to:")[1]
				} else {
					fields := strings.Split(v, ":")
					metadata[fields[0]] = fields[1]
				}
			}
		}

		typ := instabotProtocol.EventTypeMap[strings.ToUpper(notification.AmqpExchangeRoutingKey)]

		err := a.Publisher.PublishRabbitmq(notification.AmqpExchangeName, notification.AmqpExchangeRoutingKey, &instabotProtocol.Event{
			Type: typ,
			Observable: &instabotProtocol.Observable{
				Value: from,
				Relations: []*instabotProtocol.ObservableRelation{
					&instabotProtocol.ObservableRelation{
						Destination: to,
					},
				},
			},
			Metadata: metadata,
		})
		if err != nil {
			logrus.Errorln(err)
			return err
		}
		fmt.Println("[publish:rabbitmq] Received Event at : ", time.Now())
		return nil

	})

	m.Register("publish:webhook", func(args ...interface{}) error {
		// Handle "clicks" topic
		err := a.Publisher.PublishWebhook(&protocol.Event{})
		if err != nil {
			logrus.Errorln(err)
			return err
		}
		fmt.Println("[publish:webhook] Received Event at : ", time.Now())
		return nil
	})

	// Start the manager
	err := m.Start()
	if err != nil {
		panic(err)
	}
	//
	// t, err := time.Parse(time.RFC3339, "2017-07-14T09:22:00+02:00")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// t := time.Now().Add(time.Second * 5)
	// fmt.Println(time.Now().Sub(t))
	// // Add a job: It'll be added to the store and processed eventually.
	// err = m.Add(&protocol.Job{Topic: "publish:rabbitmq", Args: []interface{}{640, 480},
	// 	CorrelationID: "click_func",
	// 	// FirstStart:    t.UnixNano(),
	// 	Delay:  3,
	// 	Repeat: 5,
	// })
	// if err != nil {
	// 	logrus.Errorln(err)
	// }

	return nil
}

// Stop asynchronous tasks
func (a *baseApplication) Stop() error {
	// Stop the manager, either via Stop/Close (which stops after all workers
	// are finished) or CloseWithTimeout (which gracefully waits for a specified
	// time span)
	err := a.Scheduler.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
