type Header http.Header

func (h Header) Response(w http.ResponseWriter, _ *interface{}) error {
	for k, v := range h {
		w.Header().Del(k)
		for _, v := range v {
			w.Header().Add(k, v)
		}
	}

	return nil
}
