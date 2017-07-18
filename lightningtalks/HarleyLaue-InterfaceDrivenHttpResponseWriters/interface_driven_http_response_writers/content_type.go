type ContentType string

func (c ContentType) Response(w http.ResponseWriter, _ *interface{}) error {
	w.Header().Set("Content-Type", string(c))
	return nil
}

