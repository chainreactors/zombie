package service

import (
	"fmt"
	"strings"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/neutron/common"
	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/zombie/pkg"
)

func (r *Request) ExecuteWithResults(input *protocols.ScanContext, dynamicValues, previous map[string]interface{}, callback protocols.OutputEventCallback) error {
	sessionRaw, ok := input.Payloads["_session"]
	if !ok {
		return fmt.Errorf("service protocol: no session in scan context")
	}
	session, ok := sessionRaw.(pkg.Session)
	if !ok {
		return fmt.Errorf("service protocol: invalid session type")
	}

	if dynamicValues == nil {
		dynamicValues = make(map[string]interface{})
	}

	cliVars, _ := input.Payloads["_service_cli_vars"].(map[string]interface{})
	cliPayloads, _ := input.Payloads["_service_cli_payloads"].(map[string]interface{})
	payloadIterator, err := r.payloadIterator(cliVars, cliPayloads)
	if err != nil {
		return err
	}

	runOnce := payloadIterator == nil
	for {
		var payloadValues map[string]interface{}
		if payloadIterator == nil {
			if !runOnce {
				break
			}
			runOnce = false
		} else {
			var hasNext bool
			payloadValues, hasNext = payloadIterator.Value()
			if !hasNext {
				break
			}
		}

		var allResponses strings.Builder
		dslMap := r.responseToDSLMap("", session.Service(), input.Input)
		for k, v := range input.GlobalVars {
			dslMap[k] = v
		}
		for k, v := range previous {
			dslMap[k] = v
		}
		for k, v := range dynamicValues {
			dslMap[k] = v
		}
		for k, v := range payloadValues {
			dslMap[k] = v
		}

		for _, op := range r.Ops {
			normalized := normalizeOp(op)
			evaluated := evaluateOp(normalized, dslMap)
			response, err := executeOp(session, evaluated)
			if err != nil {
				logs.Log.Debugf("[service] op failed on %s: %v", input.Input, err)
				continue
			}

			if allResponses.Len() > 0 {
				allResponses.WriteString("\n")
			}
			allResponses.WriteString(response)

			if op.Name != "" {
				dslMap[op.Name] = response
			}
		}

		dslMap["response"] = allResponses.String()

		event := protocols.CreateEvent(r, dslMap)
		if event.OperatorsResult != nil && len(payloadValues) > 0 {
			event.OperatorsResult.PayloadValues = payloadValues
		}
		callback(event)

		if r.StopAtFirstMatch && event.OperatorsResult != nil && event.OperatorsResult.Matched {
			break
		}
	}

	return nil
}

func (r *Request) payloadIterator(varOverrides, payloadOverrides map[string]interface{}) (*protocols.Iterator, error) {
	if len(r.Payloads) == 0 && len(payloadOverrides) == 0 {
		return nil, nil
	}
	payloads := make(map[string]interface{}, len(r.Payloads)+len(payloadOverrides))
	for k, v := range r.Payloads {
		payloads[k] = v
	}
	for k, v := range varOverrides {
		if _, ok := payloads[k]; ok {
			payloads[k] = v
		}
	}
	for k, v := range payloadOverrides {
		payloads[k] = v
	}
	attack := strings.ToLower(r.AttackType)
	if attack == "" {
		attack = "pitchfork"
	}
	attackType, ok := protocols.StringToType[attack]
	if !ok {
		return nil, fmt.Errorf("unsupported attack type %q", r.AttackType)
	}
	generator, err := protocols.NewGenerator(payloads, attackType)
	if err != nil {
		return nil, err
	}
	return generator.NewIterator(), nil
}

// normalizeOp maps legacy fields to the primary shell/db/kv/file/ldap fields.
func normalizeOp(op *Op) *Op {
	if op.Shell != "" || op.DB != "" || op.KV != "" || op.File != nil || op.LDAP != nil {
		return op
	}

	n := *op
	switch {
	case n.Exec != "":
		n.Shell = n.Exec
	case n.Query != "":
		n.DB = n.Query
	case n.Get != "":
		n.KV = "GET " + n.Get
	case n.Keys != "":
		n.KV = "KEYS " + n.Keys
	case n.Cmd != "":
		n.KV = n.Cmd
	case n.List != "":
		n.File = &FileOp{List: n.List}
	case n.Read != "":
		n.File = &FileOp{Read: n.Read}
	case n.Search != nil:
		n.LDAP = n.Search
	}
	return &n
}

func evaluateOp(op *Op, values map[string]interface{}) *Op {
	evaluated := *op
	evaluated.Shell = evaluateField(op.Shell, values)
	evaluated.DB = evaluateField(op.DB, values)
	evaluated.KV = evaluateField(op.KV, values)
	if op.File != nil {
		f := *op.File
		f.List = evaluateField(op.File.List, values)
		f.Read = evaluateField(op.File.Read, values)
		f.Write = evaluateField(op.File.Write, values)
		f.Data = evaluateField(op.File.Data, values)
		evaluated.File = &f
	}
	if op.LDAP != nil {
		l := *op.LDAP
		l.BaseDN = evaluateField(op.LDAP.BaseDN, values)
		l.Filter = evaluateField(op.LDAP.Filter, values)
		evaluated.LDAP = &l
	}
	return &evaluated
}

