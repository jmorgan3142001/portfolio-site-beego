package controllers

import (
	"encoding/json"
	"net/http"
	"portfolio-site/models"
	"time"

	"github.com/beego/beego/v2/server/web"
)

// --- Controller Definition ---

type PortfolioController struct {
    web.Controller
}

// --- Helper for GitHub API --- 
func getGithubStats(username string) models.GithubProfile {
    client := http.Client{
        Timeout: 2 * time.Second,
    }
    resp, err := client.Get("https://api.github.com/users/" + username)
    if err != nil {
        return models.GithubProfile{PublicRepos: 0, Login: "System Offline"}
    }
    defer resp.Body.Close()

    var profile models.GithubProfile
    json.NewDecoder(resp.Body).Decode(&profile)
    return profile
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
    
    // Page Context
    c.Data["Page"] = "home"

    // Dynamic Data from Models
    c.Data["TechSpecs"] = models.GetTechSpecs()
    c.Data["Experience"] = models.GetExperience()
    c.Data["Projects"] = models.GetProjects()

    // Render Configuration
    c.Layout = "layout.html"
    c.TplName = "index.html"
}

// Sub page declarations
func (c *PortfolioController) About() {
    c.Data["Title"] = "System Log"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "about"
    c.Data["Email"] = "jmorgan3142001@gmail.com"
    
    // Data from Models
    c.Data["GithubStats"] = getGithubStats("jmorgan3142001")
    c.Data["SystemModules"] = models.GetSystemModules()
    c.Data["Hardware"] = models.GetHardwareProfile()

    c.Layout = "layout.html"
    c.TplName = "about.html"
}

func (c *PortfolioController) Library() {
    c.Data["Title"] = "Data Archive"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "library"
    c.Data["Email"] = "jmorgan3142001@gmail.com"

    c.Data["Books"] = models.GetBooks()

    c.Layout = "layout.html"
    c.TplName = "library.html"
}

func (c *PortfolioController) Research() {
    c.Data["Title"] = "Research Lab | Jake Morgan"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "research"
    c.Data["Email"] = "jmorgan3142001@gmail.com"

    // Data from Models
    c.Data["Papers"] = models.GetResearchPapers()
    c.Data["Algos"] = models.GetAlgorithms()
    c.Data["Meta"] = models.GetResearchMeta()
    c.Data["Complexity"] = models.GetComplexityStats()
    c.Data["NextTopic"] = "BYZANTINE FAULT TOLERANCE"

    c.Layout = "layout.html"
    c.TplName = "research.html"
}