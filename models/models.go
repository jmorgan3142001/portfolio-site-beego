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

// --- Database Models ---

type AccessLog struct {
    Id        int       `orm:"auto"`
    Name      string    `orm:"size(50)"`
    Message   string    `orm:"type(text)"`
    Signature string    `orm:"size(10)"` 
    Terminal  string    `orm:"size(30)"` 
    ProcessID int       
    Created   time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
    orm.RegisterModel(new(AccessLog))

    orm.RegisterDriver("postgres", orm.DRPostgres)

    dbUrl := os.Getenv("DATABASE_URL")

    err := orm.RegisterDataBase("default", "postgres", dbUrl)
    if err != nil {
        panic(fmt.Errorf("failed to register database: %v", err))
    }

    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        fmt.Println("Database Sync Error:", err)
    }
}

// --- Logic for Access Logs ---

func GetAccessLogs() []AccessLog {
    o := orm.NewOrm()
    var logs []AccessLog
    o.QueryTable("access_log").OrderBy("-Created").All(&logs)
    return logs
}

func AddAccessLog(name, message, userAgent string) error {
    if CheckToxicity(name) || CheckToxicity(message) {
        return errors.New("input rejected: content detected as toxic or aggressive")
    }

    o := orm.NewOrm()
    
    // Generate Signature (Pseudo-Git-SHA)
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

// --- Static Data Structs ---

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
    Company  string
    Role     string
    Duration string
    Desc     string
    Tags     []string
}

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
    Title  string
    Author string
    Desc   string
    Type   string
    Link   string
}

type ResearchPaper struct {
    Title string
    Topic string
    Note  string
    Link  string
}

type Algorithm struct {
    Name string
    Lang string
    Desc string
    Link string
}

type ResearchMeta struct {
    Label  string
    Value  string
    Status string
}

type ComplexityStat struct {
    Operation  string
    Complexity string
}

type ResourceItem struct {
    Title string
    Desc  string
    Link  string
    Icon  string
}

// --- Data Retrieval Functions ---

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
            Company:  "NATIONAL CAPTIONING INSTITUTE",
            Role:     "Software Engineer",
            Duration: "FEB 2025 - PRESENT",
            Desc:     "Building automated captioning tools to meet strict accessibility standards. I also overhauled the testing infrastructure and improved frontend performance.",
            Tags:     []string{"Python", "Performance", "Automation", "DevOps", "SQL"},
        },
        {
            Company:  "UNCOMMON GIVING",
            Role:     "Software Engineer",
            Duration: "2023 - PRESENT",
            Desc:     "Building web and mobile apps with JS, Python, and Flutter. I also tuned the CI/CD pipelines, cutting build times in half.",
            Tags:     []string{"Typescript", "Angular", "Python", "Django", "SQL"},
        },
        {
            Company:  "MUSC",
            Role:     "Systems Programmer II",
            Duration: "2023 - 2025",
            Desc:     "Led full-stack development for healthcare apps using C# and SQL. I helped modernize legacy systems by moving them to .NET Core.",
            Tags:     []string{"C#", ".Net Core", "Flutter", "SQL", "jQuery"},
        },
        {
            Company:  "District 186",
            Role:     "Computer Programmer and Software Developer",
            Duration: "2022 - 2023",
            Desc:     "Built student and staff management tools using PHP and Oracle SQL to make administrative tasks easier.",
            Tags:     []string{"PHP", "Oracle SQL", "Javascript", "Bootstrap"},
        },
    }
}

