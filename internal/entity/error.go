package entity

type ErrBadRequest struct {
	Message string
	Err     error
}

func (ebr ErrBadRequest) Error() string {
	if ebr.Err != nil {
		return ebr.Err.Error()
	} else if ebr.Message != "" {
		return ebr.Message
	} else {
		return "Invalid request"
	}
}

type ErrNotFound struct {
	Message string
	Err     error
}

func (ent ErrNotFound) Error() string {
	if ent.Err != nil {
		return ent.Err.Error()
	} else if ent.Message != "" {
		return ent.Message
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
	if evl.Err != nil {
		return evl.Err.Error()
	} else if evl.Message != "" {
		return evl.Message
	} else {
		return "Not Found"
	}
}
