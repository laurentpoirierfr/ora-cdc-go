package helper

func DieIfError(err error) {
	if err != nil {
		panic(err)
	}
}
