package api

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/urandom/readeef"
	"github.com/urandom/readeef/content"
	"github.com/urandom/readeef/content/repo"
	"github.com/urandom/readeef/log"
	"github.com/urandom/readeef/parser"
	"github.com/urandom/readeef/pool"
)

func hubbubRegistration(
	hubbub *readeef.Hubbub,
	service repo.Service,
	log log.Log,
) http.HandlerFunc {
	feedRepo := service.FeedRepo()
	subRepo := service.SubscriptionRepo()

	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		feedID, err := strconv.ParseInt(path.Base(r.URL.Path), 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f, err := feedRepo.Get(content.FeedID(feedID), content.User{})
		if err != nil {
			fatal(w, log, fmt.Sprintf("Error getting feed %d", feedID), err)
			return
		}

		s, err := subRepo.Get(f)
		if err != nil {
			fatal(w, log, "Error getting feed subscription: %+v", err)
			return
		}

		log.Infoln("Receiving hubbub event " + params.Get("hub.mode") + " for " + f.String())

		switch params.Get("hub.mode") {
		case "subscribe":
			if lease, err := strconv.Atoi(params.Get("hub.lease_seconds")); err == nil {
				s.LeaseDuration = int64(lease) * int64(time.Second)
			}
			s.VerificationTime = time.Now()
			s.SubscriptionFailure = false

			w.Write([]byte(params.Get("hub.challenge")))
		case "unsubscribe":
			s.SubscriptionFailure = true

			w.Write([]byte(params.Get("hub.challenge")))
		case "denied":
			w.Write([]byte{})
			log.Printf("Unable to subscribe to '%s': %s\n", params.Get("hub.topic"), params.Get("hub.reason"))
		default:
			w.Write([]byte{})

			buf := pool.Buffer.Get()
			defer pool.Buffer.Put(buf)

			if _, err := buf.ReadFrom(r.Body); err != nil {
				log.Print(err)
				return
			}

			var articles []content.Article
			if pf, err := parser.ParseFeed(buf.Bytes(), parser.ParseRss2, parser.ParseAtom, parser.ParseRss1); err == nil {
				f.Refresh(pf)

				if articles, err = feedRepo.Update(f); err != nil {
					log.Printf("Error updating feed %s: %+v", f, err)
					return
				}
			} else {
				log.Printf("Error parsing feed from subscription %s: %+v", s, err)
				return
			}

			hubbub.ProcessFeedUpdate(f, articles)

			return
		}

		if err = subRepo.Update(s); err != nil {
			log.Printf("Error updating subscription %s: %+v\n", s, err)
			return
		}
	}
}
