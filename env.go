package markup

import (
	"github.com/google/uuid"
)

// Env is the interface that describes ...
type Env interface {
	Component(id uuid.UUID) (c Componer, err error)

	Mount(c Componer) error

	Dismount(c Componer) error

	Update(c Componer) error
}

// NewEnv creates a new environment.
func NewEnv(b CompoBuilder) {
	return &env{
		components: make(map[uuid.UUID]Componer),
		comporoots: make(map[Componer]Tag),
	}
}

type env struct {
	components map[uuid.UUID]Componer
	comporoots map[Componer]Tag
}
