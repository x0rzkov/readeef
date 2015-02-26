package base

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/urandom/readeef/content/info"
)

type Subscription struct {
	Error

	info           info.Subscription
	callbackPrefix string
}

func (s Subscription) String() string {
	return fmt.Sprintf("Subscription for %s\n", s.info.Link)
}

func (s *Subscription) Info(in ...info.Subscription) info.Subscription {
	if s.Err() != nil {
		return s.info
	}

	if len(in) > 0 {
		s.info = in[0]
	}

	return s.info
}

func (s *Subscription) Validate() error {
	if u, err := url.Parse(s.info.Link); err != nil || !u.IsAbs() {
		return ValidationError{errors.New("Invalid subscription link")}
	}

	if s.info.FeedId == 0 {
		return ValidationError{errors.New("Invalid feed id")}
	}

	return nil
}
