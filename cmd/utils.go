package cmd

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"
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
	log.Debug().Str("frequency", fmt.Sprint(freq)).Msg("frequency of emails in repos")
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
	log.Warn().Err(err).Bool("is-abuse-rate-limited", isRateAbuseLimited(err)).Msg("Throttled by github")
	WaitIfRateLimited(err)
	return RateLimitGithubCall(fn)
}

func mapp[T any, K any](x []T, fn func(T) K) []K {
	result := make([]K, len(x))
	for i, t := range x {
		result[i] = fn(t)
	}
	return result
}

func filter[T any](x []T, fn func(T) bool) []T {
	result := make([]T, len(x))
	i := 0
	for _, t := range x {
		if fn(t) {
			result[i] = t
			i++
		}
	}
	return result[:i]
}

func reduce[T any, K any](x []T, init K, fn func(T, K) K) K {
	for _, t := range x {
		init = fn(t, init)
	}
	return init
}
