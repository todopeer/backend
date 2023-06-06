package model

/*
type Date time.Time

func (d *Date) MarshalGQLContext(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(time.Time(*d).Format(time.DateOnly)))
	return err
}

func (d *Date) UnmarshalGQLContext(ctx context.Context, v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Date must be a string")
	}

	r, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*d = Date(r)
	return nil
}

*/