func GetProjects() []Project {
    return []Project{
        // Personal Projects
        {
            ID:          "proj1",
            Title:       "Lawless Lowcountry Living",
            Description: "A live production site running on modern cloud infrastructure. It showcases my full-stack work, focusing on mobile-first design and fast content delivery.",
            Tags:        []string{"Cloud Hosting", "Web Dev"},
            Link:        "http://lawlesslowcountryliving.com",
            Icon:        "bi-box-arrow-up-right",
        },
        {
            ID:          "proj4",
            Title:       "The \"OG\" Portfolio",
            Description: "A blast from the past. This was my first portfolio site. It's a bit rough compared to my current work, but I keep it up to show how my frontend skills have evolved.",
            Tags:        []string{"Legacy", "HTML/CSS", "Progress"},
            Link:        "https://jmorgan3142001.github.io/portfolio-website/",
            Icon:        "bi-clock-history",
        },

        // NCI Open Source Contributions
        {
            ID:          "nci_os_1",
            Title:       "Django5 Forms Fieldset",
            Description: "An open-source tool I forked and maintain at NCI. It extends Django 5 to let developers group form fields semantically, making UIs cleaner and more accessible.",
            Tags:        []string{"Open Source", "Django", "Python", "NCI"},
            Link:        "https://github.com/NCIAdmin/django5-forms-fieldset",
            Icon:        "bi-github",
        },
        {
            ID:          "nci_os_2",
            Title:       "Django5 Scheduler",
            Description: "A task scheduler I forked and contributed to NCI's open source repo. It helps manage background jobs and periodic tasks directly through the Django ORM.",
            Tags:        []string{"Open Source", "Automation", "Django", "NCI"},
            Link:        "https://github.com/NCIAdmin/django5-scheduler",
            Icon:        "bi-calendar-check",
        },

        // Internal Achievements
        {
            ID:          "proj2",
            Title:       "Auto-Caption Network",
            Description: "A distributed network built for NCI to deliver low-latency captions. It handles complex data streams to ensure captions stay in sync and meet compliance rules.",
            Tags:        []string{"Distributed Systems", "Performance"},
            Link:        "#",
            Icon:        "bi-diagram-3",
        },
        {
            ID:          "proj3",
            Title:       "CRM Pipeline Opt.",
            Description: "I reworked the CI/CD workflows at Uncommon Giving to run tasks in parallel. This cut build times by 50% and helped us deploy features much faster.",
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
            Description: "Server logic, API design, and keeping things online.",
            Progress:    100,
            Tags:        []string{"Python (Django)", "Go (Beego)", "PostgreSQL"},
        },
        {
            ID:          "MODULE_02",
            Title:       "Distributed Systems",
            Icon:        "bi-diagram-3",
            Description: "Infrastructure scaling and consensus algorithms.",
            Progress:    75,
            Tags:        []string{"gRPC", "Docker", "AWS"},
        },
        {
            ID:          "MODULE_03",
            Title:       "Full Stack Integration",
            Icon:        "bi-window-stack",
            Description: "Connecting backend logic to user-friendly interfaces.",
            Progress:    85,
            Tags:        []string{"Angular", "TypeScript", "UI/UX"},
        },
        {
            ID:          "MODULE_04",
            Title:       "Data Management",
            Icon:        "bi-database-fill",
            Description: "Designing schemas and managing large datasets.",
            Progress:    100,
            Tags:        []string{"PostgreSQL", "SQL Server", "Optimization"},
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
            Title:  "Clean Code",
            Author: "Robert C. Martin",
            Desc:   "The go-to guide for writing code that is easy to read and maintain.",
            Type:   "CORE_LOGIC",
            Link:   "https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882",
        },
        {
            Title:  "Designing Data-Intensive Applications",
            Author: "Martin Kleppmann",
            Desc:   "A must-read for understanding how distributed systems actually work.",
            Type:   "DATABASE_SYS",
            Link:   "https://www.amazon.com/Designing-Data-Intensive-Applications-Reliable-Maintainable/dp/1449373321",
        },
        {
            Title:  "Operating Systems: Three Easy Pieces",
            Author: "Remzi & Andrea Arpaci-Dusseau",
            Desc:   "A great look at virtualization, concurrency, and file systems.",
            Type:   "KERNEL_OPS",
            Link:   "https://www.amazon.com/Operating-Systems-Three-Easy-Pieces/dp/198508659X",
        },
        {
            Title:  "Modern Operating Systems",
            Author: "Andrew S. Tanenbaum",
            Desc:   "A classic textbook on how computer operating systems function.",
            Type:   "KERNEL_OPS",
            Link:   "https://www.amazon.com/Modern-Operating-Systems-Andrew-Tanenbaum/dp/013359162X",
        },
        {
            Title:  "Database Internals",
            Author: "Alex Petrov",
            Desc:   "Explains how databases work under the hood, from storage engines to distributed consensus.",
            Type:   "DATABASE_SYS",
            Link:   "https://www.amazon.com/Database-Internals-Deep-Distributed-Systems/dp/1492040347",
        },
    }
}

