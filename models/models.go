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

type Challenge struct {
    Id          int       `orm:"auto"`
    Title       string    `orm:"size(255)"`
    Description string    `orm:"type(text)"`
    Difficulty  string    `orm:"size(50)"` 
    Category    string    `orm:"size(50)"` 
    Type        string    `orm:"size(50)"` 
    Language    string    `orm:"size(50);null"` 
    StarterCode string    `orm:"type(text);null"` 
    TestCases   []*TestCase `orm:"reverse(many)"`
}

func (u *Challenge) TableName() string {
    return "challenge"
}

type TestCase struct {
    Id             int        `orm:"auto"`
    Challenge      *Challenge `orm:"rel(fk);on_delete(cascade)"`
    InputArgs      string     `orm:"type(text)"`
    ExpectedOutput string     `orm:"type(text)"`
}

func (u *TestCase) TableName() string {
    return "test_case"
}

func SeedChallenges() {
    o := orm.NewOrm()
    
    // Check if data exists
    cnt, _ := o.QueryTable("challenge").Count()
    if cnt > 0 {
        fmt.Println("Skipping seed: Challenges already exist.")
        return
    }

    fmt.Println("Seeding Real-World Portfolio Challenges...")

    // =========================================================
    // 1. EASY (BUG FIX) - NCI: Legacy Timestamp Logic
    // =========================================================
    c1 := Challenge{
        Title:       "Legacy Timestamp Bug",
        Description: "At **National Captioning Institute**, we process legacy CEA-608 caption files. I discovered a bug in our ingestion pipeline where timestamps formatted as `MM:SS:mmm` were calculating incorrectly for durations over 1 minute.\n\nThe math logic below has an error - it's dividing by 60 instead of 60000 to get minutes. **Fix the math logic** so `65000ms` becomes `'01:05:000'`.",
        Difficulty:  "Easy",
        Category:    "NCI / Debugging",
        Type:        "CODE",
        Language:    "python",
        StarterCode: `def solve(ms):
    """
    Converts milliseconds to formatted string MM:SS:mmm
    Args:
        ms (int): Total milliseconds (e.g., 61500)
    Returns:
        str: Formatted string (e.g., "01:05:000")
    """
    # BUG: I found this logic in a legacy module. 
    # It fails when ms > 60000.
    minutes = ms // 60
    seconds = (ms // 1000) % 60
    rem_ms = ms % 1000
    
    return f"{minutes:02}:{seconds:02}:{rem_ms:03}"`,
    }
    id1, _ := o.Insert(&c1)

    t1_1 := TestCase{Challenge: &Challenge{Id: int(id1)}, InputArgs: "65000", ExpectedOutput: "'01:05:000'"}
    t1_2 := TestCase{Challenge: &Challenge{Id: int(id1)}, InputArgs: "125500", ExpectedOutput: "'02:05:500'"}
    o.InsertMulti(2, []TestCase{t1_1, t1_2})

    // =========================================================
    // 2. EASY (VALIDATION) - UG: Donation Input Sanitizer
    // =========================================================
    c2 := Challenge{
        Title:       "Payment Input Validator",
        Description: "At **Uncommon Giving**, ensuring data integrity before hitting the payment gateway is critical. We noticed users could bypass frontend checks and submit negative amounts or amounts with excessive decimal precision (e.g., `$10.555`), which caused downstream API failures.\n\nWrite a function that validates if a string is a valid positive money amount (max 2 decimal places). Return `True` or `False`.",
        Difficulty:  "Easy",
        Category:    "UG / Validation",
        Type:        "CODE",
        Language:    "python",
        StarterCode: `def solve(amount_str):
    """
    Validates if a string is a valid currency amount.
    Rules:
    1. Must be a positive number.
    2. No more than 2 decimal places.
    3. No currency symbols ($).
    
    Args: amount_str (str)
    Returns: bool
    """
    # TODO: Validate the string format.
    
    return False`,
    }
    id2, _ := o.Insert(&c2)

    t2_1 := TestCase{Challenge: &Challenge{Id: int(id2)}, InputArgs: "'10.50'", ExpectedOutput: "True"}
    t2_2 := TestCase{Challenge: &Challenge{Id: int(id2)}, InputArgs: "'-5.00'", ExpectedOutput: "False"}
    t2_3 := TestCase{Challenge: &Challenge{Id: int(id2)}, InputArgs: "'10.555'", ExpectedOutput: "False"}
    o.InsertMulti(3, []TestCase{t2_1, t2_2, t2_3})

    // =========================================================
    // 3. MEDIUM (BUG FIX) - NCI: Compliance Splitter
    // =========================================================
    c3 := Challenge{
        Title:       "FCC Compliance Splitter",
        Description: "Strict FCC regulations require broadcast captions to never exceed 32 characters per line. The existing library used at **NCI** would sometimes chop words in half to force a fit, violating accessibility standards.\n\nI rewrote the logic to ensure we only split on spaces. **Fix the bug** in the loop below to ensure no line exceeds 32 chars and no words are cut.",
        Difficulty:  "Medium",
        Category:    "NCI / Strings",
        Type:        "CODE",
        Language:    "python",
        StarterCode: `def solve(text):
    """
    Splits text into lines of max 32 chars without cutting words.
    Args: text (str)
    Returns: List[str]
    """
    words = text.split()
    lines = []
    current_line = ""

    for word in words:
        # BUG: This logic crashes or produces bad output on edge cases.
        # It fails to account for the added space character length.
        if len(current_line) + len(word) > 32:
            lines.append(current_line)
            current_line = word
        else:
            current_line += word
            
    lines.append(current_line)
    return lines`,
    }
    id3, _ := o.Insert(&c3)

    t3_1 := TestCase{Challenge: &Challenge{Id: int(id3)}, InputArgs: "'This is a test of the emergency broadcast system'", ExpectedOutput: "['This is a test of the emergency', 'broadcast system']"} 
    t3_2 := TestCase{Challenge: &Challenge{Id: int(id3)}, InputArgs: "'A short line'", ExpectedOutput: "['A short line']"}
    o.InsertMulti(2, []TestCase{t3_1, t3_2})

    // =========================================================
    // 4. MEDIUM (ALGO) - UG: The "Lost Penny" Problem
    // =========================================================
    c4 := Challenge{
        Title:       "The 'Lost Penny' Problem",
        Description: "In the **Uncommon Giving** workplace platform, employees can donate a bundle (e.g., $100) split across 3 charities. Since `$100 / 3 = $33.333...`, simply dividing results in a 'lost penny' ($99.99 total).\n\nI implemented a distribution algorithm to handle this. You need to split `total_cents` among `n` recipients. The remainder pennies must be distributed to the first `r` recipients so the sum matches exactly.",
        Difficulty:  "Medium",
        Category:    "UG / FinTech",
        Type:        "CODE",
        Language:    "python",
        StarterCode: `def solve(total_cents, n_charities):
    """
    Splits total_cents among n_charities.
    Args:
        total_cents (int): Total amount (e.g. 10000)
        n_charities (int): Number of recipients (e.g. 3)
    Returns:
        List[int]: List of amounts that sum exactly to total_cents.
                   (e.g. [3334, 3333, 3333])
    """
    # TODO: Calculate base amount and remainder.
    # Distribute remainder to the first few items.
    
    return []`,
    }
    id4, _ := o.Insert(&c4)

    t4_1 := TestCase{Challenge: &Challenge{Id: int(id4)}, InputArgs: "10000, 3", ExpectedOutput: "[3334, 3333, 3333]"}
    t4_2 := TestCase{Challenge: &Challenge{Id: int(id4)}, InputArgs: "100, 6", ExpectedOutput: "[17, 17, 17, 17, 16, 16]"}
    o.InsertMulti(2, []TestCase{t4_1, t4_2})

    // =========================================================
    // 5. HARD (GRAPH) - Scheduler: Cascading Cancellation
    // =========================================================
    c5 := Challenge{
        Title:       "Job Dependency Cascade",
        Description: "I contribute to the open-source **Django Scheduler** used at NCI. We had an issue where if a parent job failed, dependent child jobs would sit in a 'Pending' state forever.\n\nI implemented a graph traversal to identify downstream impacts. Given a dependency graph `{'A': ['B']}` (B depends on A) and a failed job, return a **sorted list** of all jobs that must be cancelled.",
        Difficulty:  "Hard",
        Category:    "Systems / Graph",
        Type:        "CODE",
        Language:    "python",
        StarterCode: `def solve(deps, failed_job):
    """
    Finds all downstream jobs affected by a failure.
    Args:
        deps (Dict[str, List[str]]): Key is a job, value is list of jobs that depend on it.
                                     e.g. {'A': ['B'], 'B': ['C']} means A -> B -> C
        failed_job (str): The ID of the job that crashed.
    Returns:
        List[str]: Sorted list of all impacted jobs (including failed_job).
    """
    # TODO: Perform a BFS or DFS starting at failed_job
    # to find everything downstream.
    
    impacted = []
    return impacted`,
    }
    id5, _ := o.Insert(&c5)

    t5_1 := TestCase{Challenge: &Challenge{Id: int(id5)}, InputArgs: "{'A': ['B'], 'B': ['C'], 'C': []}, 'A'", ExpectedOutput: "['A', 'B', 'C']"}
    t5_2 := TestCase{Challenge: &Challenge{Id: int(id5)}, InputArgs: "{'A': ['B', 'C'], 'B': [], 'C': ['D'], 'D': []}, 'A'", ExpectedOutput: "['A', 'B', 'C', 'D']"}
    o.InsertMulti(2, []TestCase{t5_1, t5_2})

    // =========================================================
    // 6. HARD (DATA) - UG: Distributed Ledger Fix
    // =========================================================
    c6 := Challenge{
        Title:       "Async Ledger Reconciliation",
        Description: "At **Uncommon Giving**, we rely on webhooks from payment processors. Due to distributed system lag, events often arrive out of order (e.g., a 'REVERSE' event arriving before the 'DEPOSIT' it reverses).\n\nTo ensure our ledger is accurate, we need to calculate the final balance regardless of event order. `REVERSE` events contain a `ref_id` pointing to the transaction they cancel out.",
        Difficulty:  "Hard",
        Category:    "UG / Distributed",
        Type:        "CODE",
        Language:    "python",
        StarterCode: `def solve(events):
    """
    Calculates final balance from out-of-order stream.
    Args:
        events (List[Dict]): List of dicts. 
            e.g. {"id": 1, "type": "DEPOSIT", "amount": 100}
                 {"id": 2, "type": "REVERSE", "ref_id": 1}
    Returns:
        int: Final Balance
    """
    # TODO: Process the events. 
    # Hint: Maybe two passes? Or use a map to track valid transaction IDs?
    
    return 0`,
    }
    id6, _ := o.Insert(&c6)

    t6_1 := TestCase{Challenge: &Challenge{Id: int(id6)}, InputArgs: "[{'id':1, 'type':'DEPOSIT', 'amount':100}, {'id':2, 'type':'WITHDRAW', 'amount':50}]", ExpectedOutput: "50"}
    t6_2 := TestCase{Challenge: &Challenge{Id: int(id6)}, InputArgs: "[{'id':2, 'type':'REVERSE', 'ref_id':1}, {'id':1, 'type':'DEPOSIT', 'amount':100}]", ExpectedOutput: "0"}
    o.InsertMulti(2, []TestCase{t6_1, t6_2})

    fmt.Println("Seeding complete.")
}

func init() {
	orm.RegisterModel(new(AccessLog), new(Challenge), new(TestCase))

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

	SeedChallenges()
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
	Company     string
	Role        string
	Duration    string
	Description string
	Tags        []string
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
		// Personal Projects
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

		// NCI Open Source Contributions
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

		// Internal Achievements
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

func GetChallenges() []Challenge {
    o := orm.NewOrm()
    var challenges []Challenge
    o.QueryTable("challenge").All(&challenges)
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
    o.QueryTable("test_cases").Filter("challenge_id", challengeId).All(&cases)
    return cases
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