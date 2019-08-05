package scorer

type Storage interface {
	CreateOrUpdateUser(name string, email *string) error
	UpdateScore(name string, class string, score float32, replace bool) error
	GetAllScores(class string) ([]interface {
		Score() float32
		Attempts() int
	}, error)
}

type Notifier interface {
	Started(name string, class string, ranks []float32)
	Updated(score float32)
	Finalized()
}
