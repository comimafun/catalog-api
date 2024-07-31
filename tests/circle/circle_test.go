package circle_test

import (
	"catalog-be/internal/database"
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/circle"
	"catalog-be/internal/modules/circle/bookmark"
	"catalog-be/internal/modules/circle/circle_fandom"
	"catalog-be/internal/modules/circle/circle_work_type"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"catalog-be/internal/modules/circle/referral"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"catalog-be/internal/validation"
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

func migrate(t *testing.T, db *gorm.DB) {
	sqlDir := "../../migrator/migrations"

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		t.Fatalf("Error reading directory: %s", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "up.sql") {
			sqlContent, err := os.ReadFile(fmt.Sprintf("%s/%s", sqlDir, file.Name()))
			if err != nil {
				t.Fatalf("Error reading directory: %s", err)
			}

			err = db.Exec(string(sqlContent)).Error
			if err != nil {
				t.Fatalf("Error reading directory: %s", err)
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

func seedEvent(t *testing.T, db *gorm.DB) {
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

func seedFandom(t *testing.T, db *gorm.DB) {
	fandomFile, err := os.Open("./data/fandom.json")
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

func seedWorkType(t *testing.T, db *gorm.DB) {
	workType, err := os.Open("./data/work_type.json")
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

func seedCircle(t *testing.T, db *gorm.DB) {
	file, err := os.Open("./data/circle_initial.json")
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
		// Seed the random number generator to ensure different results each time
		rand.Seed(time.Now().UnixNano())

		// Generate a random number between 1 and 3
		randomNumber := rand.Intn(3) + 1

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

func seedCircleWorkType(t *testing.T, db *gorm.DB) {
	circleWorkTypeFile, err := os.Open("./data/circle_worktype_initial.json")
	if err != nil {
		t.Fatal(err)
	}
	defer circleWorkTypeFile.Close()

	var circleWorkTypes []entity.CircleWorkType

	if err := json.NewDecoder(circleWorkTypeFile).Decode(&circleWorkTypes); err != nil {
		t.Fatal(err)
	}

	circleWorkTypeModels := []entity.CircleWorkType{}

	for _, circleWorkType := range circleWorkTypes {
		circleWorkTypeModel := entity.CircleWorkType{
			CircleID:   circleWorkType.CircleID,
			WorkTypeID: circleWorkType.WorkTypeID,
		}
		circleWorkTypeModels = append(circleWorkTypeModels, circleWorkTypeModel)
	}

	err = db.Create(&circleWorkTypeModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func seedCircleFandom(t *testing.T, db *gorm.DB) {
	circleFandomFile, err := os.Open("./data/circle_fandom_initial.json")
	if err != nil {
		t.Fatal(err)
	}

	defer circleFandomFile.Close()

	var circleFandoms []entity.CircleFandom

	if err := json.NewDecoder(circleFandomFile).Decode(&circleFandoms); err != nil {
		t.Fatal(err)
	}

	circleFandomModels := []entity.CircleFandom{}

	for _, circleFandom := range circleFandoms {
		circleFandomModel := entity.CircleFandom{
			CircleID: circleFandom.CircleID,
			FandomID: circleFandom.FandomID,
		}
		circleFandomModels = append(circleFandomModels, circleFandomModel)
	}

	err = db.Create(&circleFandomModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func seedDataForPagination(t *testing.T, db *gorm.DB) {
	seedEvent(t, db)
	seedFandom(t, db)
	seedWorkType(t, db)
	seedCircle(t, db)
	seedCircleWorkType(t, db)
	seedCircleFandom(t, db)
}

func TestMain(t *testing.M) {
	res := t.Run()
	os.Exit(res)
}

type createCircleInstance struct {
	circleService *circle.CircleService
}

func newCreateCircleInstance(db *gorm.DB) *createCircleInstance {
	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)
	circleRepo := circle.NewCircleRepo(db)
	utils := utils.NewUtils()
	refreshTokenRepo := refreshtoken.NewRefreshTokenRepo(db)
	refreshTokenService := refreshtoken.NewRefreshTokenService(refreshTokenRepo, utils)
	circleWorkTypeRepo := circle_work_type.NewCircleWorkTypeRepo(db)
	circleWorkTypeService := circle_work_type.NewCircleWorkTypeService(
		circleWorkTypeRepo)
	circleFandomRepo := circle_fandom.NewCircleFandomRepo(db)
	circleFandomService := circle_fandom.NewCircleFandomService(circleFandomRepo)
	bookmarkRepo := bookmark.NewCircleBookmarkRepo(db)
	bookmarkService := bookmark.NewCircleBookmarkService(bookmarkRepo)
	validation := validation.NewSanitizer()
	referralRepo := referral.NewReferralRepo(db)
	referralService := referral.NewReferralService(referralRepo)

	circleService := circle.NewCircleService(circleRepo, userService, utils, refreshTokenService, circleWorkTypeService, circleFandomService, bookmarkService, validation, referralService)
	return &createCircleInstance{
		circleService: circleService,
	}
}

func TestCircle(t *testing.T) {
	t.Parallel()

	t.Run("Test pagination", func(t *testing.T) {
		ctx := context.Background()
		connUrl := getConnUrl(t, ctx)
		db := setupDb(t, connUrl)
		seedDataForPagination(t, db)
		instance := newCreateCircleInstance(db)

		t.Run("Test GetPaginatedCircles should return correct metadata", func(t *testing.T) {
			data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
				Page:  1,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)
			assert.LessOrEqual(t, len(data.Data), 20)
			assert.Equal(t, 1, data.Metadata.Page)
			assert.Equal(t, 20, data.Metadata.Limit)
		})

		t.Run("No duplicated circle", func(t *testing.T) {
			t.Run("first page", func(t *testing.T) {
				data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:  1,
					Limit: 20,
				}, 0)

				assert.Nil(t, err)

				for i, circle := range data.Data {
					for j := i + 1; j < len(data.Data); j++ {
						assert.NotEqual(t, circle.ID, data.Data[j].ID)
					}
				}
			})
			t.Run("second page", func(t *testing.T) {
				data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:  2,
					Limit: 20,
				}, 0)

				assert.Nil(t, err)

				for i, circle := range data.Data {
					for j := i + 1; j < len(data.Data); j++ {
						assert.NotEqual(t, circle.ID, data.Data[j].ID)
					}
				}
			})

			t.Run("first to third page", func(t *testing.T) {
				var allDatas []circle_dto.CircleOneForPaginationResponse
				data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:  1,
					Limit: 20}, 0)
				assert.Nil(t, err)

				data2, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page: 2,
				}, 0)
				assert.Nil(t, err)

				data3, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:  3,
					Limit: 20,
				}, 0)
				assert.Nil(t, err)

				allDatas = append(allDatas, data.Data...)
				allDatas = append(allDatas, data2.Data...)
				allDatas = append(allDatas, data3.Data...)

				for i, circle := range allDatas {
					for j := i + 1; j < len(allDatas); j++ {
						assert.NotEqual(t, circle.ID, allDatas[j].ID)
					}
				}
			})
		})

		t.Run("Test last page", func(t *testing.T) {
			data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
				Page:  2,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)

			data, err = instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
				Page:  data.Metadata.TotalPages,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)
			assert.Equal(t, false, data.Metadata.HasNextPage)
		})

		t.Run("Test out of bond page", func(t *testing.T) {
			data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
				Page:  100,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(data.Data))
		})

		t.Run("Test search filter", func(t *testing.T) {
			search := "mag"
			data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
				Page:   1,
				Limit:  20,
				Search: search,
			}, 0)

			assert.Nil(t, err)

			for _, circle := range data.Data {
				assert.Contains(t, strings.ToLower(circle.Name), search)
			}
		})

		t.Run("Test rating filter", func(t *testing.T) {
			t.Run("Test single Rating", func(t *testing.T) {
				rating := "GA"
				data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:   1,
					Limit:  20,
					Rating: []string{rating},
				}, 0)

				assert.Nil(t, err)

				for _, circle := range data.Data {
					assert.Equal(t, rating, *circle.Rating)
					assert.NotEqual(t, "PG", *circle.Rating)
					assert.NotEqual(t, "M", *circle.Rating)
				}
			})
			t.Run("Test multiple Rating", func(t *testing.T) {
				rating := []string{"GA", "PG"}
				data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:   1,
					Limit:  20,
					Rating: rating,
				}, 0)

				assert.Nil(t, err)

				for _, circle := range data.Data {
					assert.Contains(t, rating, *circle.Rating)
					assert.NotEqual(t, "M", *circle.Rating)
				}
			})
		})

		t.Run("Test day filter", func(t *testing.T) {
			t.Run("Test single day", func(t *testing.T) {
				day := entity.Day("first")
				data, err := instance.circleService.GetPaginatedCircle(&circle_dto.FindAllCircleFilter{
					Page:  1,
					Limit: 20,
					Day:   &day,
				}, 0)

				assert.Nil(t, err)

				for _, circle := range data.Data {
					assert.Equal(t, day, *circle.Day)
					assert.NotEqual(t, entity.Day("second"), *circle.Day)
					assert.NotEqual(t, entity.Day("both"), *circle.Day)
				}
			})
		})
	})
}
