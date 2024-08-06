package circle_fandom

type CircleFandomService struct {
	repo *CircleFandomRepo
}

func NewCircleFandomService(repo *CircleFandomRepo) *CircleFandomService {
	return &CircleFandomService{repo: repo}
}
