package todo

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
