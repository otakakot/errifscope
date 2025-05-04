package test

func F() error {
	return nil
}

func CallF() {
	ferr := F()
	if ferr != nil { // want "ferr can be scoped with if block"
		panic(ferr)
	}

	ferr = F()
	if ferr != nil {
		panic(ferr)
	}

	if ferr := F(); ferr != nil {
		err := F()
		if err != nil { // want "err can be scoped with if block"
			panic(err)
		}
	}

	if err := F(); err != nil {
		panic(err)
	}

	if err := F(); err != nil {
		panic(err)
	}

	for range 3 {
		err := F()
		if err != nil { // want "err can be scoped with if block"
			panic(err)
		}
	}

	err := F()
	if err == nil {
		return
	}
}

func FF() (string, error) {
	return "", nil
}

func CallFF() {
	_, fferr := FF()
	if fferr != nil { // want "fferr can be scoped with if block"
		panic(fferr)
	}

	_, _ = FF()

	str, err := FF()
	if err != nil {
		panic(err)
	}

	println(str)
}

func FFF() (string, string, error) {
	return "", "", nil
}

func CallFFF() {
	_, _, ffferr := FFF()
	if ffferr != nil { // want "ffferr can be scoped with if block"
		panic(ffferr)
	}

	str1, _, err := FFF()
	if err != nil {
		panic(err)
	}

	println(str1)

	_, str2, err := FFF()
	if err != nil {
		panic(err)
	}

	println(str2)

	str3, str4, err := FFF()
	if err != nil {
		panic(err)
	}

	println(str3, str4)
}

var gerr error

func Groval() {
	gerr = F()
	if gerr != nil {
		panic(gerr)
	}
}
