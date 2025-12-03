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
			Tags:     []string{"C#", ".Net Core", "SQL", "Flutter"},
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