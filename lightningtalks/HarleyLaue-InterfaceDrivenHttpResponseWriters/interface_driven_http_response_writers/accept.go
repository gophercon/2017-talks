type Accept struct {
	ContentType string
	Marshaler   func(interface{}) ([]byte, error)
}

func (a Accept) Response(w http.ResponseWriter, v *interface{}) error {
	if a.ContentType != "" {
		if err := ContentType(a.ContentType).Response(w, nil); err != nil {
			return errors.WithStack(err)
		}
	}

	if a.Marshaler != nil {
		b, err := a.Marshaler(*v)
		if err != nil {
			return errors.WithStack(err)
		}

		*v = string(b)
	}

	return nil
}

