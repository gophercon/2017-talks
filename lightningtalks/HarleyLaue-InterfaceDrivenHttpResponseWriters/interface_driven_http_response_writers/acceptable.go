func Acceptable(r *http.Request, accepts ...Accept) Responser {
	accept := r.Header.Get("Accept")
	if (accept == "" || accept == "*/*") && len(accepts) > 0 {
		return accepts[0]
	}

	for _, a := range accepts {
		if strings.Contains(accept, a.ContentType) {
			return a
		}
	}

	err := Error{
		Error:      errors.New("unable to find an acceptable media type"),
		StatusCode: http.StatusUnsupportedMediaType,
	}
	if len(accepts) > 0 {
		err.Accept = accepts[0]
	}

	return err
}
