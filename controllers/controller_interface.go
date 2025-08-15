package controllers

import (
	"fmt"
	"net/http"

	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller[T comparable] interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	GetAll(*gin.Context)
}

type basicController[T comparable] struct {
	service       services.Service[T]
	objectParamId string
}

func (bc *basicController[T]) Create(ctx *gin.Context) {
	obj, err := bc.getObjectFromContext(ctx)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	pObj, err := bc.service.Create(*obj)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
	}

	ctx.JSON(http.StatusCreated, pObj)
}

func (bc *basicController[T]) Update(ctx *gin.Context) {
	obj, err := bc.getObjectFromContext(ctx)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	err = bc.service.Update(*obj)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
	}

	ctx.Status(http.StatusNoContent)
}

func (bc *basicController[T]) Delete(ctx *gin.Context) {
	objId, ok := ctx.Params.Get(bc.objectParamId)
	if !ok {
		types.ResponseWithError(ctx, http.StatusBadRequest, fmt.Errorf("missing teacher_id in URL parameters"))
		return
	}

	objUuid, err := uuid.Parse(objId)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	err = bc.service.Delete(objUuid)
	if err != nil {
		types.ResponseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (bc *basicController[T]) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, bc.service.GetAll())
}

func (bc *basicController[T]) getObjectFromContext(ctx *gin.Context) (*T, error) {
	var obj T
	err := ctx.ShouldBindBodyWithJSON(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
