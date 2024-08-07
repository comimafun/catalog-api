package report_test

import (
	"catalog-be/internal/database"
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/circle"
	"catalog-be/internal/modules/report"
	"catalog-be/internal/modules/user"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

func migrate(t *testing.T, db *gorm.DB) {
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

func setupDb(t *testing.T, dsn string) *gorm.DB {
	db := database.New(dsn, false)
	migrate(t, db)
	return db
}

func getConnUrl(t *testing.T, ctx context.Context) string {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
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

	return connUrl
}

type circleJson struct {
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

func seedCircle(t *testing.T, db *gorm.DB) {
	file, err := os.Open("../circle/data/circle_initial.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var circles []circleJson
	if err := json.NewDecoder(file).Decode(&circles); err != nil {
		t.Fatal(err)
	}

	circleModels := []entity.Circle{}

	for _, circle := range circles {
		// Create a new random number generator with a seed
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Generate a random number from 1 to 3
		randomNumber := rng.Intn(3) + 1

		cover := "https://cdn.innercatalog.com/development/covers/0190eefb-669a-75e9-a29c-3b244895c4fb.jpg"
		picture := "https://cdn.innercatalog.com/profiles/0190c595-24dd-79f3-8f7a-f2b8f9198d3e.png"

		circleModel := entity.Circle{
			Name:               circle.Name,
			Slug:               circle.Slug,
			Rating:             circle.Rating,
			Day:                &circle.Day,
			Published:          true,
			Verified:           true,
			UsedReferralCodeID: nil,
			CoverPictureURL:    &cover,
			PictureURL:         &picture,
			EventID:            &randomNumber,
		}

		circleModels = append(circleModels, circleModel)
	}

	err = db.Create(&circleModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func seedUser(t *testing.T, db *gorm.DB) {
	// makeCircleID := func(s int) *int {
	// 	return &s
	// }

	userModels := []entity.User{
		{
			ID:    1,
			Name:  "User 1",
			Email: "user1@example.com",
			Hash:  "hash1",
			// CircleID: makeCircleID(1),
		},
		{
			ID:    2,
			Name:  "User 2",
			Email: "user2@example.com",
			Hash:  "hash2",
			// CircleID: makeCircleID(2),
		},
		{
			ID:    3,
			Name:  "User 3",
			Email: "user3@example.com",
			Hash:  "hash3",
			// CircleID: makeCircleID(3),
		},
	}

	err := db.Create(&userModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func seedData(t *testing.T, db *gorm.DB) {
	seedCircle(t, db)
	seedUser(t, db)
}

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func TestReport(t *testing.T) {
	ctx := context.Background()
	connURL := getConnUrl(t, ctx)
	db := setupDb(t, connURL)

	result := db.Find(&entity.Circle{})
	t.Log(result.RowsAffected)

	// seeding data
	seedData(t, db)

	userRepo := user.NewUserRepo(db)
	circleRepo := circle.NewCircleRepo(db)
	repo := report.NewReportRepo(db)
	service := report.NewReportService(repo, circleRepo)

	t.Run("Create new report", func(t *testing.T) {
		user, _ := userRepo.FindOneByID(1)

		var circle entity.Circle
		err := db.Table("circle").First(circle).Error
		assert.Nil(t, err, "get circle should be success")

		reportEntity := entity.Report{
			UserID:   user.ID,
			CircleID: circle.ID,
			Reason:   "this is a test report",
		}
		domainErr := service.CreateReportCircle(&reportEntity)
		assert.Nil(t, domainErr, "create report should be success")
	})

}
