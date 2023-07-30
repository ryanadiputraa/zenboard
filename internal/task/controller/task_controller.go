package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/httpres"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
)

type ListTaskReq struct {
	BoardID string `form:"board_id" binding:"required"`
}

type CreateTaskReq struct {
	BoardID  string `json:"board_id" binding:"required"`
	TaskName string `json:"task_name" binding:"required"`
}

type taskController struct {
	conf         *config.Config
	service      domain.TaskService
	boardService domain.BoardService
}

func NewTaskController(conf *config.Config, rg *gin.RouterGroup, service domain.TaskService, boardService domain.BoardService) {
	c := &taskController{
		conf:         conf,
		service:      service,
		boardService: boardService,
	}
	r := rg.Group("/tasks")
	r.GET("/", c.ListTask)
	r.POST("/", c.CreateTask)
}

func (c *taskController) ListTask(ctx *gin.Context) {
	userID, err := jwt.ExtractUserID(ctx, c.conf.JWT)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	var listTaskReq ListTaskReq
	if err := ctx.ShouldBindQuery(&listTaskReq); err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	isAuthorized, err := c.boardService.CheckIsUserAuthorized(ctx, listTaskReq.BoardID, userID)
	if err != nil || !isAuthorized {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	tasks, err := c.service.ListBoardTasks(ctx, listTaskReq.BoardID)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusOK, tasks)
}

func (c *taskController) CreateTask(ctx *gin.Context) {
	userID, err := jwt.ExtractUserID(ctx, c.conf.JWT)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	var createTaskReq CreateTaskReq
	if err := ctx.ShouldBindJSON(&createTaskReq); err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	isAuthorized, err := c.boardService.CheckIsUserAuthorized(ctx, createTaskReq.BoardID, userID)
	if err != nil || !isAuthorized {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	task, err := c.service.AddBoardTask(ctx, createTaskReq.BoardID, createTaskReq.TaskName)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusOK, task)
}
