package cmd

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"cdr.dev/coder-cli/internal/config"
	"cdr.dev/coder-cli/internal/loginsrv"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"go.coder.com/flog"
)

func makeLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [Coder Enterprise URL eg. http://my.coder.domain/]",
		Short: "Authenticate this client for future operations",
		Args:  cobra.ExactArgs(1),
		RunE:  login,
	}
	return cmd
}

func login(cmd *cobra.Command, args []string) error {
	rawURL := args[0]
	if rawURL == "" || !strings.HasPrefix(rawURL, "http") {
		return xerrors.Errorf("invalid URL")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return xerrors.Errorf("parse url: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return xerrors.Errorf("create login server: %+v", err)
	}
	defer listener.Close()

	srv := &loginsrv.Server{
		TokenCond: sync.NewCond(&sync.Mutex{}),
	}
	go func() {
		_ = http.Serve(
			listener, srv,
		)
	}()

	err = config.URL.Write(
		(&url.URL{Scheme: u.Scheme, Host: u.Host}).String(),
	)
	if err != nil {
		return xerrors.Errorf("write url: %v", err)
	}

	authURL := url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     "/internal-auth/",
		RawQuery: "local_service=http://" + listener.Addr().String(),
	}

	err = browser.OpenURL(authURL.String())
	if err != nil {
		// Tell the user to visit the URL instead.
		flog.Info("visit %s to login", authURL.String())
	}
	srv.TokenCond.L.Lock()
	srv.TokenCond.Wait()
	err = config.Session.Write(srv.Token)
	srv.TokenCond.L.Unlock()
	if err != nil {
		return xerrors.Errorf("set session: %v", err)
	}
	flog.Success("logged in")
	return nil
}
