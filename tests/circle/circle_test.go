package circle_test

import (
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
	test_helper "catalog-be/tests/test_helper"
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedCircle(t *testing.T, db *gorm.DB) {
	file, err := os.Open("./data/circle_initial.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var circles []test_helper.CircleJson
	if err := json.NewDecoder(file).Decode(&circles); err != nil {
		t.Fatal(err)
	}

	circleModels := []entity.Circle{}

	for index, circle := range circles {
		// Create a new random number generator with a seed
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Generate a random number from 1 to 3
		randomNumber := rng.Intn(3) + 1

		cover := "https://cdn.innercatalog.com/development/covers/0190eefb-669a-75e9-a29c-3b244895c4fb.jpg"
		picture := "https://cdn.innercatalog.com/profiles/0190c595-24dd-79f3-8f7a-f2b8f9198d3e.png"

		circleModel := entity.Circle{
			ID:                 index + 1,
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
	test_helper.SeedEvent(t, db)
	test_helper.SeedFandom(t, db)
	test_helper.SeedWorkType(t, db)
	seedCircle(t, db)
	seedCircleWorkType(t, db)
	seedCircleFandom(t, db)
}

func TestMain(m *testing.M) {
	res := m.Run()
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

	ctx := context.Background()
	connURL, _ := test_helper.GetConnURL(t, ctx)
	db := test_helper.SetupDb(t, connURL)

	// Seed database with initial data
	seedDataForPagination(t, db)
	instance := newCreateCircleInstance(db)

	t.Run("Test pagination", func(t *testing.T) {
		t.Run("Test GetPaginatedCircles should return correct metadata", func(t *testing.T) {
			data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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

			t.Run("all pages", func(t *testing.T) {
				var allDatas []circle_dto.CirclePaginatedResponse
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
					Page:  1,
					Limit: 20}, 0)
				assert.Nil(t, err)

				allDatas = append(allDatas, data.Data...)

				for i := 2; i <= data.Metadata.TotalPages; i++ {
					data, err = instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
						Page:  i,
						Limit: 20,
					}, 0)
					assert.Nil(t, err)
					allDatas = append(allDatas, data.Data...)
				}

				for i, circle := range allDatas {
					for j := i + 1; j < len(allDatas); j++ {
						assert.NotEqual(t, circle.ID, allDatas[j].ID)
					}
				}
			})
		})

		t.Run("Test last page", func(t *testing.T) {
			data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
				Page:  2,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)

			data, err = instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
				Page:  data.Metadata.TotalPages,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)
			assert.Equal(t, false, data.Metadata.HasNextPage)
		})

		t.Run("Test out of bond page", func(t *testing.T) {
			data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
				Page:  100,
				Limit: 20,
			}, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(data.Data))
		})

		t.Run("Test search filter", func(t *testing.T) {
			search := "mag"
			data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
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

		t.Run("Test filter by work type", func(t *testing.T) {
			t.Run("Test single work type", func(t *testing.T) {
				workTypeIDS := []int{1}
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
					Page:        1,
					Limit:       20,
					WorkTypeIDs: workTypeIDS,
				}, 0)
				assert.Nil(t, err)

				for _, circle := range data.Data {
					workTypePaginationMap := make(map[int]bool)

					for _, workType := range circle.WorkType {
						workTypePaginationMap[workType.ID] = true
					}

					for _, workTypeID := range workTypeIDS {
						assert.Equal(t, true, workTypePaginationMap[workTypeID])
					}
				}

			})

			t.Run("Test multiple work type", func(t *testing.T) {
				workTypeIDS := []int{1, 3, 5, 10}
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
					Page:        1,
					Limit:       20,
					WorkTypeIDs: workTypeIDS,
				}, 0)
				assert.Nil(t, err)

				for _, circle := range data.Data {
					currentCircleWorkType := make(map[int]bool)
					for _, workType := range circle.WorkType {
						currentCircleWorkType[workType.ID] = true
					}

					found := false
					for _, workTypeID := range workTypeIDS {
						if currentCircleWorkType[workTypeID] {
							found = true
							break
						}
					}

					assert.True(t, found)
				}

			})

			t.Run("Test filter by fandom", func(t *testing.T) {
				t.Run("Test single fandom", func(t *testing.T) {
					// create random fandom ids from 1-94
					randomFandomID := rand.Intn(94) + 1
					fandomIDS := []int{randomFandomID}

					data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
						Page:      1,
						Limit:     20,
						FandomIDs: fandomIDS,
					}, 0)
					assert.Nil(t, err)

					for _, circle := range data.Data {
						currentCircleFandom := make(map[int]bool)
						for _, fandom := range circle.Fandom {
							currentCircleFandom[fandom.ID] = true
						}

						found := false
						for _, fandomID := range fandomIDS {
							if currentCircleFandom[fandomID] {
								found = true
								break
							}
						}

						assert.True(t, found)
					}
				})

				t.Run("Test multiple fandom", func(t *testing.T) {

					var fandomIDS []int
					maxFandomID := rand.Intn(94) + 1
					for i := 1; i <= maxFandomID; i++ {
						fandomIDS = append(fandomIDS, i)
					}

					data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
						Page:      1,
						Limit:     20,
						FandomIDs: fandomIDS,
					}, 0)
					assert.Nil(t, err)

					for _, circle := range data.Data {
						currentCircleFandom := make(map[int]bool)
						for _, fandom := range circle.Fandom {
							currentCircleFandom[fandom.ID] = true
						}

						found := false
						for _, fandomID := range fandomIDS {
							if currentCircleFandom[fandomID] {
								found = true
								break
							}
						}

						assert.True(t, found)
					}
				})
			})

		})

		t.Run("Test filter by event", func(t *testing.T) {
			t.Run("Test single event", func(t *testing.T) {
				event := "event-1"
				data, err := instance.circleService.GetPaginatedCircles(&circle_dto.GetPaginatedCirclesFilter{
					Page:  1,
					Limit: 20,
					Event: event,
				}, 0)

				assert.Nil(t, err)

				for _, circle := range data.Data {
					assert.Equal(t, event, circle.Event.Slug)
				}
			})
		})
	})
}
