package docker

func isErr(err error) bool {
	return err != nil
}

func isNil(err error) bool {
	return err == nil
}
