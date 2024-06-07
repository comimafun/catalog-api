package circleblock

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"errors"
	"strings"
)

type CircleBlockService interface {
	CreateOne(block string, circleID int) (*entity.CircleBlock, *domain.Error)
	UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	DeleteByID(id int) *domain.Error
	GetOneByBlock(block string) (*entity.CircleBlock, *domain.Error)
	transformBlockString(block string) (*entity.CircleBlock, *domain.Error)
	UpsertOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
}

type circleBlockService struct {
	repo CircleBlockRepo
}

// transformBlockString implements CircleBlockService.
func (c *circleBlockService) transformBlockString(block string) (*entity.CircleBlock, *domain.Error) {
	var postfix string
	var prefix string

	splitted := strings.SplitN(block, "-", 2)
	if len(splitted) != 2 {
		return nil, domain.NewError(400, errors.New("INVALID_BLOCK_FORMAT"), nil)
	}

	prefix = strings.ToUpper(splitted[0])
	postfix = strings.ToLower(splitted[1])

	// check prefix length
	if len(postfix) > 8 || len(prefix) > 2 {
		return nil, domain.NewError(400, errors.New("INVALID_BLOCK_FORMAT"), nil)
	}

	return &entity.CircleBlock{
		Prefix:  prefix,
		Postfix: postfix,
	}, nil
}

// UpsertOne implements CircleBlockService.
func (c *circleBlockService) UpsertOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	panic("implement me")
	// data, err := c.repo.GetOneByBlock(block.Prefix, block.Postfix)
	// if err != nil && !errors.Is(err.Err, gorm.ErrRecordNotFound) {
	// 	return nil, err
	// }
	// if data == nil {
	// 	return c.CreateOne(block)
	// }
	// return c.UpdateOne(data.ID, block)
}

// CreateOne implements CircleBlockService.
func (c *circleBlockService) CreateOne(block string, circleID int) (*entity.CircleBlock, *domain.Error) {
	transform, transformErr := c.transformBlockString(block)
	if transformErr != nil {
		return nil, transformErr
	}

	data, err := c.repo.CreateOne(entity.CircleBlock{
		Prefix:   transform.Prefix,
		Postfix:  transform.Postfix,
		CircleID: circleID,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteByID implements CircleBlockService.
func (c *circleBlockService) DeleteByID(id int) *domain.Error {
	err := c.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

// GetOneByBlock implements CircleBlockService.
func (c *circleBlockService) GetOneByBlock(block string) (*entity.CircleBlock, *domain.Error) {
	transform, transformErr := c.transformBlockString(block)
	if transformErr != nil {
		return nil, transformErr
	}
	data, err := c.repo.GetOneByBlock(transform.Prefix, transform.Postfix)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateOne implements CircleBlockService.
func (c *circleBlockService) UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	data, err := c.repo.UpdateOne(id, block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewCircleBlockService(repo CircleBlockRepo) CircleBlockService {
	return &circleBlockService{repo: repo}
}
