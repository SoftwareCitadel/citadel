package controllers

import (
	"citadel/app/repositories"

	caesar "github.com/caesar-rocks/core"
)

type AnalyticsWebsitesController struct {
	*caesar.BaseResourceController

	websitesRepo *repositories.AnalyticsWebsitesRepository
}

func NewAnalyticsWebsitesController(websitesRepo *repositories.AnalyticsWebsitesRepository) *AnalyticsWebsitesController {
	return &AnalyticsWebsitesController{websitesRepo: websitesRepo}
}

// Define controller methods here
// func (c *AnalyticsWebsitesController) Index(ctx *caesar.Context) error {
// 	// Implement the controller method here

// 	return nil
// }
