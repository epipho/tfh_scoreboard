package scorer

type InMemoryScorer struct {
	storage  Storage  // persistence for completed scores
	notifier Notifier // notification for score changes
}

// NewInMemoryScorer creates a scorer that maintains the intermediate scores
// in an in0memory cache
func NewInMemoryScorer(storage Storage, notifier Notifier) *InMemoryScorer {
	return &InMemoryScorer{
		storage:  storage,
		notifier: notifier,
	}
}

func (s *InMemoryScorer) Create(name string, email string, class string) (string, error) {
	return "", nil
}

func (s *InMemoryScorer) Update(id string, score float32, incremental bool) error {
	return nil
}

func (s *InMemoryScorer) Finalize(id string, replace bool) error {
	return nil
}

func (s *InMemoryScorer) Cancel(id string) error {
	return nil
}
