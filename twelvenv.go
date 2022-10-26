package twelvenv // import "go.izbudki.dev/twelvenv"

import "go.izbudki.dev/twelvenv/internal"

func FromEnv(s interface{}) error {
	fields, err := internal.ParseFields(s)
	if err != nil {
		return err
	}

	values, err := internal.ReadEnvironment(fields)
	if err != nil {
		return err
	}

	err = internal.SetValues(s, values)
	if err != nil {
		return err
	}

	return nil
}
