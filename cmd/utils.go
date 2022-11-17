package cmd

import (
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"
)

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
