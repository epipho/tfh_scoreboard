package scorer

type Storage interface {
	CreateOrUpdateUser(name string, email *string) error
	UpdateScore(name string, class string, score float32, incremental bool, replace bool) error
	GetAllScores(class string) error
}

type Notifier interface {
	Started(name string, class string, ranks []float32)
	Updated(score float32)
	Finalized()
}
