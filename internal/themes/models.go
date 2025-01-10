package themes

import "net/url"

type GithubRepo struct {
	URL   *url.URL
	Owner string
	Name  string
}
