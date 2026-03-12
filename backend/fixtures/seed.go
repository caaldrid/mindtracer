package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/setup"
)

func dueIn(days int) *time.Time {
	t := time.Now().AddDate(0, 0, days)
	return &t
}

const testUserEmail string = "seed@test.local"

// JSON structs — seed-data shape, lives here not in models

type todoJSON struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	DueInDays       *int      `json:"due_in_days"`
	CompletedInDays *int      `json:"completed_in_days"`
	Prerequisite    *todoJSON `json:"prerequisite"`
}

type projectJSON struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Tasks        []todoJSON   `json:"tasks"`
	Prerequisite *projectJSON `json:"prerequisite"`
}

type areaJSON struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Projects    []projectJSON `json:"projects"`
}

func (t *todoJSON) toSeed() todoSeed {
	seed := todoSeed{
		ToDo: &models.ToDo{
			Title:       t.Title,
			Description: t.Description,
			Status:      models.Status(t.Status),
		},
	}
	if t.DueInDays != nil {
		seed.DueDate = dueIn(*t.DueInDays)
	}
	if t.CompletedInDays != nil {
		seed.CompletedAt = dueIn(*t.CompletedInDays)
	}
	if t.Prerequisite != nil {
		prereq := t.Prerequisite.toSeed()
		seed.prerequisite = &prereq
	}
	return seed
}

func (p *projectJSON) toSeed() projectSeed {
	seed := projectSeed{
		Project: &models.Project{
			Name:        p.Name,
			Description: p.Description,
		},
	}
	for _, t := range p.Tasks {
		seed.tasks = append(seed.tasks, t.toSeed())
	}
	if p.Prerequisite != nil {
		prereq := p.Prerequisite.toSeed()
		seed.prerequisite = &prereq
	}
	return seed
}

func (a *areaJSON) toSeed() areaSeed {
	seed := areaSeed{
		Area: &models.Area{
			Name:        a.Name,
			Description: a.Description,
		},
	}
	for _, p := range a.Projects {
		seed.projects = append(seed.projects, p.toSeed())
	}
	return seed
}

// Seed structs — used during DB insertion

type todoSeed struct {
	*models.ToDo
	prerequisite *todoSeed
}

type projectSeed struct {
	*models.Project
	tasks        []todoSeed
	prerequisite *projectSeed
}

type areaSeed struct {
	*models.Area
	projects []projectSeed
}

func loadAreas(path string) ([]areaSeed, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw []areaJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	areas := make([]areaSeed, len(raw))
	for i, a := range raw {
		areas[i] = a.toSeed()
	}
	return areas, nil
}

func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	bytePassword := make([]byte, length)
	_, err := rand.Read(bytePassword)
	if err != nil {
		return "", err
	}
	for i := range length {
		bytePassword[i] = charset[int(bytePassword[i])%len(charset)]
	}
	return string(bytePassword), nil
}

func main() {
	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	teardown := flag.Bool("teardown", false, "Truncate all tables and exit")
	flag.Parse()

	DB, err := setup.ConnectDB(config)
	if err != nil {
		log.Fatal("Could connect to database instance", err)
	}
	cntx := context.Background()

	if *teardown {
		tables, err := DB.Migrator().GetTables()
		if err != nil {
			log.Fatalf("Could not retrieve tables: %s", err.Error())
		}
		if len(tables) == 0 {
			log.Println("No tables found — nothing to clear.")
			return
		}
		for _, table := range tables {
			if err := DB.WithContext(cntx).Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
				log.Fatalf("Failed to truncate %s: %s", table, err.Error())
			}
			log.Printf("Truncated %s", table)
		}
		log.Println("Database cleared.")
		return
	}

	areas, err := loadAreas("fixtures/seed_data.json")
	if err != nil {
		log.Fatalf("Could not load seed data: %s", err.Error())
	}

	usr := &models.User{Email: testUserEmail, UserName: "Seed User"}

	result := DB.WithContext(cntx).Where("email = ? ", testUserEmail).First(usr)

	if result.Error == nil {
		log.Println("Seed data is already present in the DB")
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("Seeding Db...")

		randPass, err := generatePassword(32)
		if err != nil {
			log.Fatalf("Failed to generate a password: %s", err.Error())
		}
		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(randPass),
			bcrypt.DefaultCost,
		)
		if err != nil {
			log.Fatalf("Could not encrypt seed user pswd: %s", err.Error())
		}
		usr.Password = string(passwordHash)

		result = DB.WithContext(cntx).Create(usr)

		if result.Error != nil {
			log.Fatalf("Failed to create a user in the DB: %s", result.Error.Error())
		}

		createProject := func(project *projectSeed, areaID uuid.UUID, userID uuid.UUID) {
			project.AreaID = areaID
			project.UserID = userID

			result := DB.WithContext(cntx).Create(project.Project)
			if result.Error != nil {
				log.Fatalf(
					"Failed to create a Project %s in the DB: %s",
					project.Name,
					result.Error.Error(),
				)
			}

			createToDo := func(todo *todoSeed) {
				todo.ProjectID = project.ID

				result := DB.WithContext(cntx).Create(todo.ToDo)
				if result.Error != nil {
					log.Fatalf(
						"Failed to create a ToDo %s in the DB: %s",
						todo.Title,
						result.Error.Error(),
					)
				}
			}

			for _, task := range project.tasks {
				if task.prerequisite != nil {
					createToDo(task.prerequisite)
					task.PrerequisiteID = &task.prerequisite.ID
				}
				createToDo(&task)
			}
		}

		for _, area := range areas {
			area.UserID = usr.ID

			result = DB.WithContext(cntx).Create(area.Area)
			if result.Error != nil {
				log.Fatalf(
					"Failed to create a Area %s in the DB: %s",
					area.Name,
					result.Error.Error(),
				)
			}

			for _, project := range area.projects {
				if project.prerequisite != nil {
					createProject(project.prerequisite, area.ID, usr.ID)
					project.PrerequisiteID = &project.prerequisite.ID
				}

				createProject(&project, area.ID, usr.ID)
			}
		}

	} else {
		log.Fatalf("%s", result.Error.Error())
	}
}
