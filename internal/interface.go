package internal

type Creator interface {
	Create(r Identifiable) (Identifiable, error)
}

type Deleter interface {
	Delete(id int64) error
}

type Retriever interface {
	Retrieve(id int64) (Identifiable, error)
}

type RetrieverAll interface {
	RetrieveAll() ([]Identifiable, error)
}

type Service interface {
	Creator
	Deleter
	Retriever
	RetrieverAll
}

type Identifiable interface {
	GID() int
}
