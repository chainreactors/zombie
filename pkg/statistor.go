package pkg

import (
	"errors"
	"fmt"
	"strings"
)

type ErrCategory int

const (
	ErrCatOther ErrCategory = iota
	ErrCatTimeout
	ErrCatRefused
	ErrCatAuth
)

type Statistor struct {
	Total   int
	Success int
	Cur     string
	Tasks   map[string]int

	ErrTimeout int
	ErrRefused int
	ErrAuth    int
	ErrOther   int

	Extracteds int
	Loot       int
}

func (stat *Statistor) RecordResult(result *Result) {
	if result.OK {
		stat.Success++
		stat.Extracteds += len(result.Extracteds)
		stat.Loot += len(result.Loot)
	} else {
		stat.RecordError(result.Err)
	}
}

func (stat *Statistor) RecordError(err error) {
	if err == nil {
		return
	}
	switch ClassifyError(err) {
	case ErrCatTimeout:
		stat.ErrTimeout++
	case ErrCatRefused:
		stat.ErrRefused++
	case ErrCatAuth:
		stat.ErrAuth++
	default:
		stat.ErrOther++
	}
}

func (stat *Statistor) ErrorString() string {
	total := stat.ErrTimeout + stat.ErrRefused + stat.ErrAuth + stat.ErrOther
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("errors: timeout=%d, refused=%d, auth_fail=%d, other=%d",
		stat.ErrTimeout, stat.ErrRefused, stat.ErrAuth, stat.ErrOther)
}

func (stat *Statistor) SummaryString() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("total: %d, success: %d", stat.Total, stat.Success))
	if stat.Extracteds > 0 || stat.Loot > 0 {
		parts = append(parts, fmt.Sprintf("extracteds: %d, loot: %d", stat.Extracteds, stat.Loot))
	}
	if errStr := stat.ErrorString(); errStr != "" {
		parts = append(parts, errStr)
	}
	return strings.Join(parts, ", ")
}

func (stat *Statistor) TaskString() string {
	var s strings.Builder
	for k, v := range stat.Tasks {
		s.WriteString(fmt.Sprintf("%s:%d ", k, v))
	}
	return s.String()
}

func ClassifyError(err error) ErrCategory {
	if err == nil {
		return ErrCatOther
	}

	var te TimeoutError
	if errors.As(err, &te) {
		return ErrCatTimeout
	}
	if errors.Is(err, ErrorWrongUserOrPwd) {
		return ErrCatAuth
	}

	msg := err.Error()
	switch {
	case strings.Contains(msg, "i/o timeout"),
		strings.Contains(msg, "deadline exceeded"),
		strings.Contains(msg, "context deadline exceeded"):
		return ErrCatTimeout
	case strings.Contains(msg, "connection refused"),
		strings.Contains(msg, "connection reset"),
		strings.Contains(msg, "no route to host"):
		return ErrCatRefused
	case strings.Contains(msg, "unable to authenticate"),
		strings.Contains(msg, "Access denied"),
		strings.Contains(msg, "authentication fail"),
		strings.Contains(msg, "login fail"),
		strings.Contains(msg, "wrong username"),
		strings.Contains(msg, "invalid password"):
		return ErrCatAuth
	}
	return ErrCatOther
}
