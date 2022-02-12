package internal

type Creator interface {
	Create(r Identifiable) (Identifiable, error)
}

type Service interface {
	Creator
	// Updater
	// Deleter
}

type Identifiable interface {
	GID() int
}
