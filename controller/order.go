package controller

import (
	"dbmsbackend/domain"
	"dbmsbackend/service"
	"dbmsbackend/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Order struct {
	orderService *service.Order
}

func (controller *Order) Initialize(config *util.Config) (err error) {
	controller.orderService = new(service.Order)
	err = controller.orderService.Initialize(config)

	if err != nil {
		err = fmt.Errorf("initializing order controller: %w", err)
	}

	return err
}

func (controller *Order) New(c *gin.Context) {

	userID := c.GetInt("userID")

	var dto []domain.OrderItemDTO

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	entity, err := controller.orderService.New(userID, dto)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, entity.ToRespDTO())
}

func (controller *Order) Update(c *gin.Context) {

	userID := c.GetInt("userID")

	rawID := c.Param("id")

	id, err := strconv.Atoi(rawID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity, err := controller.orderService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	if entity.Buyer.ID != userID {
		c.JSON(http.StatusForbidden, util.GeneralAPIResponse{
			Status:  http.StatusForbidden,
			Message: "unauthorized operation",
		})
		return
	}

	type updateOrderDTO struct {
		Products []domain.OrderItemDTO
	}

	dto := new(updateOrderDTO)
	err = c.BindJSON(dto)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
	}

	entity, err = controller.orderService.Update(id, dto.Products)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *entity.ToRespDTO())
}

func (controller *Order) Delete(c *gin.Context) {

	userID := c.GetInt("userID")

	rawID := c.Param("id")

	id, err := strconv.Atoi(rawID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity, err := controller.orderService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	if entity.Buyer.ID != userID {
		c.JSON(http.StatusForbidden, util.GeneralAPIResponse{
			Status:  http.StatusForbidden,
			Message: "unauthorized operation",
		})
		return
	}

	err = controller.orderService.Delete(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	type respDTO struct {
		Data string `json:"data"`
	}

	c.JSON(http.StatusOK, respDTO{
		Data: fmt.Sprintf("order %v deleted", id),
	})

}

func (controller *Order) GetByID(c *gin.Context) {

	userID := c.GetInt("userID")

	rawID := c.Param("id")

	id, err := strconv.Atoi(rawID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity, err := controller.orderService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	fmt.Println(entity.Buyer)
	fmt.Println(entity.Buyer.ID)

	if entity.Buyer.ID != userID {
		c.JSON(http.StatusForbidden, util.GeneralAPIResponse{
			Status:  http.StatusForbidden,
			Message: "unauthorized operation",
		})
		return
	}

	fmt.Println("AAA")

	c.JSON(http.StatusOK, *entity.ToRespDTO())
}

func (controller *Order) Query(c *gin.Context) {

	userID := c.GetInt("userID")

	orderEntities, err := controller.orderService.GetByBuyerID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	var orderDTOs []domain.OrderDTO

	for _, dto := range orderEntities {
		orderDTOs = append(orderDTOs, *dto.ToDTO())
	}

	c.JSON(http.StatusOK, domain.OrderListRespDTO{Data: orderDTOs})
}
