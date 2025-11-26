package domain

type Status int

const (
	Unknown Status = iota
	Pending
	Active
	Inactive
	Banned
)

func NewStatus(s string) Status {
	switch s {
	case "pending":
		return Pending
	case "active":
		return Active
	case "inactive":
		return Inactive
	case "banned":
		return Banned
	default:
		return Unknown
	}
}

//nolint:exhaustive
func (s Status) String() string {
	switch s {
	case Pending:
		return "pending"
	case Active:
		return "active"
	case Inactive:
		return "inactive"
	case Banned:
		return "banned"
	default:
		return "unknown"
	}
}
