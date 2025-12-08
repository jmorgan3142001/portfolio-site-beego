package models

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

// ===================================================================================
// SECTION 1: DATABASE MODELS & SCHEMA CONFIGURATION
// ===================================================================================

// --- Access Log Model ---
type AccessLog struct {
    Id        int       `orm:"auto"`
    Name      string    `orm:"size(50)"`
    Message   string    `orm:"type(text)"`
    Signature string    `orm:"size(10)"`
    Terminal  string    `orm:"size(30)"`
    ProcessID int
    Created   time.Time `orm:"auto_now_add;type(datetime)"`
}

func (u *AccessLog) TableName() string {
    return "access_log"
}

// --- Challenge Model ---
type Challenge struct {
    Id          int         `orm:"auto"`
    Title       string      `orm:"size(255)"`
    Description string      `orm:"type(text)"`
	InputHint   string      `orm:"size(255)"`
    Difficulty  string      `orm:"size(50)"`
    Category    string      `orm:"size(50)"`
    Type        string      `orm:"size(50)"`
    Language    string      `orm:"size(50);null"`
    StarterCode string      `orm:"type(text);null"`
    TestCases   []*TestCase `orm:"reverse(many)"`
}

// Force table name to singular 'challenge' to match DB creation default
func (u *Challenge) TableName() string {
    return "challenge"
}

// --- Test Case Model ---
type TestCase struct {
    Id             int        `orm:"auto"`
    Challenge      *Challenge `orm:"rel(fk);on_delete(cascade)"`
    InputArgs      string     `orm:"type(text)"`
    ExpectedOutput string     `orm:"type(text)"`
}

// Force table name to singular 'test_case' to match DB creation default
func (u *TestCase) TableName() string {
    return "test_case"
}

// ===================================================================================
// SECTION 2: DATABASE INITIALIZATION
// ===================================================================================

func init() {
    // 1. Register Models
    orm.RegisterModel(new(AccessLog), new(Challenge), new(TestCase))

    // 2. Register Driver
    orm.RegisterDriver("postgres", orm.DRPostgres)

    // 3. Connect to Database
    dbUrl := os.Getenv("DATABASE_URL")
    err := orm.RegisterDataBase("default", "postgres", dbUrl)
    if err != nil {
        panic(fmt.Errorf("failed to register database: %v", err))
    }

    // 4. Run Synchronize (Create tables if they don't exist)
    // force=false (do not drop tables), verbose=true
    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        fmt.Println("Database Sync Error:", err)
    }

    // 5. Seed Initial Data
    SeedChallenges()
}

// ===================================================================================
// SECTION 3: SEEDING LOGIC
// ===================================================================================

