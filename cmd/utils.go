package cmd

import (
	"github.com/google/go-github/v48/github"
)

// Find statistical mode from list.
func mode[T comparable](l []T) T {
	var maxV T
	freq := make(map[T]int)

	for _, x := range l {
		freq[x] = freq[x] + 1

		if freq[x] > freq[maxV] {
			maxV = x
		}
	}
	return maxV
}

type GetFunction[T any] func() (T, *github.Response, error)

func RateLimitGithubCall[T any](fn GetFunction[T]) (T, *github.Response, error) {
	value, resp, err := fn()
	if err == nil {
		return value, resp, err
	}
	if !isRateLimited(err) && !isRateAbuseLimited(err) {
		return value, resp, err
	}
	WaitIfRateLimited(err)
	return RateLimitGithubCall(fn)
}
