package controllers

import "github.com/gin-gonic/gin"

type Controller[T any] interface {
	Create(*gin.Context) (*T, error)
	Update(*gin.Context) error
	Delete(*gin.Context) error
	GetAll() []T
}
