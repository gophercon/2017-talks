type StatusCode int

func (s StatusCode) Response(w http.ResponseWriter, _ *interface{}) error {
	w.WriteHeader(int(s))
	return nil
}

