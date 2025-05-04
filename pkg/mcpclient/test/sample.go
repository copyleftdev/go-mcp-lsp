package test

func SampleWithError() error {
	err := performOperation()
	return err
}

func SampleWithoutError() error {
	err := performOperation()
	if err != nil {
		return err
	}
	return nil
}

func performOperation() error {
	return nil
}
