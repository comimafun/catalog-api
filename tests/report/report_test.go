package report_test

import (
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/circle"
	"catalog-be/internal/modules/report"
	test_helper "catalog-be/tests/test_helper"
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedCircle(t *testing.T, db *gorm.DB) {
	file, err := os.Open("../circle/data/circle_initial.json")
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
		}

		circleModels = append(circleModels, circleModel)
	}

	err = db.Create(&circleModels).Error
	if err != nil {
		t.Fatal(err)
	}
}

func seedData(t *testing.T, db *gorm.DB) {
	seedCircle(t, db)
	test_helper.SeedUser(t, db)
}

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func TestReport(t *testing.T) {
	ctx := context.Background()
	connURL, _ := test_helper.GetConnURL(t, ctx)
	db := test_helper.SetupDb(t, connURL)

	// seeding data
	seedData(t, db)

	circleRepo := circle.NewCircleRepo(db)
	repo := report.NewReportRepo(db)
	service := report.NewReportService(repo, circleRepo)

	t.Run("Create new report", func(t *testing.T) {
		reportEntity := entity.Report{
			UserID:   1,
			CircleID: 1,
			Reason:   "this is a test report",
		}

		err := service.CreateReportCircle(&reportEntity)
		assert.Nil(t, err)
	})

	t.Run("Circle not found", func(t *testing.T) {
		reportEntity := entity.Report{
			UserID:   1,
			CircleID: 1000,
			Reason:   "this is a test report",
		}

		err := service.CreateReportCircle(&reportEntity)
		assert.NotNil(t, err)
		assert.Equal(t, 404, err.Code)
		assert.Equal(t, errors.New("CIRCLE_NOT_FOUND"), err.Err)
	})

}
