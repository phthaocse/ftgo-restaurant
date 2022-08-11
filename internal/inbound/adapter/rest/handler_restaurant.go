package rest

import (
	"ftgo-restaurant/internal/core/model"
	reqModel "ftgo-restaurant/internal/inbound/adapter/rest/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (gs *ginServer) createRestaurant(c *gin.Context) {
	restaurantReq := reqModel.Restaurant{}
	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	var menu []*model.MenuItem
	for _, item := range restaurantReq.Menu {
		menu = append(menu, &model.MenuItem{
			Price: item.Price,
			Name:  item.Name,
		})
	}
	gs.BusinessService.RestaurantService.Create(model.Restaurant{
		Name:    restaurantReq.Name,
		Address: restaurantReq.Address,
		Menu:    menu,
	})
}
