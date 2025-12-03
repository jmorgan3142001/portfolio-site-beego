package controllers

import (
	"portfolio-site/models"

	"github.com/beego/beego/v2/server/web"
)

// --- Controller Definition ---

type PortfolioController struct {
	web.Controller
}

// --- Request Handling ---

func (c *PortfolioController) Get() {
	// Static Personal Info
	c.Data["Title"] = "Jake Morgan | Software Engineer"
	c.Data["Name"] = "Jake Morgan"
	c.Data["Career"] = "Software Engineer"
	c.Data["Location"] = "Charleston, SC"
	c.Data["Email"] = "jmorgan3142001@gmail.com"
	c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
	c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jmorgan3142001"

	// Dynamic Data from Models
	c.Data["TechSpecs"] = models.GetTechSpecs()
	c.Data["Experience"] = models.GetExperience()
	c.Data["Projects"] = models.GetProjects()

	// Render Configuration
	c.Layout = "layout.html"
	c.TplName = "index.html"
}