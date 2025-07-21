package controllers

import (
	"net/http"

	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LessonController interface {
	Controller[types.Lesson]
	SwapSlots(*gin.Context)
}

func NewLessonController(s services.LessonService) LessonController {
	lc := lessonController{
		basicController: basicController[types.Lesson]{
			service:       s,
			objectParamId: "lesson_id",
		},
		service: s,
	}

	return &lc
}

type lessonController struct {
	basicController[types.Lesson]
	service services.LessonService
}

func (lc *lessonController) SwapSlots(ctx *gin.Context) {
	type LessonPair struct {
		First  uuid.UUID `json:"first" binding:"required"`
		Second uuid.UUID `json:"second" binding:"required"`
	}

	var pair LessonPair
	err := ctx.ShouldBindBodyWithJSON(&pair)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	err = lc.service.SwapSlots(pair.First, pair.Second)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
