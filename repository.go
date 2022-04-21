package gorm_generics

import (
	"context"

	"gorm.io/gorm"
)

type GormModel[E any] interface {
	ToEntity() E
}

type ModelFactoryMethod[E any, M GormModel[E]] func(entity E) M

func NewRepository[E any, M GormModel[E]](db *gorm.DB, creator ModelFactoryMethod[E, M]) *GormRepository[E, M] {
	return &GormRepository[E, M]{
		creator: creator,
		db:      db,
	}
}

type GormRepository[E any, M GormModel[E]] struct {
	creator ModelFactoryMethod[E, M]
	db      *gorm.DB
}

func (r *GormRepository[E, M]) Insert(ctx context.Context, entity *E) error {
	model := r.creator(*entity)

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

func (r *GormRepository[E, M]) FindByID(ctx context.Context, id uint) (E, error) {
	var model M
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return *new(E), err
	}

	return model.ToEntity(), nil
}

func (r *GormRepository[E, M]) Find(ctx context.Context, specification Specification) ([]E, error) {
	var models []M
	err := r.db.WithContext(ctx).Where(specification.GetQuery(), specification.GetValues()...).Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]E, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToEntity())
	}

	return result, nil
}
