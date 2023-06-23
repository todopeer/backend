package highorder

// for function compositions

// All runs fns one by one, would return the first error it encounters
func All(fns ...func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

func Branch(cond bool, trueBranch func() error, falseBranch func() error) func() error {
	return BranchF(func() bool { return cond }, trueBranch, falseBranch)
}

func BranchF(condF func() bool, trueBranch func() error, falseBranch func() error) func() error {
	return func() error {
		if condF() {
			return trueBranch()
		}
		if falseBranch != nil {
			return falseBranch()
		}
		return nil
	}
}
