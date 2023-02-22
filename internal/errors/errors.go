package errors

func E(err error) {
	if err != nil {
		panic(err)
	}
}