func SeedChallenges() {
    o := orm.NewOrm()
    fmt.Println("Running Seeder for Challenges...")

    // Definition struct to keep the loop clean
    type SeedData struct {
        Title       string
        Desc        string
        Diff        string
        Cat         string
        Code        string
        Tests       []TestCase
    }

    seeds := []SeedData{
        // --- 1. EASY (BUG FIX) - NCI ---
        {
            Title: "Legacy Timestamp Bug",
            Diff:  "Easy",
            Cat:   "NCI / Debugging",
            Desc:  "In my work at NCI, we process legacy CEA-608 caption files. For these files, and any caption file, correct timing and time formatting is of paramount importance.\n\nThe output and logic of this code is incorrect. Find the bug and fix the output.",
            Code:  `def convert_time(ms):
    """
    Converts milliseconds to formatted string MM:SS:mmm
    Args:
        ms (int): Total milliseconds (e.g., 61500)
    Returns:
        str: Formatted string (e.g., "01:05:000")
    """
    # BUG: Code fails when ms > 60000. Not sure why?
    minutes = ms // 60
    seconds = (ms // 1000) % 60
    rem_ms = ms % 1000
    
    return f"{minutes:02}:{seconds:02}:{rem_ms:03}"`,
            Tests: []TestCase{
                {InputArgs: "65000", ExpectedOutput: "01:05:000"}, 
                {InputArgs: "125500", ExpectedOutput: "02:05:500"}, 
            },
        },

        // --- 2. EASY (VALIDATION) - UG ---
        {
            Title: "Payment Input Validator",
            Diff:  "Easy",
            Cat:   "UG / Validation",
            Desc:  "At Uncommon Giving, ensuring data integrity before hitting the payment gateway is critical. If not careful at our work, bad data can cause downstream API failures and reconciliation issues.\n\nThe current validator is too permissive. Update the logic to reject negative numbers, non-numeric characters, and excessive decimal precision.",
            Code:  `def validate_currency(amount_str):
    """
    Validates if a string is a valid currency amount.
    Rules:
    1. Must be a positive number.
    2. No more than 2 decimal places.
    3. No currency symbols ($).
    
    Args: amount_str (str)
    Returns: bool
    """
    # TODO: Implement proper validation of value.
    return False`,
            Tests: []TestCase{
                {InputArgs: "'10.50'", ExpectedOutput: "True"},
                {InputArgs: "'-5.00'", ExpectedOutput: "False"},
                {InputArgs: "'10.555'", ExpectedOutput: "False"},
                {InputArgs: "'abc'", ExpectedOutput: "False"},
            },
        },

        // --- 3. MEDIUM (BUG FIX) - NCI ---
        {
            Title: "FCC Compliance Splitter",
            Diff:  "Medium",
            Cat:   "NCI / Strings",
            Desc:  "At NCI, broadcast captions for many of our clients must adhere to strict FCC regulations (max 32 characters per line). Accessibility compliance is non-negotiable and must always be accounted and tested for.\n\nThe current implementation occasionally breaks words in half to force a fit or miscalculates line length. Fix the logic to ensure cleaner line breaks.",
            Code:  `def wrap_text_lines(text):
    """
    Splits text into lines of max 32 chars without cutting words.
    Args: text (str)
    Returns: List[str]
    """
    words = text.split()
    lines = []
    current_line = ""

    for word in words:
        # BUG: Logic allows lines > 32 chars in specific edge cases.
        if len(current_line) + len(word) > 32:
            lines.append(current_line)
            current_line = word
        else:
            current_line += word 
            
    lines.append(current_line)
    return lines`,
            Tests: []TestCase{
                {InputArgs: "'This is a test of the emergency broadcast system'", ExpectedOutput: "['This is a test of the emergency', 'broadcast system']"},
                {InputArgs: "'A short line'", ExpectedOutput: "['A short line']"},
            },
        },

        // --- 4. MEDIUM (ALGO) - UG ---
        {
            Title: "The 'Lost Penny' Problem",
            Diff:  "Medium",
            Cat:   "UG / FinTech",
            Desc:  "At Uncommon Giving, bundled donations can be split among multiple charities as specified by the user. Financial accuracy must always be kept, and standard division can result in 'lost pennies' (e.g., $100 / 3).\n\nImplement a distribution algorithm that splits `total_cents` among `n` recipients, ensuring every penny is accounted for.",
            Code:  `def distribute_pennies(total_cents, n_charities):
    """
    Splits total_cents among n_charities.
    Args:
        total_cents (int): Total amount (e.g. 10000)
        n_charities (int): Number of recipients (e.g. 3)
    Returns:
        List[int]: List of amounts that sum exactly to total_cents.
                   (e.g. [3334, 3333, 3333])
    """
    # TODO: Implement fair distribution logic.
    return []`,
            Tests: []TestCase{
                {InputArgs: "10000, 3", ExpectedOutput: "[3334, 3333, 3333]"},
                {InputArgs: "100, 6", ExpectedOutput: "[17, 17, 17, 17, 16, 16]"},
            },
        },

        // --- 5. HARD (GRAPH) - Scheduler ---
        {
            Title: "Job Dependency Cascade",
            Diff:  "Hard",
            Cat:   "Systems / Graph",
            Desc:  "We contribute to the open-source Django5 Scheduler used by NCI. Preventing zombie processes and handling cascading failures is critical for system stability.\n\nImplement a dependency resolver that identifies all downstream jobs that must be cancelled when a parent job fails.",
            Code:  `def find_impacted_jobs(deps, failed_job):
    """
    Finds all downstream jobs affected by a failure.
    Args:
        deps (Dict[str, List[str]]): Key is a job, value is list of jobs that depend on it.
                                     e.g. {'A': ['B'], 'B': ['C']} means A -> B -> C
        failed_job (str): The ID of the job that crashed.
    Returns:
        List[str]: Sorted list of all impacted jobs (including failed_job).
    """
    # TODO: Identify all downstream dependencies.
    return []`,
            Tests: []TestCase{
                {InputArgs: "{'A': ['B'], 'B': ['C'], 'C': []}, 'A'", ExpectedOutput: "['A', 'B', 'C']"},
                {InputArgs: "{'A': ['B', 'C'], 'B': [], 'C': ['D'], 'D': []}, 'A'", ExpectedOutput: "['A', 'B', 'C', 'D']"},
            },
        },

        // --- 6. HARD (DATA) - UG ---
        {
            Title: "Async Ledger Reconciliation",
            Diff:  "Hard",
            Cat:   "UG / Distributed",
            Desc:  "At Uncommon Giving, we process asynchronous webhooks from payment processors. Ledger accuracy must be kept at all times with no errors, but events can often arrive out of order (e.g., REVERSE before DEPOSIT).\n\nImplement reconciliation logic that calculates the correct final balance regardless of event arrival order.",
            Code:  `def reconcile_ledger(events):
    """
    Calculates final balance from out-of-order stream.
    Args:
        events (List[Dict]): List of dicts. 
            e.g. {"id": 1, "type": "DEPOSIT", "amount": 100}
                 {"id": 2, "type": "REVERSE", "ref_id": 1}
    Returns:
        int: Final Balance
    """
    # TODO: Reconcile the ledger.
    return 0`,
            Tests: []TestCase{
                {InputArgs: "[{'id':1, 'type':'DEPOSIT', 'amount':100}, {'id':2, 'type':'WITHDRAW', 'amount':50}]", ExpectedOutput: "50"},
                {InputArgs: "[{'id':2, 'type':'REVERSE', 'ref_id':1}, {'id':1, 'type':'DEPOSIT', 'amount':100}]", ExpectedOutput: "0"},
            },
        },
    }

    // --- UPSERT LOGIC ---
    for _, s := range seeds {
        c := Challenge{Title: s.Title}
        
        // 1. Check existence
        err := o.Read(&c, "Title")

        // 2. Set/Update fields
        c.Description = s.Desc
        c.Difficulty = s.Diff
        c.Category = s.Cat
        c.Type = "CODE"
        c.Language = "python"
        c.StarterCode = s.Code

        // 3. Persist Challenge
        if err == orm.ErrNoRows {
            if _, err := o.Insert(&c); err != nil {
                fmt.Printf("Error inserting %s: %v\n", s.Title, err)
                continue
            }
        } else {
            if _, err := o.Update(&c); err != nil {
                fmt.Printf("Error updating %s: %v\n", s.Title, err)
                continue
            }
        }

        // 4. Reset Test Cases (Delete old, Insert new)
        o.QueryTable("test_case").Filter("challenge_id", c.Id).Delete()
        
        var newTests []TestCase
        for _, t := range s.Tests {
            t.Challenge = &c
            newTests = append(newTests, t)
        }
        o.InsertMulti(len(newTests), newTests)
    }
    
    fmt.Println("Seed Complete.")
}

