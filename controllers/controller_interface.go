package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller[T any] interface {
	Create(*gin.Context) (*T, error)
	Update(*gin.Context) error
	Delete(*gin.Context) error
	GetAll() []T
}

func NewBasicMiddleware[T any](controller Controller[T]) *BasicMiddleware[T] {
	return &BasicMiddleware[T]{
		controller: controller,
	}
}

type BasicMiddleware[T any] struct {
	controller Controller[T]
}

func (bm *BasicMiddleware[T]) Create(ctx *gin.Context) {
	teacher, err := bm.controller.Create(ctx)
	if err != nil {
		responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, teacher)
}

func (bm *BasicMiddleware[T]) Update(ctx *gin.Context) {
	err := bm.controller.Update(ctx)
	if err != nil {
		responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (bm *BasicMiddleware[T]) Delete(ctx *gin.Context) {
	err := bm.controller.Delete(ctx)
	if err != nil {
		responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (bm *BasicMiddleware[T]) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, bm.controller.GetAll())
}

func responseWithError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}
