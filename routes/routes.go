package routes

import (
	"app/auth"
	"app/db"
	lg "app/logging"
	"app/structs"
	"app/utils"
	"app/validations"
	"app/weberrors"
	"net/http"
	"os"

	_ "app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitApp(app *gin.Engine) {
	app.Use(gin.Recovery())
	app.Use(lg.Middleware())
	app.Use(weberrors.JSONAppErrorReporter())

	app.GET("/healthcheck", HealthCheckHandler)
	app.POST("/event", CreateEventHandler)
	app.GET("/event/:id", GetEventHandler)

	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adminGroup := app.Group("/")
	adminGroup.Use(auth.Middleware())
	adminGroup.DELETE("/event/:id", DeleteEventHandler)

	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NotFound", "description": ""})
	})
	validations.BindCustomValidators()
}

// CreateEventHandler creates event.
// @Summary	Creates event to database
// @Description NOTE: videoQuality & audioQuality default to 720p & Low, respectively if not provided. If provided default value is first in the list.
// @Tags		event
// @Accept json
// @Produce json
// @Param event body structs.EventData true "Event Data"
// @Success	201 {object} structs.EventData
// @Failure 400,404,500 {object} weberrors.AppError
// @Router		/event [post]
var CreateEventHandler = func(ctx *gin.Context) {
	eventData := structs.EventData{}
	bindError := ctx.ShouldBind(&eventData)
	if bindError != nil {
		if parsedErr := validations.GetBindErrors(bindError); parsedErr != nil {
			utils.AppendContextError(ctx, parsedErr)
			return
		}
		utils.AppendContextError(ctx, &weberrors.InvalidPayload)
		return
	}
	// setting default values if not provided in payload
	if len(eventData.VideoQuality) == 0 {
		eventData.VideoQuality = []string{utils.DEFAULT_RESOLUTION}
	}
	if len(eventData.AudioQuality) == 0 {
		eventData.AudioQuality = []string{utils.DEVAULT_AUDIO}
	}
	id, err := db.CreateEvent(eventData)
	if err != nil {
		utils.AppendContextError(ctx, &weberrors.InternalError)
		return
	}
	eventData.Id = id
	ctx.JSON(http.StatusCreated, eventData)
}

// GetEventHandler retrieves event.
// @Summary	Retrieves event from database
// @Tags		event
// @Param id path string true "Event ID (uuid)"
// @Produce json
// @Success	200
// @Failure 404,500 {object} weberrors.AppError
// @Router		/event/{id} [get]
var GetEventHandler = func(ctx *gin.Context) {
	id := ctx.Param("id")
	if !validations.CheckUuidFormat(id) {
		utils.AppendContextError(ctx, &weberrors.NotFound)
		return
	}
	response, err := db.GetEvent(id)
	if err != nil {
		if err.Error() == "not found" {
			utils.AppendContextError(ctx, &weberrors.NotFound)
			return
		}
		utils.AppendContextError(ctx, &weberrors.InternalError)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// DeleteEventHandler removes event.
// @Summary	Delete event from database
// @Tags		event
// @Param id path string true "Event ID (uuid)"
// @Success	204
// @Failure 500 {object} weberrors.AppError
// @Router		/event/{id} [delete]
var DeleteEventHandler = func(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
	id := ctx.Param("id")
	if !validations.CheckUuidFormat(id) {
		return
	}
	err := db.DeleteEvent(id)
	if err != nil {
		if err.Error() != "not found" {
			utils.AppendContextError(ctx, &weberrors.InternalError)
		}
	}
}

// HealthCheckHandler checks the status of the server.
// @Summary	Check health of this service
// @Tags		Health check
// @Produce	json
// @Success	200	{object}	structs.JsonHealthCheckStatus
// @Router		/healthcheck [get]
var HealthCheckHandler = func(ctx *gin.Context) {
	var status structs.JsonHealthCheckStatus
	status.Result = "ok"
	status.Version = os.Getenv("COMMIT_TAG")     // would be set in pipeline
	status.DeployDate = os.Getenv("DEPLOY_DATE") // would be set in pipeline
	ctx.JSON(http.StatusOK, status)
}