// ===================================================================================
// SECTION 4: DATABASE RETRIEVAL LOGIC
// ===================================================================================

func GetChallenges() []Challenge {
    o := orm.NewOrm()
    var challenges []Challenge
    o.QueryTable("challenge").OrderBy("id").All(&challenges)
    return challenges
}

func GetChallengeById(id int) (Challenge, error) {
    o := orm.NewOrm()
    var challenge Challenge
    err := o.QueryTable("challenge").Filter("Id", id).One(&challenge)
    return challenge, err
}

func GetTestCases(challengeId int) []TestCase {
    o := orm.NewOrm()
    var cases []TestCase
    o.QueryTable("test_case").Filter("challenge_id", challengeId).All(&cases)
    return cases
}

func GetAccessLogs() []AccessLog {
    o := orm.NewOrm()
    var logs []AccessLog
    o.QueryTable("access_log").OrderBy("-created").All(&logs)
    return logs
}

func AddAccessLog(name, message, userAgent string) error {
    if CheckToxicity(name) || CheckToxicity(message) {
        return errors.New("input rejected: content detected as toxic or aggressive")
    }

    o := orm.NewOrm()

    // Generate Signature
    h := sha1.New()
    h.Write([]byte(message + time.Now().String()))
    bs := h.Sum(nil)
    sig := fmt.Sprintf("%x", bs)[:7]

    // Generate PID
    pid := rand.Intn(8999) + 1000

    // Parse Terminal
    terminal := "Unknown/Term"
    if len(userAgent) > 0 {
        if len(userAgent) > 30 {
            terminal = userAgent[:27] + "..."
        } else {
            terminal = userAgent
        }
    }

    log := AccessLog{
        Name:      name,
        Message:   message,
        Signature: sig,
        Terminal:  terminal,
        ProcessID: pid,
    }

    _, err := o.Insert(&log)
    return err
}

