type Responser interface {
	Response(http.ResponseWriter, *interface{}) error
}

func Write(w http.ResponseWriter, v interface{}, responses ...Responser) error {
	for _, r := range responses {
		if err := r.Response(w, &v); err != nil {
			return errors.WithStack(err)
		}
	}

	_, err := fmt.Fprint(w, v)
	return errors.WithStack(err)
}
