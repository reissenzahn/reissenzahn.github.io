package main

import "errors"

// align the happy path to the left;  you should quickly be able to scan down one column to see the expected execution flow

func concat(s, t string) (string, error) {
	return s + t, nil
}

func problem(s, t string, max int) (string, error) {
	if s == "" {
		return "", errors.New("s is empty")
	} else {
		if t == "" {
			return "", errors.New("t is empty")
		} else {
			r, err := concat(s, t)
			if err != nil {
				return "", err
			} else {
				if len(r) > max {
					return r[:max], nil
				} else {
					return r, nil
				}
			}
		}
	}
}

// reduce the number of nested blocks, align the happy path to the left and return as early as possible
func solution(s, t string, max int) (string, error) {
	if s == "" {
		return "", errors.New("s is empty")
	}

	if t == "" {
		return "", errors.New("t is empty")
	}

	r, err := concat(s, t)
	if err != nil {
		return "", err
	}

	if len(r) > max {
		return r[:max], nil
	}

	return r, nil
}

func main() {

}
