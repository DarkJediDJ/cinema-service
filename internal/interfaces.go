package internal

import "context"

type Creator interface {
	Create(r Identifiable, ctx context.Context) (Identifiable, error)
}

type Deleter interface {
	Delete(id int64, ctx context.Context) error
}

type Retriever interface {
	Retrieve(id int64, ctx context.Context) (Identifiable, error)
}

type RetrieverAll interface {
	RetrieveAll(ctx context.Context) ([]Identifiable, error)
}

type Service interface {
	Creator
	Deleter
	Retriever
	RetrieverAll
}

type Identifiable interface {
	GID() int64
}
