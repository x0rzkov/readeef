package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"readeef"
	"readeef/parser"
	"strconv"
	"strings"

	"github.com/urandom/webfw"
	"github.com/urandom/webfw/context"
	"github.com/urandom/webfw/util"
)

type Feed struct {
	webfw.BaseController
	fm *readeef.FeedManager
}

func NewFeed(fm *readeef.FeedManager) Feed {
	return Feed{
		webfw.NewBaseController("/v:version/feed/*action", webfw.MethodAll, ""),
		fm,
	}
}

type feed struct {
	Id          int64
	Title       string
	Description string
	Link        string
	Image       parser.Image
	Articles    []readeef.Article
}

func (con Feed) Handler(c context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		db := readeef.GetDB(c)
		user := readeef.GetUser(c, r)

		actionParam := webfw.GetParams(c, r)
		parts := strings.Split(actionParam["action"], "/")
		action := parts[0]

	SWITCH:
		switch action {
		case "list":
			var feeds []readeef.Feed
			feeds, err = db.GetUserFeeds(user)

			if err != nil {
				break
			}

			type response struct {
				Feeds []feed
			}

			resp := response{}
			for _, f := range feeds {
				resp.Feeds = append(resp.Feeds, feed{
					Id: f.Id, Title: f.Title, Description: f.Description,
					Link: f.Link, Image: f.Image,
				})
			}

			var b []byte
			b, err = json.Marshal(resp)
			if err != nil {
				break
			}

			w.Write(b)
		case "discover":
			r.ParseForm()

			link := r.FormValue("url")
			var u *url.URL

			/* TODO: non-fatal error */
			if u, err = url.Parse(link); err != nil {
				break
				/* TODO: non-fatal error */
			} else if !u.IsAbs() {
				u.Scheme = "http"
				if u.Host == "" {
					parts := strings.SplitN(u.Path, "/", 2)
					u.Host = parts[0]
					if len(parts) > 1 {
						u.Path = "/" + parts[1]
					} else {
						u.Path = ""
					}
				}
				link = u.String()
			}

			var feeds []readeef.Feed
			feeds, err = con.fm.DiscoverFeeds(link)
			if err != nil {
				break
			}

			type response struct {
				Feeds []feed
			}
			resp := response{}
			for _, f := range feeds {
				resp.Feeds = append(resp.Feeds, feed{
					Id: f.Id, Title: f.Title, Description: f.Description,
					Link: f.Link, Image: f.Image,
				})
			}

			var b []byte
			b, err = json.Marshal(resp)
			if err != nil {
				break
			}

			w.Write(b)
		case "add":
			r.ParseForm()
			links := r.Form["url"]
			success := false

			for _, link := range links {
				/* TODO: non-fatal error */
				var u *url.URL
				if u, err = url.Parse(link); err != nil {
					break SWITCH
					/* TODO: non-fatal error */
				} else if !u.IsAbs() {
					err = errors.New("Feed has no link")
					break SWITCH
				}

				var f readeef.Feed
				f, err = con.fm.AddFeedByLink(link)
				if err != nil {
					break SWITCH
				}

				_, err = db.CreateUserFeed(readeef.GetUser(c, r), f)
				if err != nil {
					break SWITCH
				}
				success = true
			}

			type response struct {
				Success bool
			}
			resp := response{success}

			var b []byte
			b, err = json.Marshal(resp)
			if err != nil {
				break
			}

			w.Write(b)
		case "remove":
			var id int64
			id, err = strconv.ParseInt(parts[1], 10, 64)

			/* TODO: non-fatal error */
			if err != nil {
				break
			}

			var feed readeef.Feed
			feed, err = db.GetUserFeed(id, user)
			/* TODO: non-fatal error */
			if err != nil {
				break
			}

			err = db.DeleteUserFeed(feed)
			/* TODO: non-fatal error */
			if err != nil {
				break
			}

			con.fm.RemoveFeed(feed)

			type response struct {
				Success bool
			}
			resp := response{true}

			var b []byte
			b, err = json.Marshal(resp)
			if err != nil {
				break
			}

			w.Write(b)
		case "opml":
			buf := util.BufferPool.GetBuffer()
			defer util.BufferPool.Put(buf)

			buf.ReadFrom(r.Body)

			var opml parser.Opml
			opml, err = parser.ParseOpml(buf.Bytes())
			if err != nil {
				break
			}

			var feeds []readeef.Feed
			for _, opmlFeed := range opml.Feeds {
				var discovered []readeef.Feed

				discovered, err = con.fm.DiscoverFeeds(opmlFeed.Url)
				if err != nil {
					continue
				}

				feeds = append(feeds, discovered...)
			}

			type response struct {
				Feeds []feed
			}
			resp := response{}
			for _, f := range feeds {
				resp.Feeds = append(resp.Feeds, feed{
					Id: f.Id, Title: f.Title, Description: f.Description,
					Link: f.Link, Image: f.Image,
				})
			}

			var b []byte
			b, err = json.Marshal(resp)
			if err != nil {
				break
			}

			w.Write(b)
		default:
			var articles []readeef.Article

			var limit, offset int

			if len(parts) != 5 {
				err = errors.New(fmt.Sprintf("Expected 5 arguments, got %d", len(parts)))
				break
			}

			limit, err = strconv.Atoi(parts[1])
			/* TODO: non-fatal error */
			if err != nil {
				break
			}

			offset, err = strconv.Atoi(parts[2])
			/* TODO: non-fatal error */
			if err != nil {
				break
			}

			newerFirst := parts[3] == "true"
			unreadOnly := parts[4] == "true"

			if limit > 50 {
				limit = 50
			}

			if action == "__all__" {
				if newerFirst {
					if unreadOnly {
						articles, err = db.GetUnreadUserArticlesDesc(user, limit, offset)
					} else {
						articles, err = db.GetUserArticlesDesc(user, limit, offset)
					}
				} else {
					if unreadOnly {
						articles, err = db.GetUnreadUserArticles(user, limit, offset)
					} else {
						articles, err = db.GetUserArticles(user, limit, offset)
					}
				}
				if err != nil {
					break
				}
			} else if action == "__favorite__" {
				if newerFirst {
					articles, err = db.GetUserFavoriteArticlesDesc(user, limit, offset)
				} else {
					articles, err = db.GetUserFavoriteArticles(user, limit, offset)
				}
				if err != nil {
					break
				}
			} else {
				var f readeef.Feed

				var id int64
				id, err = strconv.ParseInt(action, 10, 64)

				if err != nil {
					err = errors.New("Unknown action " + action)
					break
				}

				f, err = db.GetFeed(id)
				/* TODO: non-fatal error */
				if err != nil {
					break
				}

				f.User = user

				if newerFirst {
					if unreadOnly {
						f, err = db.GetUnreadFeedArticlesDesc(f, limit, offset)
					} else {
						f, err = db.GetFeedArticlesDesc(f, limit, offset)
					}
				} else {
					if unreadOnly {
						f, err = db.GetUnreadFeedArticles(f, limit, offset)
					} else {
						f, err = db.GetFeedArticles(f, limit, offset)
					}
				}
				if err != nil {
					break
				}

				articles = f.Articles
			}

			type response struct {
				Articles []readeef.Article
			}

			resp := response{Articles: articles}

			var b []byte
			b, err = json.Marshal(resp)
			if err != nil {
				break
			}

			w.Write(b)
		}

		if err != nil {
			webfw.GetLogger(c).Print(err)

			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (con Feed) AuthRequired(c context.Context, r *http.Request) bool {
	return true
}
