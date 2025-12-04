package models

// --- Data Structures ---

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
    PublicRepos     int    `json:"public_repos"`
    Login           string `json:"login"`
    LatestCommitMsg string `json:"latest_commit_msg"`
    LatestCommitUrl string `json:"latest_commit_url"` 
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
}

type ResearchPaper struct {
    Title string
    Topic string
    Note  string
}

type Algorithm struct {
    Name string
    Lang string
    Desc string
}

type ResearchMeta struct {
    Label  string
    Value  string
    Status string // "success", "warning", or "accent" for color logic
}

type ComplexityStat struct {
    Operation  string
    Complexity string
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
			Desc:     "Engineered high-precision automated captioning solutions exceeding accessibility compliance standards. Re-architected testing infrastructure and optimized frontend performance.",
			Tags:     []string{"Python", "Performance", "Automation", "DevOps", "SQL"},
		},
		{
			Company:  "UNCOMMON GIVING",
			Role:     "Software Engineer",
			Duration: "2023 - PRESENT",
			Desc:     "Architect of scalable web and mobile solutions using JS, Python (Django), and Flutter. Optimized CI/CD pipelines reducing build times by 50%.",
			Tags:     []string{"Typescript", "Angular", "Python", "Django", "SQL"},
		},
		{
			Company:  "MUSC",
			Role:     "Systems Programmer II",
			Duration: "2023 - 2025",
			Desc:     "Spearheaded full-stack development of enterprise healthcare applications using C# and SQL. Modernized legacy applications to .NET Core.",
			Tags:     []string{"C#", ".Net Core", "Flutter", "SQL"},
		},
	}
}

func GetProjects() []Project {
    return []Project{
        {
            ID:          "proj1",
            Title:       "Lawless Lowcountry Living",
            Description: "A comprehensive production website architected on modern cloud infrastructure. This project demonstrates end-to-end full-stack capabilities, featuring a responsive mobile-first design, optimized content delivery networks (CDN) for rapid asset loading, and a scalable backend system.",
            Tags:        []string{"Cloud Hosting", "Web Dev"},
            Link:        "http://lawlesslowcountryliving.com",
            Icon:        "bi-box-arrow-up-right",
        },
        {
            ID:          "proj2",
            Title:       "Auto-Caption Network",
            Description: "A distributed node network designed to provide ultra-low-latency captions in partnership with NCI. The system implements high-precision automation logic to ensure strict accessibility compliance, handling complex data streams with high availability and minimal delay.",
            Tags:        []string{"Distributed Systems", "Performance"},
            Link:        "#",
            Icon:        "bi-diagram-3",
        },
        {
            ID:          "proj3",
            Title:       "CRM Pipeline Opt.",
            Description: "Re-engineered CI/CD workflows for Uncommon Giving by implementing advanced parallel execution strategies. This optimization initiative streamlined the development lifecycle, successfully reducing build times by 50% and significantly increasing deployment velocity.",
            Tags:        []string{"DevOps", "CI/CD"},
            Link:        "#",
            Icon:        "bi-gear-wide-connected",
        },
        {
            ID:          "proj4",
            Title:       "The \"OG\" Portfolio",
            Description: "A fun blast from the past! This was my very first attempt at a portfolio site. While it's a little rough around the edges compared to my current work, I keep it online as a benchmark to highlight just how far my frontend skills and design sensibilities have evolved.",
            Tags:        []string{"Legacy", "HTML/CSS", "Progress"},
            Link:        "https://jmorgan3142001.github.io/portfolio-website/",
            Icon:        "bi-clock-history",
        },
    }
}

func GetSystemModules() []SystemModule {
    return []SystemModule{
        {
            ID:          "MODULE_01",
            Title:       "Backend Architecture",
            Icon:        "bi-hdd-network",
            Description: "High-availability server logic and API design.",
            Progress:    90,
            Tags:        []string{"Python (Django)", "Go (Beego)", "PostgreSQL"},
        },
        {
            ID:          "MODULE_02",
            Title:       "Distributed Systems",
            Icon:        "bi-diagram-3",
            Description: "Scalable infrastructure and consensus logic.",
            Progress:    45,
            Tags:        []string{"gRPC", "Docker", "AWS"},
        },
        {
            ID:          "MODULE_03",
            Title:       "Full Stack Integration",
            Icon:        "bi-window-stack",
            Description: "Bridging complex backend logic with user interfaces.",
            Progress:    67,
            Tags:        []string{"Angular", "TypeScript", "UI/UX"},
        },
		{
            ID:          "MODULE_04",
            Title:       "Data Management",
            Icon:        "bi-database-fill",
            Description: "Architect of efficient schemas and handling large-scale datasets.",
            Progress:    100,
            Tags:        []string{"PostgreSQL", "SQL Server", "SQLite"},
        },
    }
}

func GetHardwareProfile() []HardwareInfo {
    return []HardwareInfo{
        {Category: "Computer", Item: "MAC/WINDOWS/LINUX"},
        {Category: "Monitor", Item: "INNOCN 24.5\" 240Hz"},
        {Category: "Input", Item: "Custom Ergo Keyboard"},
        {Category: "Companions", Item: "1x Wife and 2x Staffordshire Terriers"},
    }
}

func GetBooks() []Book {
    return []Book{
        {Title: "Clean Code", Author: "Robert C. Martin", Desc: "The standard protocol for writing maintainable software.", Type: "CORE_LOGIC"},
        {Title: "Designing Data-Intensive Applications", Author: "Martin Kleppmann", Desc: "Essential for understanding distributed systems and scalability.", Type: "DATABASE_SYS"},
        {Title: "Operating Systems: Three Easy Pieces", Author: "Remzi & Andrea Arpaci-Dusseau", Desc: "Deep dive into virtualization, concurrency, and persistence.", Type: "KERNEL_OPS"},
        {Title: "Modern Operating Systems", Author: "Andrew S. Tanenbaum", Desc: "The definitive guide to underlying computer architecture.", Type: "KERNEL_OPS"},
    }
}

func GetResearchPapers() []ResearchPaper {
    return []ResearchPaper{
        {Title: "MapReduce: Simplified Data Processing", Topic: "Distributed Systems", Note: "Analysis of Google's implementation of map and reduce primitives for large clusters. Key focus on fault tolerance."},
        {Title: "Time, Clocks, and Ordering", Topic: "Concurrency", Note: "Leslie Lamport's seminal work on partial ordering and logical clocks in distributed systems."},
    }
}

func GetAlgorithms() []Algorithm {
    return []Algorithm{
        {Name: "Barrier Synchronization", Lang: "C++ / OpenMP", Desc: "Implementing thread barriers without standard libraries to understand race conditions."},
        {Name: "gRPC Store", Lang: "C++", Desc: "A distributed key-value store utilizing thread pools and custom replication logic."},
    }
}

func GetResearchMeta() []ResearchMeta {
    return []ResearchMeta{
        {Label: "SYSTEM", Value: "ACTIVE", Status: "success"},
        {Label: "THREADS", Value: "8", Status: "warning"},
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