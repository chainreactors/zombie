package smb

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chainreactors/utils/encode"
	"github.com/chainreactors/zombie/pkg"
	"github.com/hirochachacha/go-smb2"
)

// smbSession implements pkg.FileSession over an authenticated SMB2 session.
type smbSession struct {
	service string
	conn    *smb2.Session
}

func (s *smbSession) Service() string  { return s.service }
func (s *smbSession) Raw() interface{} { return s.conn }

func (s *smbSession) Close() error {
	if s.conn != nil {
		return s.conn.Logoff()
	}
	return nil
}

// parseSharePath splits a path like "SHARE/dir/file.txt" into share name and
// the remainder. If no separator is found, the whole string is the share name
// and the relative path is empty.
func parseSharePath(path string) (share, rel string) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "\\")
	idx := strings.IndexAny(path, "/\\")
	if idx < 0 {
		return path, ""
	}
	return path[:idx], path[idx+1:]
}

func (s *smbSession) List(path string) ([]string, error) {
	share, rel := parseSharePath(path)
	if share == "" {
		// No share specified: list available shares.
		names, err := s.conn.ListSharenames()
		if err != nil {
			return nil, err
		}
		return names, nil
	}
	mount, err := s.conn.Mount(share)
	if err != nil {
		return nil, fmt.Errorf("mount %q: %w", share, err)
	}
	defer mount.Umount()

	if rel == "" {
		rel = "."
	}
	entries, err := mount.ReadDir(rel)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

func (s *smbSession) Read(path string) ([]byte, error) {
	share, rel := parseSharePath(path)
	if share == "" || rel == "" {
		return nil, fmt.Errorf("path must include share and file: %q", path)
	}
	mount, err := s.conn.Mount(share)
	if err != nil {
		return nil, fmt.Errorf("mount %q: %w", share, err)
	}
	defer mount.Umount()

	f, err := mount.Open(rel)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// SmbPlugin is stateless; all connection state lives in smbSession.
type SmbPlugin struct{}

func (p *SmbPlugin) Name() string { return "smb" }

// dial establishes a raw TCP connection and performs the SMB2 handshake.
func (p *SmbPlugin) dial(task *pkg.Task, dialer *smb2.Dialer) (*smb2.Session, error) {
	c, err := task.DialTimeout("tcp", task.Address(), time.Duration(task.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	conn, err := dialer.Dial(c)
	if err != nil {
		return nil, err
	}
	// Validate the session by listing shares.
	if _, err := conn.ListSharenames(); err != nil {
		conn.Logoff()
		return nil, err
	}
	return conn, nil
}

func (p *SmbPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	user, domain := pkg.SplitUserDomain(task.Username)

	dialer := &smb2.Dialer{}
	method, pwd := pkg.ParseMethod(task.Password)
	if method == "hash" {
		dialer.Initiator = &smb2.NTLMInitiator{
			User:   user,
			Domain: domain,
			Hash:   encode.HexDecode(pwd),
		}
	} else {
		dialer.Initiator = &smb2.NTLMInitiator{
			User:     user,
			Domain:   domain,
			Password: task.Password,
		}
	}

	conn, err := p.dial(task, dialer)
	if err != nil {
		return nil, err
	}
	return &smbSession{service: task.Service, conn: conn}, nil
}

func (p *SmbPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	user, domain := pkg.SplitUserDomain(task.Username)

	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Domain:   domain,
			Password: "",
		},
	}

	conn, err := p.dial(task, dialer)
	if err != nil {
		return nil, err
	}
	return &smbSession{service: task.Service, conn: conn}, nil
}
