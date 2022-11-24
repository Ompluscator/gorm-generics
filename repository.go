package gorm_generics

import (
	"context"
	"gorm.io/gorm"
)

type GormModel[E any] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}

func NewRepository[M GormModel[E], E any](db *gorm.DB) *GormRepository[M, E] {
	return &GormRepository[M, E]{
		db: db,
	}
}

type GormRepository[M GormModel[E], E any] struct {
	db *gorm.DB
}

func (r *GormRepository[M, E]) Insert(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

// WIP
func (r *GormRepository[M, E]) BatchInsert(ctx context.Context, entities []*E) error {
	convertion := make([]*E, len(entities))
	var entity *E
	for _, chunk := range ChunkSlice(entities, 1000) {
		convertedChunk := make([]M, len(chunk)) //should be declared outside for and emptied here
		for _, e := range chunk {
			var start M
			model := start.FromEntity(*e).(M)
			convertedChunk = append(convertedChunk, model)
		}
		err := r.db.WithContext(ctx).Create(convertedChunk).Error
		if err != nil {
			return err
		}
		unconvertedChunk := make([]*E, len(convertedChunk)) //should be declared outside for and emptied here
		for _, e := range convertedChunk {
			*entity = e.ToEntity() //major error
			unconvertedChunk = append(unconvertedChunk, entity)
		}
		convertion = append(convertion, unconvertedChunk...)
	}

	entities = convertion
	return nil
}

func (r *GormRepository[M, E]) Delete(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)
	err := r.db.WithContext(ctx).Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormRepository[M, E]) DeleteById(ctx context.Context, id any) error {
	var start M
	err := r.db.WithContext(ctx).Delete(&start, &id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GormRepository[M, E]) Update(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

func (r *GormRepository[M, E]) FindByID(ctx context.Context, id any) (E, error) {
	var model M
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return *new(E), err
	}

	return model.ToEntity(), nil
}

func (r *GormRepository[M, E]) Find(ctx context.Context, specification Specification) ([]E, error) {
	return r.FindWithLimit(ctx, &specification, -1, -1)
}

func (r *GormRepository[M, E]) Count(ctx context.Context, specification *Specification) (i int64, err error) {
	model := new(M)
	err = r.getPreWarmDbForSelect(ctx, specification).Model(model).Count(&i).Error
	return
}

func (r *GormRepository[M, E]) getPreWarmDbForSelect(ctx context.Context, specification *Specification) *gorm.DB {
	dbPrewarm := r.db.WithContext(ctx)
	if specification != nil {
		dbPrewarm = dbPrewarm.Where((*specification).GetQuery(), (*specification).GetValues()...)
	}
	return dbPrewarm
}
func (r *GormRepository[M, E]) FindWithLimit(ctx context.Context, specification *Specification, limit int, offset int) ([]E, error) {
	var models []M
	if limit == 0 {
		limit = -1
	}
	if offset == 0 {
		offset = -1
	}

	dbPrewarm := r.getPreWarmDbForSelect(ctx, specification)
	err := dbPrewarm.Limit(limit).Offset(offset).Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make([]E, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToEntity())
	}

	return result, nil
}

func (r *GormRepository[M, E]) FindAll(ctx context.Context) ([]E, error) {
	return r.FindWithLimit(ctx, nil, -1, -1)
}

func (r *GormRepository[M, E]) BatchDelete(ctx context.Context, entities []*E) error {
	for _, chunk := range ChunkSlice(entities, 1000) {
		convertedChunk := make([]M, len(chunk)) //should be declared outside for and emptied here
		for _, e := range chunk {
			var start M
			model := start.FromEntity(*e).(M)
			convertedChunk = append(convertedChunk, model)
		}
		err := r.db.WithContext(ctx).Delete(convertedChunk).Error
		if err != nil {
			return err
		}
	}
	return nil
}
