package controller

import (
	"dbmsbackend/domain"
	"dbmsbackend/service"
	"dbmsbackend/util"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	productService *service.Product
}

func (controller *Product) Initialize(config *util.Config) (err error) {
	controller.productService = new(service.Product)
	err = controller.productService.Initialize(config)

	if err != nil {
		err = fmt.Errorf("initializing product controller: %w", err)
	}

	return err
}

func (controller *Product) New(c *gin.Context) {

	userID := c.GetInt("userID")

	type newProductDTO struct {
		Name          string     `json:"name"`
		Description   string     `json:"description"`
		Inventory     int        `json:"inventory"`
		Price         int        `json:"price"`
		StartSaleTime *time.Time `json:"startSaleTime"`
		EndSaleTime   *time.Time `json:"endSaleTime"`
		Picture       string     `json:"picture"`
	}

	dto := new(newProductDTO)

	err := c.BindJSON(dto)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity := controller.productService.New(dto.Name, dto.Description, dto.Picture, dto.Inventory, dto.Price, dto.StartSaleTime, dto.EndSaleTime, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, entity.ToRespDTO())
}

func (controller *Product) Update(c *gin.Context) {

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

	entity, err := controller.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	if entity.OwnerID != userID {
		c.JSON(http.StatusForbidden, util.GeneralAPIResponse{
			Status:  http.StatusForbidden,
			Message: "unauthorized operation",
		})
		return
	}

	type updateProductDTO struct {
		Name          string     `json:"name"`
		Description   string     `json:"description"`
		Inventory     int        `json:"inventory"`
		Price         int        `json:"price"`
		StartSaleTime *time.Time `json:"startSaleTime"`
		EndSaleTime   *time.Time `json:"endSaleTime"`
		Picture       string     `json:"picture"`
	}

	dto := new(updateProductDTO)

	err = c.BindJSON(dto)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity = controller.productService.Update(id, dto.Name, dto.Description, dto.Picture, dto.Inventory, dto.Price, dto.StartSaleTime, dto.EndSaleTime)

	c.JSON(http.StatusOK, entity.ToRespDTO())
}

func (controller *Product) Delete(c *gin.Context) {

	userID := c.GetInt("userID")

	rawID := c.Param("id")
	id, err := strconv.Atoi(rawID)

	type respDTO struct {
		Data string `json:"data"`
	}

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity, err := controller.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	if entity.OwnerID != userID {
		c.JSON(http.StatusForbidden, util.GeneralAPIResponse{
			Status:  http.StatusForbidden,
			Message: "unauthorized operation",
		})
		return
	}

	err = controller.productService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, respDTO{Data: fmt.Sprintf("product %v deleted", id)})

}

func (controller *Product) GetByID(c *gin.Context) {

	rawID := c.Param("id")

	id, err := strconv.Atoi(rawID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity, err := controller.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *entity.ToRespDTO())
}

func (controller *Product) Query(c *gin.Context) {

	conditions := make(map[string]interface{})

	productEntities, err := controller.productService.Query(conditions)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	var productDTOs []domain.ProductDTO

	for _, dto := range productEntities {
		productDTOs = append(productDTOs, *dto.ToDTO())
	}

	c.JSON(http.StatusOK, domain.ProductListRespDTO{Data: productDTOs})
}

func (controller *Product) NewImage(config *util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		type imageUrlDTO struct {
			Url string `json:"url"`
		}

		type imageUrlRespDTO struct {
			Data imageUrlDTO `json:"data"`
		}

		img, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		link, err := controller.productService.UploadToImgur(config.ImgurUploadURL, config.ImgurClientID, img)

		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, imageUrlRespDTO{Data: imageUrlDTO{Url: link}})
	}
}