func GetNextReads() []Book {
    return []Book{
        {
            Title:  "TCP/IP Illustrated, Vol. 1",
            Author: "W. Richard Stevens",
            Desc:   "A visual guide to the protocols that run the internet.",
            Type:   "NETWORKING",
            Link:   "https://www.amazon.com/TCP-Illustrated-Vol-Protocols-Addison-Wesley/dp/0201633469",
        },
    }
}

func GetDigitalResources() []ResourceItem {
    return []ResourceItem{
        {
            Title: "The Go Blog",
            Desc: "Official news and insights from the Go team.",
            Link: "https://go.dev/blog/",
            Icon: "bi-google",
        },
        {
            Title: "High Scalability",
            Desc: "Building bigger, faster, more reliable websites.",
            Link: "http://highscalability.com/",
            Icon: "bi-graph-up-arrow",
        },
        {
            Title: "Django Documentation",
            Desc: "The Model layer and ORM deep dives.",
            Link: "https://docs.djangoproject.com/en/stable/",
            Icon: "bi-filetype-py",
        },
    }
}

func GetResearchPapers() []ResearchPaper {
    return []ResearchPaper{
        {
            Title: "MapReduce: Simplified Data Processing",
            Topic: "Distributed Systems",
            Note:  "How Google processes massive datasets on commodity hardware with built-in fault tolerance.",
            Link:  "https://research.google.com/archive/mapreduce-osdi04.pdf",
        },
        {
            Title: "Time, Clocks, and Ordering",
            Topic: "Concurrency",
            Note:  "Lamport's classic paper on ordering events in distributed systems using logical clocks.",
            Link:  "https://lamport.azurewebsites.net/pubs/time-clocks.pdf",
        },
        {
            Title: "The Google File System",
            Topic: "Storage Systems",
            Note:  "The architecture behind Google's scalable file system for data-intensive apps.",
            Link:  "https://research.google.com/archive/gfs-sosp2003.pdf",
        },
        {
            Title: "In Search of an Understandable Consensus Algorithm (Raft)",
            Topic: "Consensus",
            Note:  "A consensus algorithm built to be easier to understand than Paxos. It's the backbone of many modern distributed systems.",
            Link:  "https://raft.github.io/raft.pdf",
        },
    }
}

func GetAlgorithms() []Algorithm {
    return []Algorithm{
        {
            Name: "Barrier Synchronization",
            Lang: "C++ / OpenMP",
            Desc: "Building thread barriers from scratch to better understand race conditions.",
            Link: "https://github.com/OpenMP/examples", 
        },
        {
            Name: "gRPC Store",
            Lang: "C++",
            Desc: "A distributed key-value store I built using thread pools and custom replication.",
            Link: "https://github.com/grpc/grpc/tree/master/examples/cpp", 
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

func GetComplexityStats() []ComplexityStat {
    return []ComplexityStat{
        {Operation: "Hash Map Access", Complexity: "O(1)"},
        {Operation: "Bin. Tree Search", Complexity: "O(log n)"},
        {Operation: "Quick Sort", Complexity: "O(n log n)"},
    }
}