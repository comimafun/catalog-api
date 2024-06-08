package seed

import (
	"catalog-be/internal/entity"
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Seed struct {
	db *gorm.DB
}

func (s *Seed) SeedFandom() {
	file, err := os.Open("public/fandom.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	fandomStrings := make([]string, 0)
	decodeERr := json.NewDecoder(file).Decode(&fandomStrings)

	if decodeERr != nil {
		fmt.Println(decodeERr)
		return
	}

	println("START SEEDING FANDOM")
	for _, fandom := range fandomStrings {
		err := s.db.Create(&entity.Fandom{
			Name:    fandom,
			Visible: true,
		}).Error
		if err != nil {
			fmt.Println(err)
			return
		}
		println("FANDOM SEEDED: ", fandom)
	}

	println("FINISH SEEDING FANDOM")

}

func (s *Seed) SeedWorkType() {
	file, err := os.Open("public/worktype.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	workTypeStrings := make([]string, 0)
	decodeERr := json.NewDecoder(file).Decode(&workTypeStrings)

	if decodeERr != nil {
		fmt.Println(decodeERr)
		return
	}

	println("START SEEDING WORK_TYPE")
	for _, workType := range workTypeStrings {
		err := s.db.Create(&entity.WorkType{
			Name: workType,
		}).Error
		if err != nil {
			fmt.Println(err)
			return
		}
		println("WORK_TYPE SEEDED: ", workType)
	}

	println("FINISH SEEDING WORK_TYPE")

}

func (s *Seed) Run() {
	is_seeding := os.Getenv("SEED")

	fmt.Printf("IS SEEDING: %s\n", is_seeding)
	if is_seeding != "true" {
		return
	}

	fmt.Println("START_SEEDING")

	// s.seedFandom()
	// s.seedWorkType()

	fmt.Println("FINISH_SEEDING")
}

func NewSeed(db *gorm.DB) *Seed {
	return &Seed{db: db}
}
