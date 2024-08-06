package circle_work_type

type CircleWorkTypeService struct {
	repo *CircleWorkTypeRepo
}

func NewCircleWorkTypeService(repo *CircleWorkTypeRepo) *CircleWorkTypeService {
	return &CircleWorkTypeService{repo: repo}
}
