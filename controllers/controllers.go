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
    client := http.Client{Timeout: 2 * time.Second}
    var profile models.GithubProfile

    // Fetch User Profile
    resp, err := client.Get("https://api.github.com/users/" + username)
    if err == nil {
        defer resp.Body.Close()
        json.NewDecoder(resp.Body).Decode(&profile)
    } else {
        profile.Login = "System Offline"
    }

    // Fetch Public Events
    respEvents, errEvents := client.Get("https://api.github.com/users/" + username + "/events/public")
    if errEvents == nil {
        defer respEvents.Body.Close()
        
        var rawEvents []models.GithubEvent
        
        if json.NewDecoder(respEvents.Body).Decode(&rawEvents) == nil {
            for _, event := range rawEvents {
                if event.Type == "PushEvent" {
                    
                    // Get the Code
                    if event.Payload.Head != "" {
                        // Use the first 7 characters for the "Short SHA"
                        if len(event.Payload.Head) >= 7 {
                            profile.LatestCommitCode = event.Payload.Head[:7]
                        } else {
                            profile.LatestCommitCode = event.Payload.Head
                        }
                        
                        // Construct URL
                        if event.Repo.Name != "" {
                            profile.LatestCommitUrl = "https://github.com/" + event.Repo.Name + "/commit/" + event.Payload.Head
                        }
                    }
                    
                    break 
                }
            }
        }
    }

	// Fallback if nothing is returned
    if profile.LatestCommitUrl == "" {
        profile.LatestCommitCode = "N/A"
        profile.LatestCommitUrl = "#"
    }

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
    c.Data["Logs"] = models.GetAccessLogs()

    // Render Configuration
    c.Layout = "layout.html"
    c.TplName = "index.html"
}

func (c *PortfolioController) SubmitLog() {
    name := c.GetString("username")
    message := c.GetString("payload")
    userAgent := c.Ctx.Input.UserAgent()

    if name != "" && message != "" {
        models.AddAccessLog(name, message, userAgent)
    }

    c.Redirect("/", 302)
}

// Sub page declarations
func (c *PortfolioController) About() {
    c.Data["Title"] = "System Log"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "about"
    c.Data["Email"] = "jmorgan3142001@gmail.com"
	c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
    
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
	c.Data["NextBooks"] = models.GetNextReads()
	c.Data["Resources"] = models.GetDigitalResources()

    c.Layout = "layout.html"
    c.TplName = "library.html"
}

func (c *PortfolioController) Research() {
    c.Data["Title"] = "Research Lab"
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

func (c *PortfolioController) Terminal() {
    c.Data["Title"] = "Interactive Shell"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "terminal"
    c.Data["Email"] = "jmorgan3142001@gmail.com"

    c.Layout = "layout.html"
    c.TplName = "terminal.html"
}