package test

func f() error {
	return nil
}

func testf() {
	ferr := f()
	if ferr != nil { // want "ferr can be scoped with if block"
		return
	}

	if ferr := f(); ferr != nil {
		err := f()
		if err != nil { // want "err can be scoped with if block"
			return
		}
	}

	if err := f(); err != nil {
		println(err)
	}

	if err := f(); err != nil {
		println(err)
	}
}

func ff() (string, error) {
	return "", nil
}

func testff() {
	_, fferr := ff()
	if fferr != nil { // want "fferr can be scoped with if block"
		return
	}

	str, err := ff()
	if err != nil {
		return
	}

	println(str)
}

func fff() (string, string, error) {
	return "", "", nil
}

func testfff() {
	_, _, ffferr := fff()
	if ffferr != nil { // want "ffferr can be scoped with if block"
		return
	}

	str1, _, err := fff()
	if err != nil {
		return
	}

	println(str1)

	_, str2, err := fff()
	if err != nil {
		return
	}

	println(str2)

	str3, str4, err := fff()
	if err != nil {
		return
	}

	println(str3, str4)
}
