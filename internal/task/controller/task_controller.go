package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/httpres"
)

type ListTaskReq struct {
	BoardID string `form:"board_id" binding:"required"`
}

type taskController struct {
	conf    *config.Config
	service domain.TaskService
}

func NewTaskController(conf *config.Config, rg *gin.RouterGroup, service domain.TaskService) {
	c := &taskController{
		conf:    conf,
		service: service,
	}
	r := rg.Group("/task")
	r.GET("/", c.ListTask)
}

func (c *taskController) ListTask(ctx *gin.Context) {
	var listTaskReq ListTaskReq
	if err := ctx.ShouldBindQuery(&listTaskReq); err != nil {
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
