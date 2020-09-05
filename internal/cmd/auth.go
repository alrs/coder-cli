package cmd

import (
	"net/url"

	"cdr.dev/coder-cli/coder-sdk"
	"cdr.dev/coder-cli/internal/config"
	"golang.org/x/xerrors"
)

func newClient() (*coder.Client, error) {
	sessionToken, err := config.Session.Read()
	if err != nil {
		return nil, xerrors.Errorf("read session: %v (did you run coder login?)", err)
	}

	rawURL, err := config.URL.Read()
	if err != nil {
		return nil, xerrors.Errorf("read url: %v (did you run coder login?)", err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, xerrors.Errorf("url misformatted: %v (try runing coder login)", err)
	}

	client := &coder.Client{
		BaseURL: u,
		Token:   sessionToken,
	}

	return client, nil
}