// ===================================================================================
// SECTION 5: EXTERNAL API HELPERS
// ===================================================================================

// Check for toxicity using the Perspective API
func CheckToxicity(content string) bool {
    apiKey := os.Getenv("PERSPECTIVE_API_KEY")
    url := "https://commentanalyzer.googleapis.com/v1alpha1/comments:analyze?key=" + apiKey

    requestBody, err := json.Marshal(map[string]interface{}{
        "comment": map[string]string{
            "text": content,
        },
        "requestedAttributes": map[string]interface{}{
            "TOXICITY": map[string]interface{}{},
        },
    })

    if err != nil {
        return false // Handle JSON marshaling error
    }

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    var result struct {
        AttributeScores struct {
            Toxicity struct {
                SummaryScore struct {
                    Value float64 `json:"value"`
                } `json:"summaryScore"`
            } `json:"TOXICITY"`
        } `json:"attributeScores"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false
    }

    return result.AttributeScores.Toxicity.SummaryScore.Value > 0.75
}

// --- Piston API Payloads ---

type PistonRequest struct {
    Language string        `json:"language"`
    Version  string        `json:"version"`
    Files    []PistonFile  `json:"files"`
    Stdin    string        `json:"stdin"`
    Args     []string      `json:"args"`
}

type PistonFile struct {
    Name    string `json:"name"`
    Content string `json:"content"`
}

type PistonResponse struct {
    Run struct {
        Stdout string `json:"stdout"`
        Stderr string `json:"stderr"`
        Code   int    `json:"code"`
    } `json:"run"`
}

// --- GitHub API Models ---

type GithubProfile struct {
    PublicRepos      int    `json:"public_repos"`
    Login            string `json:"login"`
    LatestCommitCode string `json:"latest_commit_code"`
    LatestCommitUrl  string `json:"latest_commit_url"`
}

type GithubEvent struct {
    Type string `json:"type"`
    Repo struct {
        Name string `json:"name"`
    } `json:"repo"`
    Payload struct {
        Head    string `json:"head"`
        Commits []struct {
            Message string `json:"message"`
        } `json:"commits"`
    } `json:"payload"`
}

// ===================================================================================
// SECTION 6: STATIC PORTFOLIO DATA & STRUCTS
// ===================================================================================

type TechItem struct {
    Category string
    Items    []string
}

type Project struct {
    ID          string
    Title       string
    Description string
    Tags        []string
    Link        string
    Icon        string
}

type Experience struct {
    Company     string
    Role        string
    Duration    string
    Description string
    Tags        []string
}

type SystemModule struct {
    ID          string
    Title       string
    Icon        string
    Description string
    Progress    int
    Tags        []string
}

type HardwareInfo struct {
    Category string
    Item     string
}

type Book struct {
    Title       string
    Author      string
    Description string
    Type        string
    Link        string
}

type ResearchPaper struct {
    Title string
    Topic string
    Note  string
    Link  string
}

type Algorithm struct {
    Name        string
    Lang        string
    Description string
    Link        string
}

type ResearchMeta struct {
    Label  string
    Value  string
    Status string
}

type ResourceItem struct {
    Title       string
    Description string
    Link        string
    Icon        string
}

// --- Static Data Retrieval ---

func GetTechSpecs() []TechItem {
    return []TechItem{
        {Category: "Languages", Items: []string{"Python", "C#", "JavaScript", "SQL", "PHP", "C++", "Go"}},
        {Category: "Frameworks", Items: []string{"Django", "Angular", ".NET Core", "Flutter", "Beego"}},
        {Category: "Databases", Items: []string{"PostgreSQL", "Oracle", "MySQL", "SQLite", "Distributed DB"}},
        {Category: "Infrastructure", Items: []string{"AWS", "Azure", "Google Cloud", "Docker", "Git", "gRPC"}},
    }
}

func GetExperience() []Experience {
    return []Experience{
        {
            Company:     "NATIONAL CAPTIONING INSTITUTE",
            Role:        "Software Engineer",
            Duration:    "FEB 2025 - PRESENT",
            Description: "Developing automated captioning systems that meet strict accessibility and broadcast standards. I also redesigned our testing pipeline and improved frontend performance across key workflows.",
            Tags:        []string{"Python", "Performance", "Automation", "DevOps", "SQL"},
        },
        {
            Company:     "UNCOMMON GIVING",
            Role:        "Software Engineer",
            Duration:    "2023 - PRESENT",
            Description: "Building and maintaining web and mobile applications using JS, Python, and Flutter. I optimized the CI/CD pipelines to run tasks in parallel, reducing build and deployment times by over 50%.",
            Tags:        []string{"Typescript", "Angular", "Python", "Django", "SQL"},
        },
        {
            Company:     "MUSC",
            Role:        "Systems Programmer II",
            Duration:    "2023 - 2025",
            Description: "Led full-stack development for healthcare software using C# and SQL. I helped modernize legacy systems by migrating core applications to .NET Core and improving data workflows.",
            Tags:        []string{"C#", ".Net Core", "Flutter", "SQL", "jQuery"},
        },
        {
            Company:     "District 186",
            Role:        "Computer Programmer and Software Developer",
            Duration:    "2022 - 2023",
            Description: "Built internal tools for student and staff management using PHP and Oracle SQL, streamlining repetitive administrative processes and improving data accuracy.",
            Tags:        []string{"PHP", "Oracle SQL", "Javascript", "Bootstrap"},
        },
    }
}

func GetProjects() []Project {
    return []Project{
        {
            ID:          "proj1",
            Title:       "Lawless Lowcountry Living",
            Description: "A production site hosted on modern cloud infrastructure that demonstrates my full-stack work. Mobile-first design and optimized delivery are core priorities, with attention to accessibility and performance.",
            Tags:        []string{"Cloud Hosting", "Web Dev"},
            Link:        "https://lawlesslowcountryliving.com",
            Icon:        "bi-box-arrow-up-right",
        },
        {
            ID:          "proj4",
            Title:       "The \"OG\" Portfolio",
            Description: "My first portfolio site. A snapshot of where my frontend skills started. I keep it preserved as a reminder of the progression of my design and implementation choices over time. What a ride!",
            Tags:        []string{"Legacy", "HTML/CSS", "Progress"},
            Link:        "https://jmorgan3142001.github.io/portfolio-website/",
            Icon:        "bi-clock-history",
        },
        {
            ID:          "nci_os_1",
            Title:       "Django5 Forms Fieldset",
            Description: "An open-source extension I forked and now maintain at NCI that groups Django form fields semantically to simplify UI structure and improve accessibility for form-heavy applications.",
            Tags:        []string{"Open Source", "Django", "Python", "NCI"},
            Link:        "https://github.com/NCIAdmin/django5-forms-fieldset",
            Icon:        "bi-github",
        },
        {
            ID:          "nci_os_2",
            Title:       "Django5 Scheduler",
            Description: "A task scheduler I forked and contributed to for NCI. It provides simple, ORM-driven scheduling for background jobs and periodic tasks, keeping operations manageable from within Django.",
            Tags:        []string{"Open Source", "Automation", "Django", "NCI"},
            Link:        "https://github.com/NCIAdmin/django5-scheduler",
            Icon:        "bi-calendar-check",
        },
        {
            ID:          "proj2",
            Title:       "Auto-Caption Network",
            Description: "A distributed captioning network built for low latency and high reliability. It manages complex streams to keep captions synchronized and compliant across multiple endpoints.",
            Tags:        []string{"Distributed Systems", "Performance"},
            Link:        "#",
            Icon:        "bi-diagram-3",
        },
        {
            ID:          "proj3",
            Title:       "CRM Pipeline Opt.",
            Description: "Helped to refactor the CI/CD workflows at Uncommon Giving to run staged tasks in parallel, reducing build times by ~50% and enabling faster, more reliable feature deployments.",
            Tags:        []string{"DevOps", "CI/CD"},
            Link:        "#",
            Icon:        "bi-gear-wide-connected",
        },
    }
}

func GetSystemModules() []SystemModule {
    return []SystemModule{
        {
            ID:          "MODULE_01",
            Title:       "Backend Architecture",
            Icon:        "bi-hdd-network",
            Description: "Server design, API construction, and operational practices that keep services secure and highly available.",
            Progress:    100,
            Tags:        []string{"Python (Django)", "Go (Beego)", "PostgreSQL"},
        },
        {
            ID:          "MODULE_02",
            Title:       "Full Stack Integration",
            Icon:        "bi-window-stack",
            Description: "Bridging backend services with user-facing interfaces, focusing on reliable data flows, performance, and usable UI/UX.",
            Progress:    85,
            Tags:        []string{"Angular", "TypeScript", "UI/UX"},
        },
        {
            ID:          "MODULE_03",
            Title:       "Data Management",
            Icon:        "bi-database-fill",
            Description: "Schema design, query optimization, and handling large datasets for robust and maintainable data platforms.",
            Progress:    100,
            Tags:        []string{"PostgreSQL", "SQL Server", "Optimization"},
        },
        {
            ID:          "MODULE_04",
            Title:       "Distributed and Cloud Systems",
            Icon:        "bi-diagram-3",
            Description: "Designing scalable infrastructure, replication, and fault-tolerance strategies for resilient distributed systems.",
            Progress:    75,
            Tags:        []string{"gRPC", "Docker", "AWS"},
        },
    }
}

func GetHardwareProfile() []HardwareInfo {
    return []HardwareInfo{
        {Category: "Computer", Item: "MAC/WINDOWS/LINUX"},
        {Category: "Input", Item: "Custom Ergo Keyboard"},
        {Category: "Pair Programmers", Item: "1x Wife and 2x Staffordshire Terriers"},
    }
}

func GetBooks() []Book {
    return []Book{
        {
            Title:       "Clean Code",
            Author:      "Robert C. Martin",
            Description: "The practical guide I return to for discipline and patterns that make code readable, testable, and maintainable.",
            Type:        "CORE_LOGIC",
            Link:        "https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882",
        },
        {
            Title:       "Designing Data-Intensive Applications",
            Author:      "Martin Kleppmann",
            Description: "A foundational book on the principles and tradeoffs of building scalable, reliable data systems which is essential for designing real-world distributed applications.",
            Type:        "DATABASE_SYS",
            Link:        "https://www.amazon.com/Designing-Data-Intensive-Applications-Reliable-Maintainable/dp/1449373321",
        },
        {
            Title:       "Modern Operating Systems",
            Author:      "Andrew S. Tanenbaum",
            Description: "A comprehensive textbook covering OS concepts and design (that I used in my undergrad os course). Useful when reasoning about scheduling, memory, and system interactions.",
            Type:        "KERNEL_OPS",
            Link:        "https://www.amazon.com/Modern-Operating-Systems-Andrew-Tanenbaum/dp/013359162X",
        },
        {
            Title:       "Operating Systems: Three Easy Pieces",
            Author:      "Remzi & Andrea Arpaci-Dusseau",
            Description: "A clear, accessible exploration of virtualization, concurrency, and file systems that shaped my understanding even further of low-level systems behavior.",
            Type:        "KERNEL_OPS",
            Link:        "https://www.amazon.com/Operating-Systems-Three-Easy-Pieces/dp/198508659X",
        },
        {
            Title:       "Database Internals",
            Author:      "Alex Petrov",
            Description: "An in-depth look at how databases function under the hood, from storage engines to replication and distributed consensus.",
            Type:        "DATABASE_SYS",
            Link:        "https://www.amazon.com/Database-Internals-Deep-Distributed-Systems/dp/1492040347",
        },
    }
}

func GetNextReads() []Book {
    return []Book{
        {
            Title:       "TCP/IP Illustrated, Vol. 1",
            Author:      "W. Richard Stevens",
            Description: "A visual and protocol-level guide to the core networking stacks that power the internet. This is a great reference for systems and network debugging. Excited to dive deep into network computing!",
            Type:        "NETWORKING",
            Link:        "https://www.amazon.com/TCP-Illustrated-Vol-Protocols-Addison-Wesley/dp/0201633469",
        },
    }
}

func GetDigitalResources() []ResourceItem {
    return []ResourceItem{
        {
            Title:       "The Go Blog",
            Description: "Official news and insights from the Go team.",
            Link:        "https://go.dev/blog/",
            Icon:        "bi-google",
        },
        {
            Title:       "High Scalability",
            Description: "Building bigger, faster, more reliable websites.",
            Link:        "https://highscalability.com/",
            Icon:        "bi-graph-up-arrow",
        },
        {
            Title:       "Django Documentation",
            Description: "The Model layer and ORM deep dives.",
            Link:        "https://docs.djangoproject.com/en/stable/",
            Icon:        "bi-filetype-py",
        },
    }
}

func GetResearchPapers() []ResearchPaper {
    return []ResearchPaper{
        {
            Title: "MapReduce: Simplified Data Processing",
            Topic: "Distributed Systems",
            Note:  "How Google processes massive datasets on commodity hardware with a fault-tolerant, distributed model; a cornerstone for batch processing systems.",
            Link:  "https://research.google.com/archive/mapreduce-osdi04.pdf",
        },
        {
            Title: "Time, Clocks, and Ordering",
            Topic: "Concurrency",
            Note:  "Lamport's classic on ordering events in distributed systems and the rationale for logical clocks.",
            Link:  "https://lamport.azurewebsites.net/pubs/time-clocks.pdf",
        },
        {
            Title: "The Google File System",
            Topic: "Storage Systems",
            Note:  "Architecture and design decisions behind a scalable file system used for large data processing workloads.",
            Link:  "https://research.google.com/archive/gfs-sosp2003.pdf",
        },
        {
            Title: "Dynamo: Amazon's Highly Available Key-Value Store",
            Topic: "Distributed Systems",
            Note:  "A foundational exploration of eventual consistency, gossip protocols, and partition-tolerant system design.",
            Link:  "https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf",
        },
        {
            Title: "In Search of an Understandable Consensus Algorithm (Raft)",
            Topic: "Consensus",
            Note:  "A consensus algorithm designed for clarity and practical implementation which is often used as an approachable alternative to Paxos.",
            Link:  "https://raft.github.io/raft.pdf",
        },
    }
}

func GetAlgorithms() []Algorithm {
    return []Algorithm{
        {
            Name:        "Circular Stream Buffer",
            Lang:        "TypeScript",
            Description: "A fixed-size sliding window implementation designed to handle high-velocity caption ingestion while maintaining constant memory usage in the browser.",
            Link:        "https://github.com/xtermjs/xterm.js",
        },
        {
            Name:        "Barrier Synchronization",
            Lang:        "C++ / OpenMP",
            Description: "Implementing thread barriers to study synchronization primitives and expose race conditions in concurrent programs.",
            Link:        "https://github.com/OpenMP/examples",
        },
        {
            Name:        "gRPC Store",
            Lang:        "C++",
            Description: "A distributed key-value store built with thread pools and replication primitives to explore consistency and performance tradeoffs.",
            Link:        "https://github.com/grpc/grpc/tree/master/examples/cpp",
        },
        {
            Name:        "MapReduce Coordinator",
            Lang:        "Go",
            Description: "A fault-tolerant master-worker implementation that manages task distribution, handles worker failure detection, and aggregates intermediate data reduction.",
            Link:        "https://github.com/google/go-cloud",
        },
        {
            Name:        "Persistent Priority Queue",
            Lang:        "Python",
            Description: "An ordering algorithm for the Django Scheduler that manages task execution based on priority weights and time-windows to prevent job starvation.",
            Link:        "https://github.com/NCIAdmin/django5-scheduler",
        },
    }
}

func GetResearchMeta() []ResearchMeta {
    return []ResearchMeta{
        {Label: "SYSTEM", Value: "ACTIVE", Status: "success"},
        {Label: "THREADS", Value: "16", Status: "warning"},
        {Label: "FOCUS", Value: "DISTRIBUTED", Status: "accent"},
    }
}