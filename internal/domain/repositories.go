package domain

type HabitRepository interface {
	Save(habit Habit) (Habit, error)
	FindByID(id int64) (Habit, bool, error)
	FindByNameAndPillar(name string, pillarID int64) (Habit, bool, error)
	FindAll() ([]Habit, error)
	FindByPillar(pillarID int64) ([]Habit, error)
	Delete(id int64) error
}

type TaskRepository interface {
	Save(task Task) (Task, error)
	FindByTitle(title string) (Task, bool, error)
	FindInProgress() ([]Task, error)
	Update(task Task) error
	Delete(id int64) error
}

type PillarRepository interface {
	Save(pillar Pillar) (Pillar, error)
	FindByID(id int64) (Pillar, bool, error)
	FindByName(name string) (Pillar, bool, error)
	FindAll() ([]Pillar, error)
	Delete(id int64) error
}
