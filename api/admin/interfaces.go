package admin

type Scorer interface {
	Create(name string, email *string, class string) (string, error)
	Update(id string, score float32) error
	Finalize(id string, replace bool) error
	Cancel(id string) error
}
