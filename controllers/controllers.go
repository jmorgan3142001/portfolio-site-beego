package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"portfolio-site/models"
	"strings"
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
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"
    
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
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"
    
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
    c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"

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
    c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"

    // Data from Models
    c.Data["Papers"] = models.GetResearchPapers()
    c.Data["Algos"] = models.GetAlgorithms()
    c.Data["Meta"] = models.GetResearchMeta()
    c.Data["NextTopic"] = "BYZANTINE FAULT TOLERANCE"

    c.Layout = "layout.html"
    c.TplName = "research.html"
}

func (c *PortfolioController) Terminal() {
    c.Data["Title"] = "Interactive Shell"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "terminal"
    c.Data["Email"] = "jmorgan3142001@gmail.com"
    c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"

    c.Layout = "layout.html"
    c.TplName = "terminal.html"
}

func (c *PortfolioController) Challenge() {
    c.Data["Title"] = "Skill Check"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "challenges"
    c.Data["Email"] = "jmorgan3142001@gmail.com"
    c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"

    c.Data["Challenges"] = models.GetChallenges()
    
    c.Layout = "layout.html"
    c.TplName = "challenges.html"
}

// RunCode sends the user's submission to the Piston Execution Engine
func (c *PortfolioController) RunCode() {
    // Parse the incoming JSON payload
    var req struct {
        ChallengeID int    `json:"challenge_id"`
        UserCode    string `json:"user_code"`
    }
    
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
        c.Data["json"] = map[string]interface{}{"passed": false, "output": "System Error: Malformed JSON."}
        c.ServeJSON()
        return
    }

    // Validate Challenge Existence
    if _, err := models.GetChallengeById(req.ChallengeID); err != nil {
        c.Data["json"] = map[string]interface{}{
            "passed": false, 
            "output": "System Error: Challenge ID not found in database.",
        }
        c.ServeJSON()
        return
    }

    // Fetch Test Cases
    testCases := models.GetTestCases(req.ChallengeID)
    if len(testCases) == 0 {
        c.Data["json"] = map[string]interface{}{"passed": false, "output": "System Error: No test cases found for this challenge."}
        c.ServeJSON()
        return
    }

    // Configuration - using Python 3.10 via Piston API
    version := "3.10.0" 
    
    allPassed := true
    var outputLog strings.Builder

    for _, tc := range testCases {
        
        // Wrap User Code - append a print statement that calls their function with our hidden input
        // Expected format - print(solve(ARGUMENTS))
        fullCode := fmt.Sprintf("%s\n\nprint(solve(%s))", req.UserCode, tc.InputArgs)

        // Construct Piston Payload
        pistonReq := models.PistonRequest{
            Language: "python",
            Version:  version,
            Files: []models.PistonFile{
                {Name: "main.py", Content: fullCode},
            },
        }

        reqBody, _ := json.Marshal(pistonReq)
        
        // Execute Request
        resp, err := http.Post("https://emkc.org/api/v2/piston/execute", "application/json", bytes.NewBuffer(reqBody))
        if err != nil {
            c.Data["json"] = map[string]interface{}{"error": "External Execution Engine (Piston) is unreachable."}
            c.ServeJSON()
            return
        }
        defer resp.Body.Close()

        // Parse Response
        var pistonResp models.PistonResponse
        if err := json.NewDecoder(resp.Body).Decode(&pistonResp); err != nil {
            outputLog.WriteString("Error parsing execution response.\n")
            allPassed = false
            continue
        }

        // Clean Output Strings
        actualOutput := strings.TrimSpace(pistonResp.Run.Stdout)
        expected := strings.TrimSpace(tc.ExpectedOutput)

        // Check for Runtime/Syntax Errors
        if pistonResp.Run.Stderr != "" {
            allPassed = false
            // Remove Piston's internal file paths to make the error cleaner for the user
            cleanErr := strings.Replace(pistonResp.Run.Stderr, "/piston/jobs/", "", -1)
            outputLog.WriteString(fmt.Sprintf("Runtime Error on input (%s):\n%s\n", tc.InputArgs, cleanErr))
            break // Stop on first error to prevent log spam
        }

        // Compare Results
        if actualOutput == expected {
            outputLog.WriteString(fmt.Sprintf("✓ PASS: Input(%s) -> Output(%s)\n", tc.InputArgs, actualOutput))
        } else {
            allPassed = false
            outputLog.WriteString(fmt.Sprintf("✗ FAIL: Input(%s)\n  Expected: %s\n  Got:      %s\n", tc.InputArgs, expected, actualOutput))
        }
    }

    // Return Result
    c.Data["json"] = map[string]interface{}{
        "passed": allPassed,
        "output": outputLog.String(),
    }
    c.ServeJSON()
}