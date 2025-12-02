package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type PortfolioController struct {
	web.Controller
}

// Project struct to hold data for your portfolio items
type Project struct {
	Title       string
	Description string
	TechStack   []string
	Link        string
	Repo        string
}

func (c *PortfolioController) Get() {
	// Title and Personal Info
	c.Data["Title"] = "Jake Morgan | Software Engineer"
	c.Data["Name"] = "Jake Morgan"
	c.Data["Headline"] = "Innovative Software Engineer"
	
	// Social Links
	c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
	c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jmorgan3142001" 
	c.Data["Email"] = "jmorgan3142001@gmail.com" 

	// Project Data
	projects := []Project{
		{
			Title:       "Lawless Lowcountry Living",
			Description: "A responsive real estate platform integrating custom navigation, contact forms, and a modern UI for the Charleston market.",
			TechStack:   []string{"HTML/CSS", "JavaScript", "AWS"},
			Link:        "http://lawlesslowcountryliving.com/",
			Repo:        "https://github.com/jmorgan3142001/LawlessLowcountyLiving",
		},
		{
			Title:       "Go/Beego Portfolio",
			Description: "A scalable web application built with the Beego framework, utilizing MVC architecture and Bootstrap 5.",
			TechStack:   []string{"Go", "Beego", "Bootstrap 5"},
			Link:        "#",
			Repo:        "https://github.com/jmorgan3142001",
		},
		{
			Title:       "Data Structures & Algos",
			Description: "Optimized implementations of core computer science algorithms and data structures.",
			TechStack:   []string{"Python", "C++"},
			Link:        "#",
			Repo:        "https://github.com/jmorgan3142001",
		},
	}

	// Pass data to the template
	c.Data["Projects"] = projects
	
	// Render view
	c.TplName = "index.html"
}