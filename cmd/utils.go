package cmd

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"
)

// Find frequency of occurrences in array.
func frequency[T comparable](l []T) map[T]int {
	freq := make(map[T]int)
	for _, x := range l {
		freq[x] = freq[x] + 1
	}
	return freq
}

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

// map elements in an array.
func mapp[T any, K any](x []T, fn func(T) K) []K {
	result := make([]K, len(x))
	for i, t := range x {
		result[i] = fn(t)
	}
	return result
}

// filter out values from array.
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

// reduce an array to a single value, given an initial value and a combiner function. This is a left-side aggregation.
// i.e. in [a, b, c] we combine (a, b), the combine ((a, b), c)
func reduce[T any, K any](x []T, init K, fn func(T, K) K) K {
	for _, t := range x {
		init = fn(t, init)
	}
	return init
}

// uniq elements from an array.
func uniq[T comparable](x []T) []T {
	exist := make(map[T]uint8)
	var uniqs []T

	for _, t := range x {
		_, inMap := exist[t]
		if !inMap {
			exist[t] = 1
			uniqs = append(uniqs, t)
		}
	}
	return uniqs
}

// reverseGroupBy in a map(K, []V) -> map(V, []K). k in []K does not contain duplicates.
func reverseGroupBy[K comparable, V comparable](x map[K][]V) map[V][]K {
	result := make(map[V][]K)
	for k, vs := range x {
		for _, v := range vs {
			result[v] = append(result[v], k)
		}
	}

	// If v was in array of multiple K's, keep unique list.
	for v, ks := range result {
		result[v] = uniq[K](ks)
	}
	return result
}
