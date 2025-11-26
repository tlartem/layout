package postgres

type Postgres struct{}

func New() *Postgres {
	return &Postgres{}
}
