package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
    c.Data["Title"] = "User Log"
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

func (c *PortfolioController) Directory() {
    c.Data["Title"] = "File Server"
    c.Data["Name"] = "Jake Morgan"
    c.Data["Page"] = "directory"
    c.Data["Email"] = "jmorgan3142001@gmail.com"
    c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
    c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jake-morgan-/"

    // Library Data
    c.Data["Books"] = models.GetBooks()
    c.Data["NextBooks"] = models.GetNextReads()
    c.Data["Resources"] = models.GetDigitalResources()
    c.Data["Creators"] = models.GetCreators()

    // Research Data
    c.Data["Papers"] = models.GetResearchPapers()
    c.Data["Algos"] = models.GetAlgorithms()
    c.Data["Meta"] = models.GetResearchMeta()
    c.Data["NextTopic"] = "BYZANTINE FAULT TOLERANCE"

    c.Layout = "layout.html"
    c.TplName = "directory.html"
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
    // 1. Parse Payload
    var req struct {
        ChallengeID int    `json:"challenge_id"`
        UserCode    string `json:"user_code"`
    }
    
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
        c.Data["json"] = map[string]interface{}{"passed": false, "output": "System Error: Bad JSON"}
        c.ServeJSON()
        return
    }

    // 2. Fetch Challenge & Test Cases
    challenge, err := models.GetChallengeById(req.ChallengeID) 
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "passed": false, 
            "output": "System Error: Challenge ID not found in database.",
        }
        c.ServeJSON()
        return
    }

    testCases := models.GetTestCases(req.ChallengeID)
    if len(testCases) == 0 {
        c.Data["json"] = map[string]interface{}{"passed": false, "output": "System Error: No test cases found in DB."}
        c.ServeJSON()
        return
    }

    // 3. Prepare Execution
    version := "3.10.0" 
    allPassed := true
    var outputLog strings.Builder

    // Default to 'solve' if database field is empty
    funcName := challenge.FunctionName
    if funcName == "" {
        funcName = "solve"
    }

    // 4. Execution Loop
    for i, tc := range testCases {
        
        // --- RATE LIMIT FIX ---
        if i > 0 {
            time.Sleep(250 * time.Millisecond)
        }

        fullCode := fmt.Sprintf("%s\n\nprint(%s(%s))", req.UserCode, funcName, tc.InputArgs)

        pistonReq := models.PistonRequest{
            Language: "python",
            Version:  version,
            Files: []models.PistonFile{
                {Name: "main.py", Content: fullCode},
            },
        }

        reqBody, _ := json.Marshal(pistonReq)
        
        client := &http.Client{Timeout: 10 * time.Second}
        r, _ := http.NewRequest("POST", "https://emkc.org/api/v2/piston/execute", bytes.NewBuffer(reqBody))
        r.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(r)
        
        if err != nil {
            c.Data["json"] = map[string]interface{}{"error": "Execution Engine Offline"}
            c.ServeJSON()
            return
        }
        defer resp.Body.Close()

        var bodyBytes []byte
        if resp.Body != nil {
            bodyBytes, _ = io.ReadAll(resp.Body)
        }

        // Check for API errors
        if resp.StatusCode != 200 {
            outputLog.WriteString(fmt.Sprintf("API Error (Case %d): %s\n", i+1, string(bodyBytes)))
            allPassed = false
            break 
        }

        var pistonResp models.PistonResponse
        if err := json.Unmarshal(bodyBytes, &pistonResp); err != nil {
            allPassed = false
            outputLog.WriteString("Error parsing execution response.\n")
            continue
        }

        actualOutput := strings.TrimSpace(pistonResp.Run.Stdout)
        expected := strings.TrimSpace(tc.ExpectedOutput)

        if pistonResp.Run.Stderr != "" {
            allPassed = false
            cleanErr := strings.Replace(pistonResp.Run.Stderr, "/piston/jobs/", "", -1)
            outputLog.WriteString(fmt.Sprintf("ERROR on input (%s):\n%s\n", tc.InputArgs, cleanErr))
            break 
        }

        if actualOutput == expected {
            outputLog.WriteString(fmt.Sprintf("✓ PASS: Input(%s) -> Output(%s)\n", tc.InputArgs, actualOutput))
        } else {
            allPassed = false
            outputLog.WriteString(fmt.Sprintf("✗ FAIL: Input(%s)\n  Expected: %s\n  Got:      %s\n", tc.InputArgs, expected, actualOutput))
        }
    }

    c.Data["json"] = map[string]interface{}{
        "passed": allPassed,
        "output": outputLog.String(),
    }
    c.ServeJSON()
}