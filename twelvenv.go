package twelvenv

import "github.com/izbudki/twelvenv/internal"

func Unmarshal(s interface{}) error {
	fields, err := internal.ParseFields(s, nil)
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
