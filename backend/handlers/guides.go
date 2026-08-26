package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"prep-dashboard/backend/models"
)

// topicGuides is the curated set of study materials for cracking good companies.
var topicGuides = []models.TopicGuide{
	{
		Key:     "dsa",
		Name:    "DSA",
		Icon:    "🧮",
		Summary: "Data Structures & Algorithms — the core of every coding round. Aim for pattern recognition, not memorization.",
		Roadmap: []string{
			"Arrays, Strings, Hashing, Two Pointers, Sliding Window",
			"Stacks, Queues, Linked Lists, Recursion & Backtracking",
			"Binary Search, Sorting, Greedy",
			"Trees, BST, Tries, Heaps",
			"Graphs (BFS/DFS, Topo sort, Dijkstra, Union-Find)",
			"Dynamic Programming (1D, 2D, knapsack, LIS, intervals)",
		},
		Resources: []models.Resource{
			{Type: "sheet", Title: "Striver's SDE Sheet (A2Z)", URL: "https://takeuforward.org/strivers-a2z-dsa-course/strivers-a2z-dsa-course-sheet-2/", Note: "Most complete free roadmap"},
			{Type: "sheet", Title: "NeetCode 150", URL: "https://neetcode.io/practice", Note: "Pattern-wise, great for FAANG"},
			{Type: "playlist", Title: "Striver DSA Playlist", URL: "https://www.youtube.com/@takeUforward", Note: "Video explanations"},
			{Type: "practice", Title: "LeetCode Top Interview 150", URL: "https://leetcode.com/studyplan/top-interview-150/", Note: "Curated must-do problems"},
			{Type: "practice", Title: "Codeforces", URL: "https://codeforces.com/", Note: "For strong problem-solving"},
		},
	},
	{
		Key:     "lld",
		Name:    "LLD",
		Icon:    "🧱",
		Summary: "Low Level Design — OOP, design patterns, and building extensible class structures. Common in SDE-2 and product companies.",
		Roadmap: []string{
			"OOP fundamentals + SOLID principles",
			"Design patterns (Strategy, Factory, Observer, Singleton, Decorator)",
			"UML & class relationships",
			"Practice: Parking Lot, BookMyShow, Splitwise, Elevator, Rate Limiter",
			"Concurrency basics in design",
		},
		Resources: []models.Resource{
			{Type: "repo", Title: "Awesome LLD (GitHub)", URL: "https://github.com/ashishps1/awesome-low-level-design", Note: "Problems + solutions"},
			{Type: "course", Title: "Grokking the OOD Interview", URL: "https://www.designgurus.io/course/grokking-the-object-oriented-design-interview", Note: "Structured LLD prep"},
			{Type: "playlist", Title: "Concept && Coding - LLD", URL: "https://www.youtube.com/@ConceptandCoding", Note: "Free LLD series"},
			{Type: "article", Title: "Refactoring Guru - Patterns", URL: "https://refactoring.guru/design-patterns", Note: "Best patterns reference"},
		},
	},
	{
		Key:     "projects",
		Name:    "Projects",
		Icon:    "🚀",
		Summary: "Strong, deployed projects with clean READMEs set you apart. Depth over quantity — know every layer.",
		Roadmap: []string{
			"Pick 2-3 substantial full-stack projects",
			"Add auth, DB, tests, and CI/CD",
			"Deploy live (Vercel / Render / Fly.io)",
			"Write clear README with architecture diagram",
			"Be ready to explain trade-offs & scaling",
		},
		Resources: []models.Resource{
			{Type: "article", Title: "How to build resume-worthy projects", URL: "https://roadmap.sh/", Note: "Roadmaps per stack"},
			{Type: "repo", Title: "Awesome App Ideas", URL: "https://github.com/florinpop17/app-ideas", Note: "Project inspiration tiers"},
			{Type: "article", Title: "System Design Primer", URL: "https://github.com/donnemartin/system-design-primer", Note: "Design projects to scale"},
			{Type: "practice", Title: "Deploy on Render (free)", URL: "https://render.com/", Note: "Host backend + DB"},
		},
	},
	{
		Key:     "os",
		Name:    "OS Fundamentals",
		Icon:    "💻",
		Summary: "Operating Systems — processes, threads, memory, scheduling. Heavily asked in core CS interview rounds.",
		Roadmap: []string{
			"Process vs Thread, context switching",
			"CPU scheduling algorithms",
			"Concurrency: mutex, semaphore, deadlock",
			"Memory management: paging, segmentation, virtual memory",
			"File systems & I/O",
		},
		Resources: []models.Resource{
			{Type: "playlist", Title: "Gate Smashers - OS", URL: "https://www.youtube.com/playlist?list=PLxCzCOWd7aiGz9donHRrE9I3Mwn6XdP8p", Note: "Complete OS course"},
			{Type: "book", Title: "OSTEP (free book)", URL: "https://pages.cs.wisc.edu/~remzi/OSTEP/", Note: "Three Easy Pieces"},
			{Type: "sheet", Title: "OS Interview Questions", URL: "https://www.geeksforgeeks.org/commonly-asked-operating-systems-interview-questions-set-1/", Note: "Quick revision"},
		},
	},
	{
		Key:     "dbms",
		Name:    "DBMS + SQL",
		Icon:    "🗄️",
		Summary: "Database concepts + writing SQL queries. Normalization, transactions, indexing, and hands-on query practice.",
		Roadmap: []string{
			"ER model, normalization (1NF-BCNF)",
			"Joins, subqueries, aggregations, window functions",
			"Transactions, ACID, isolation levels",
			"Indexing & query optimization",
			"Practice SQL problems daily",
		},
		Resources: []models.Resource{
			{Type: "practice", Title: "LeetCode SQL 50", URL: "https://leetcode.com/studyplan/top-sql-50/", Note: "Must-do query set"},
			{Type: "practice", Title: "SQLZoo", URL: "https://sqlzoo.net/", Note: "Interactive SQL practice"},
			{Type: "practice", Title: "HackerRank SQL", URL: "https://www.hackerrank.com/domains/sql", Note: "Graded challenges"},
			{Type: "playlist", Title: "Gate Smashers - DBMS", URL: "https://www.youtube.com/playlist?list=PLxCzCOWd7aiFAN6I8CuViBuCdJgiOkT2y", Note: "Full DBMS theory"},
		},
	},
	{
		Key:     "hr",
		Name:    "HR Questions",
		Icon:    "🤝",
		Summary: "Behavioral & HR rounds — tell your story with the STAR method. Prepare, don't wing it.",
		Roadmap: []string{
			"Craft your 'Tell me about yourself' pitch",
			"Prepare STAR stories (conflict, failure, leadership)",
			"Know your resume end to end",
			"Research the company & role",
			"Prepare questions to ask the interviewer",
		},
		Resources: []models.Resource{
			{Type: "article", Title: "STAR Method Guide", URL: "https://www.themuse.com/advice/star-interview-method", Note: "Answer framework"},
			{Type: "sheet", Title: "Top 50 HR Questions", URL: "https://www.geeksforgeeks.org/hr-interview-questions/", Note: "With sample answers"},
			{Type: "article", Title: "Amazon Leadership Principles", URL: "https://www.amazon.jobs/content/en/our-workplace/leadership-principles", Note: "Behavioral gold standard"},
		},
	},
}

// GetGuides returns all curated topic study guides.
func GetGuides(c *gin.Context) {
	c.JSON(http.StatusOK, topicGuides)
}
