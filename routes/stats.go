package routes

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/omarghader/jobqueue/protocol"
	"github.com/omarghader/jobqueue/services"
	"github.com/zenazn/goji/web"
	"zenithar.org/go/common/web/rest/jsonld"
	"zenithar.org/go/common/web/utils"
)

// StatsController is the condition search service entrypoint
type StatsController struct {
	schedulerService services.SchedulerService
}

// NewStatsControllerAutowired returns a new Stats controller autowired with the registrar
// func NewStatsControllerAutowired() *models.Stat {
// service := shared.Registrar.Lookup("botService")
// return NewStatsController(service.(api.ApiInstagram))
// }

// NewStatsController returns a new condition controller
func NewStatsController(schedulerService services.SchedulerService) *StatsController {
	return &StatsController{
		schedulerService: schedulerService,
	}
}

//------------------------------------------------------------------------------
// Resources
//------------------------------------------------------------------------------
type ListResource struct {
	Jobs  []protocol.Job
	Total int
}

//------------------------------------------------------------------------------
// Handlers
//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
// GetStats returns the Stats address condition
func (ctrl *StatsController) Get(c web.C, w http.ResponseWriter, r *http.Request) {

	// Do the query
	res, err := ctrl.schedulerService.Stats(&protocol.StatsRequest{})
	if err != nil {
		logrus.WithError(err).Error("Unable to retrieve Stats condition from database !")

		status := jsonld.NewError("JSONLDContext", utils.GetHTTPUrl(r), "Stats/GET/01", err.Error())
		status.Write(w)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"Failed":    res.Failed,
		"Succeeded": res.Succeeded,
		"Waiting":   res.Waiting,
		"Working":   res.Working,
	})
}

func (ctrl *StatsController) List(c web.C, w http.ResponseWriter, r *http.Request) {

	input := &protocol.ListRequest{
		CorrelationID:    r.URL.Query().Get("correlation_id"),
		CorrelationGroup: r.URL.Query().Get("correlation_group"),
		// Limit:            r.URL.Query().Get("limit"),
		// Offset:           r.URL.Query().Get("offset"),
		State: r.URL.Query().Get("state"),
		Topic: r.URL.Query().Get("topic"),
	}
	// Do the query
	res, err := ctrl.schedulerService.List(input)
	if err != nil {
		logrus.WithError(err).Error("Unable to retrieve Stats condition from database !")

		status := jsonld.NewError("JSONLDContext", utils.GetHTTPUrl(r), "Stats/GET/01", err.Error())
		status.Write(w)
		return
	}
	utils.JSONResponse(w, http.StatusOK, ListResource{
		Jobs:  res.Jobs,
		Total: res.Total,
	})
}

func (ctrl *StatsController) Add(c web.C, w http.ResponseWriter, r *http.Request) {

	job := &protocol.Job{}
	err := utils.ParseRequest(r, job)
	if err != nil {
		fmt.Println("Error Parsing data", err)
	}
	// Do the query
	err = ctrl.schedulerService.Add(job)
	if err != nil {
		logrus.WithError(err).Error("Unable to retrieve Stats condition from database !")

		status := jsonld.NewError("JSONLDContext", utils.GetHTTPUrl(r), "Stats/GET/01", err.Error())
		status.Write(w)
		return
	}
	utils.JSONResponse(w, http.StatusOK, ListResource{
	// Jobs:  res.Jobs,
	// Total: res.Total,
	})
}
