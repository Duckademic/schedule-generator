package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller[T any] interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	GetAll(*gin.Context)
}

type ServiceController[T any] interface {
	Create(*gin.Context) (*T, error)
	Update(*gin.Context) error
	Delete(*gin.Context) error
	GetAll() []T
}

type BasicController[T any] struct {
	serviceController ServiceController[T]
}

func (bm *BasicController[T]) Create(ctx *gin.Context) {
	teacher, err := bm.serviceController.Create(ctx)
	if err != nil {
		responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, teacher)
}

func (bm *BasicController[T]) Update(ctx *gin.Context) {
	err := bm.serviceController.Update(ctx)
	if err != nil {
		responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (bm *BasicController[T]) Delete(ctx *gin.Context) {
	err := bm.serviceController.Delete(ctx)
	if err != nil {
		responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (bm *BasicController[T]) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, bm.serviceController.GetAll())
}

func responseWithError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}
