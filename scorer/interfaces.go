package scorer

type Storage interface {
	CreateOrUpdateUser(name string, email string) error
	ReplaceScore(name string, class string, score float32) error
	UpdateScoreIfHigher(name string, class string, score float32) error
	GetAllScores(class string) error
}

type Notifier interface {
	Started(name string, class string, ranks []float32)
	Updated(name string, score float32)
	Finalized(name string)
}
