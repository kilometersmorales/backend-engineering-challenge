package util

// Panic on error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
