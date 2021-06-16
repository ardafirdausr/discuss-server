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

type ErrValidation struct {
	Message          string
	ValidationErrors []map[string]string
	Err              error
}

func (evl ErrValidation) Error() string {
	if evl.Message != "" {
		return evl.Message
	} else if evl.Err != nil {
		return evl.Err.Error()
	} else {
		return "Not Found"
	}
}
