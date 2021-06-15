package entity

type ErrNotFound struct {
	Message string
	Err     error
}

func (ent ErrNotFound) Error() string {
	if ent.Message != "" {
		return ent.Message
	} else if ent.Err != nil {
		return ent.Err.Error()
	} else {
		return "Not Found"
	}
}
