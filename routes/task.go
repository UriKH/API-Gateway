package routes

import (
	"net/http"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/schemas"
	tasks "github.com/TekClinic/Tasks-MicroService/tasks_protobuf"
	"github.com/gin-gonic/gin"
)

const resourceNameTask = "task"

type TasksParams struct {
	Skip   int32  `form:"skip,default=0"`
	Limit  int32  `form:"limit,default=20"`
	Search string `form:"search" binding:"omitempty,min=1,max=100"`
}

func getTasks(service tasks.TasksServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the query
		var params TasksParams
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call task microservice
		response, err := service.GetTasksIDs(ctx, &tasks.GetTasksIDsRequest{
			Token:  ctx.GetString(middlewares.TokenKey),
			Limit:  params.Limit,
			Offset: params.Skip,
			Search: params.Search,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			CreateNamedAPIResourceList(ctx, resourceNameTask,
				params.Skip, params.Limit, response.GetCount(), response.GetResults()))
	}
}

type TaskParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getTask(service tasks.TasksServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var uriParams TaskParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call task microservice
		response, err := service.GetTask(ctx, &tasks.GetTaskRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		if response.GetTask() == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Message: "Invalid response from the server.",
			})
			return
		}

		task := response.GetTask()

		ctx.JSON(http.StatusOK,
			schemas.Task{
                TaskBase: schemas.TaskBase{
                    PatientId: task.GetPatientId(),
                    Expertise: task.GetExpertise(),
                    Title: task.GetTitle(),
                    Description: task.GetDescription(),
                },
                Id: task.GetId(),
                CreatedAt: task.GetCreatedAt(),
                Complete: task.GetComplete(),
			})
	}
}

func createTask(service tasks.TasksServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the body
		var bodyParams schemas.TaskBase
		err := ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call task microservice
		response, err := service.CreateTask(ctx, &tasks.CreateTaskRequest{
            Token: ctx.GetString(middlewares.TokenKey),
            Title: bodyParams.Title,
            Description: bodyParams.Description,
            Expertise: bodyParams.Expertise,
            PatientId: bodyParams.PatientId,
        })
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusCreated, schemas.IDHolder{
			ID: response.GetId(),
		})
	}
}

func deleteTask(service tasks.TasksServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var uriParams TaskParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call task microservice
		_, err = service.DeleteTask(ctx, &tasks.DeleteTaskRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

type UpdateTaskParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func updateTask(service tasks.TasksServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams UpdateTaskParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var bodyParams schemas.TaskUpdate
		err = ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call task microservice
		response, err := service.UpdateTask(ctx, &tasks.UpdateTaskRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Task: &tasks.Task{
				Id:     uriParams.ID,
                // TODO: We need to get this from the request...
				// Complete: bodyParams.Complete,
                Complete: false,
                Title:  bodyParams.Title,
                Description: bodyParams.Description,
                Expertise: bodyParams.Expertise,
                PatientId: bodyParams.PatientId,
                // TODO: remove this!
                CreatedAt: "2020-02-20",
			},
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, schemas.IDHolder{
			ID: response.GetId(),
		})
	}
}

func RegisterTaskRoutes(router *gin.Engine) {
	client := InitiateClient(resourceNameTask, tasks.NewTasksServiceClient)

	router.GET("/tasks", getTasks(client))
	router.POST("/tasks", createTask(client))
	router.GET("/tasks/:id", getTask(client))
	router.PUT("/tasks/:id", updateTask(client))
	router.DELETE("/tasks/:id", deleteTask(client))
}
