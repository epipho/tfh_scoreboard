package scorer

import (
	"errors"
	"math"
	"sync"

	"github.com/google/uuid"
)

type score struct {
	name        string
	class       string
	val         float32
	incremental bool
}

type InMemoryScorer struct {
	storage  Storage  // persistence for completed scores
	notifier Notifier // notification for score changes

	partialScores map[string]score // scores that are being worked on right now
	lock          sync.RWMutex
}

// NewInMemoryScorer creates a scorer that maintains the intermediate scores
// in an in0memory cache
func NewInMemoryScorer(storage Storage, notifier Notifier) *InMemoryScorer {
	return &InMemoryScorer{
		storage:       storage,
		notifier:      notifier,
		partialScores: map[string]score{},
	}
}

func (s *InMemoryScorer) Create(name string, email *string, class string) (string, error) {
	err := s.storage.CreateOrUpdateUser(name, email)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.partialScores[id] = score{
		val:   -math.MaxFloat32,
		name:  name,
		class: class,
	}

	return id, nil
}

func (s *InMemoryScorer) Update(id string, score float32, incremental bool) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	cur, ok := s.partialScores[id]
	if !ok {
		return errors.New("Invalid score id")
	}
	if incremental {
		if cur.val == -math.MaxFloat32 {
			cur.val = 0
		}
		cur.val += score
		cur.incremental = true
		s.partialScores[id] = cur
	} else {
		if cur.val < score {
			cur.val = score
			s.partialScores[id] = cur
		}
	}

	return nil
}

func (s *InMemoryScorer) Finalize(id string, replace bool) error {
	score, ok := s.partialScores[id]
	if !ok {
		return errors.New("Invalid score id")
	}
	if score.val == -math.MaxFloat32 {
		return errors.New("Score was never set")
	}
	err := s.storage.UpdateScore(score.name, score.class, score.val, score.incremental, replace)
	if err != nil {
		return err
	}
	// remove from map
	s.Cancel(id)
	return nil
}

func (s *InMemoryScorer) Cancel(id string) error {
	delete(s.partialScores, id)
	return nil
}