func evaluateField(field string, values map[string]interface{}) string {
	if field == "" {
		return field
	}
	field = replacePayloadMarkers(field, values)
	if strings.Contains(field, "{{") {
		result, err := common.Evaluate(field, values)
		if err != nil {
			return field
		}
		return result
	}
	return field
}

func replacePayloadMarkers(field string, values map[string]interface{}) string {
	if !strings.Contains(field, "§") {
		return field
	}
	for k, v := range values {
		field = strings.ReplaceAll(field, "§"+k+"§", common.ToString(v))
	}
	return field
}

func executeOp(session pkg.Session, op *Op) (string, error) {
	switch {
	case op.Shell != "":
		return execShell(session, op.Shell)
	case op.DB != "":
		return execDB(session, op.DB)
	case op.KV != "":
		return execKV(session, op.KV)
	case op.File != nil:
		return execFile(session, op.File)
	case op.LDAP != nil:
		return execLDAP(session, op.LDAP)
	default:
		return "", fmt.Errorf("no operation specified in op")
	}
}

func execShell(session pkg.Session, cmd string) (string, error) {
	sh, ok := session.(pkg.ShellSession)
	if !ok {
		return "", fmt.Errorf("session does not support shell")
	}
	out, err := sh.Exec(cmd)
	return string(out), err
}

func execDB(session pkg.Session, query string) (string, error) {
	sq, ok := session.(pkg.SQLSession)
	if !ok {
		return "", fmt.Errorf("session does not support db")
	}
	rows, err := sq.Query(query)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, row := range rows {
		b.WriteString(strings.Join(row, "\t"))
		b.WriteString("\n")
	}
	return b.String(), nil
}

func execKV(session pkg.Session, expr string) (string, error) {
	kv, ok := session.(pkg.KVSession)
	if !ok {
		return "", fmt.Errorf("session does not support kv")
	}

	parts, err := parseCommandFields(expr)
	if err != nil {
		return "", err
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("empty kv expression")
	}

	verb := strings.ToUpper(parts[0])
	arg := strings.Join(parts[1:], " ")

	switch verb {
	case "GET":
		val, err := kv.Get(arg)
		return string(val), err
	case "KEYS":
		if arg == "" {
			arg = "*"
		}
		keys, err := kv.Keys(arg)
		if err != nil {
			return "", err
		}
		return strings.Join(keys, "\n"), nil
	default:
		result, err := kv.Command(parts[0], parts[1:]...)
		if err != nil {
			return "", err
		}
		return formatCommandResult(result), nil
	}
}

func formatCommandResult(result interface{}) string {
	switch v := result.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case []string:
		return strings.Join(v, "\n")
	case []interface{}:
		items := make([]string, 0, len(v))
		for _, item := range v {
			items = append(items, formatCommandResult(item))
		}
		return strings.Join(items, "\n")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func parseCommandFields(expr string) ([]string, error) {
	var fields []string
	var b strings.Builder
	var quote rune
	var escaped bool
	var tokenStarted bool

	flush := func() {
		if !tokenStarted {
			return
		}
		fields = append(fields, b.String())
		b.Reset()
		tokenStarted = false
	}

	for _, r := range expr {
		if escaped {
			switch r {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			default:
				b.WriteRune(r)
			}
			escaped = false
			tokenStarted = true
			continue
		}

		if quote != 0 {
			switch r {
			case '\\':
				escaped = true
			case quote:
				quote = 0
			default:
				b.WriteRune(r)
			}
			tokenStarted = true
			continue
		}

		switch r {
		case '\'', '"':
			quote = r
			tokenStarted = true
		case ' ', '\t', '\n', '\r':
			flush()
		default:
			b.WriteRune(r)
			tokenStarted = true
		}
	}

	if escaped {
		b.WriteRune('\\')
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quoted string in kv expression")
	}
	flush()
	return fields, nil
}

func execFile(session pkg.Session, op *FileOp) (string, error) {
	fs, ok := session.(pkg.FileSession)
	if !ok {
		return "", fmt.Errorf("session does not support file")
	}
	switch {
	case op.List != "":
		entries, err := fs.List(op.List)
		if err != nil {
			return "", err
		}
		return strings.Join(entries, "\n"), nil
	case op.Read != "":
		data, err := fs.Read(op.Read)
		return string(data), err
	case op.Write != "":
		if err := fs.Write(op.Write, []byte(op.Data)); err != nil {
			return "", err
		}
		return "OK", nil
	default:
		return "", fmt.Errorf("file op: set list, read, or write")
	}
}

func execLDAP(session pkg.Session, op *LDAPOp) (string, error) {
	dir, ok := session.(pkg.DirectorySession)
	if !ok {
		return "", fmt.Errorf("session does not support ldap")
	}
	results, err := dir.Search(op.BaseDN, op.Filter, op.Attrs)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, entry := range results {
		for attr, vals := range entry {
			for _, v := range vals {
				b.WriteString(attr)
				b.WriteString(": ")
				b.WriteString(v)
				b.WriteString("\n")
			}
		}
		b.WriteString("\n")
	}
	return b.String(), nil
}
