package routes

//
// import (
// 	"net/http"
//
// 	"github.com/omarghader/jobqueue/api"
// 	// "github.com/omarghader/jobqueue/shared"
//
// 	"github.com/Sirupsen/logrus"
// 	"github.com/zenazn/goji/web"
// 	"zenithar.org/go/common/web/rest/jsonld"
// 	"zenithar.org/go/common/web/utils"
// )
//
// // UserController is the condition search service entrypoint
// type UserController struct {
// 	// service api.ApiInstagram
// }
//
// // NewUserControllerAutowired returns a new User controller autowired with the registrar
// func NewUserControllerAutowired() *UserController {
// 	// service := shared.Registrar.Lookup("botService")
// 	// return NewUserController(service.(api.ApiInstagram))
// }
//
// // NewUserController returns a new condition controller
// func NewUserController(service api.ApiInstagram) *UserController {
// 	return &UserController{
// 		service: service,
// 	}
// }
//
// //------------------------------------------------------------------------------
// // Resources
// //------------------------------------------------------------------------------
//
// //------------------------------------------------------------------------------
// // Handlers
// //------------------------------------------------------------------------------
//
// //------------------------------------------------------------------------------
// // GetUser returns the User address condition
// func (ctrl *UserController) Get(c web.C, w http.ResponseWriter, r *http.Request) {
//
// 	// Do the query
// 	_, err := ctrl.service.GetUser(ctrl.service.GetSelf().Username, true)
// 	if err != nil {
// 		logrus.WithError(err).Error("Unable to retrieve User condition from database !")
//
// 		status := jsonld.NewError("JSONLDContext", utils.GetHTTPUrl(r), "User/GET/01", err.Error())
// 		status.Write(w)
// 		return
// 	}
//
// 	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{})
// }
