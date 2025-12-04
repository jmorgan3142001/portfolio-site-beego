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

    // 1. Fetch User Profile (for Repo count)
    resp, err := client.Get("https://api.github.com/users/" + username)
    if err == nil {
        defer resp.Body.Close()
        json.NewDecoder(resp.Body).Decode(&profile)
    } else {
        profile.Login = "System Offline"
    }

    // 2. Fetch Public Events (for Latest Commit Link)
    respEvents, errEvents := client.Get("https://api.github.com/users/" + username + "/events/public")
    if errEvents == nil {
        defer respEvents.Body.Close()
        
        var rawEvents []models.GithubEvent // Using the internal struct defined in models
        
        // We decode into the helper struct we defined
        if json.NewDecoder(respEvents.Body).Decode(&rawEvents) == nil {
            for _, event := range rawEvents {
                // We only want PushEvents (commits)
                if event.Type == "PushEvent" {
                    
                    // 1. Get the Message
                    if len(event.Payload.Commits) > 0 {
                        profile.LatestCommitMsg = event.Payload.Commits[0].Message
                    } else {
                        profile.LatestCommitMsg = "Update repository"
                    }
                    
                    // Truncate message if it's too long
                    if len(profile.LatestCommitMsg) > 50 {
                        profile.LatestCommitMsg = profile.LatestCommitMsg[:47] + "..."
                    }

                    // 2. Construct the Public URL
                    // Format: https://github.com/{Owner}/{Repo}/commit/{SHA}
                    if event.Repo.Name != "" && event.Payload.Head != "" {
                        profile.LatestCommitUrl = "https://github.com/" + event.Repo.Name + "/commit/" + event.Payload.Head
                    }
                    
                    break // Stop after finding the most recent one
                }
            }
        }
    }

    // Fallback if no public commits found
    if profile.LatestCommitUrl == "" {
        profile.LatestCommitMsg = "No recent public commits"
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