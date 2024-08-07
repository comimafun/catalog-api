package testhelper

import (
	"catalog-be/internal/database"
	"catalog-be/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

type CircleJson struct {
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Rating           *string    `json:"rating"` // enum GA, PG, M
	Day              entity.Day `json:"day"`
	Comission        bool       `json:"comission"`
	Comic            bool       `json:"comic"`
	Artbook          bool       `json:"artbook"`
	PhotobookGeneral bool       `json:"photobook_general"`
	Novel            bool       `json:"novel"`
	Game             bool       `json:"game"`
	Music            bool       `json:"music"`
	Goods            bool       `json:"goods"`
	HandmadeCrafts   bool       `json:"handmade_crafts"`
	PhotobookCosplay bool       `json:"photobook_cosplay"`
	Fandom           string     `json:"fandom"`
	WorkTypeIDs      []int      `json:"work_type_ids"`
}

func Migrate(t *testing.T, db *gorm.DB) {
	sqlDir := "../../migrator/migrations"

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "up.sql") {
			sqlContent, err := os.ReadFile(fmt.Sprintf("%s/%s", sqlDir, file.Name()))
			if err != nil {
				t.Fatalf("Error reading file: %s", err)
			}

			err = db.Exec(string(sqlContent)).Error
			if err != nil {
				t.Fatalf("Error executing SQL: %s", err)
			}
		}
	}
}

func SetupDb(t *testing.T, dsn string) *gorm.DB {
	db := database.New(dsn, false)
	Migrate(t, db)
	return db
}

func SeedEvent(t *testing.T, db *gorm.DB) {
	err := db.Create([]entity.Event{
		{
			Name: "Event 1",
			Slug: "event-1",
		},
		{
			Name: "Event 2",
			Slug: "event-2",
		}, {
			Name: "Event 3",
			Slug: "event-3",
		},
	}).Error

	if err != nil {
		t.Fatal(err)
	}
}

func SeedFandom(t *testing.T, db *gorm.DB) {
	fandomFile, err := os.Open("../circle/data/fandom.json")
	if err != nil {
		t.Fatal(err)
	}
	defer fandomFile.Close()

	var fandoms []entity.Fandom
	if err := json.NewDecoder(fandomFile).Decode(&fandoms); err != nil {
		t.Fatal(err)
	}

	fandomModels := []entity.Fandom{}
	for _, fandom := range fandoms {
		fandomModel := entity.Fandom{
			Name:    fandom.Name,
			Visible: true,
		}
		fandomModels = append(fandomModels, fandomModel)
	}
	err = db.Create(&fandomModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func SeedWorkType(t *testing.T, db *gorm.DB) {
	workType, err := os.Open("../circle/data/work_type.json")
	if err != nil {
		t.Fatal(err)
	}
	defer workType.Close()

	var workTypes []entity.WorkType
	if err := json.NewDecoder(workType).Decode(&workTypes); err != nil {
		t.Fatal(err)
	}

	workTypeModels := []entity.WorkType{}

	for _, workType := range workTypes {
		workTypeModel := entity.WorkType{
			Name: workType.Name,
		}
		workTypeModels = append(workTypeModels, workTypeModel)
	}

	err = db.Create(&workTypeModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func SeedUser(t *testing.T, db *gorm.DB) {
	makeCircleID := func(s int) *int {
		return &s
	}

	userModels := []entity.User{
		{
			ID:       1,
			Name:     "User 1",
			Email:    "user1@example.com",
			Hash:     "hash1",
			CircleID: makeCircleID(1),
		},
		{
			ID:       2,
			Name:     "User 2",
			Email:    "user2@example.com",
			Hash:     "hash2",
			CircleID: makeCircleID(2),
		},
		{
			ID:       3,
			Name:     "User 3",
			Email:    "user3@example.com",
			Hash:     "hash3",
			CircleID: makeCircleID(3),
		},
	}

	err := db.Create(&userModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func GetConnURL(t *testing.T, ctx context.Context) (string, testcontainers.Container) {
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("foobar"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
	}

	connUrl, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Could not get connection string: %s", err)
	}

	return connUrl, container
}
