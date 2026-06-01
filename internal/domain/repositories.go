package domain

type HabitRepository interface {
	Save(habit Habit) (Habit, error)
	FindByNameAndPillar(name string, pillarID int64) (Habit, bool, error)
	FindAll() ([]Habit, error)
	FindByPillar(pillarID int64) ([]Habit, error)
	Delete(id int64) error
}

type PillarRepository interface {
	Save(pillar Pillar) (Pillar, error)
	FindByID(id int64) (Pillar, bool, error)
	FindByName(name string) (Pillar, bool, error)
	FindAll() ([]Pillar, error)
	Delete(id int64) error
}
