//go:build !emptytemplates
// +build !emptytemplates

package pkg

import (
	_ "embed"

	"github.com/chainreactors/utils/encode"
)

//go:embed data/zombie_default.bin
var zombieDefaultData []byte

//go:embed data/zombie_common.bin
var zombieCommonData []byte

//go:embed data/zombie_rule.bin
var zombieRuleData []byte

//go:embed data/zombie_template.bin
var zombieTemplateData []byte

//go:embed data/zombie_service.bin
var zombieServiceData []byte

//go:embed data/zombie_loot.bin
var zombieLootData []byte

//go:embed data/zombie_audit.bin
var zombieAuditData []byte

//go:embed data/port.bin
var portData []byte

//go:embed data/socket.bin
var socketData []byte

//go:embed data/http.bin
var httpData []byte

func loadEmbeddedConfig(typ string) []byte {
	if typ == "zombie_default" {
		return encode.MustDeflateDeCompress(zombieDefaultData)
	} else if typ == "zombie_common" {
		return encode.MustDeflateDeCompress(zombieCommonData)
	} else if typ == "zombie_rule" {
		return encode.MustDeflateDeCompress(zombieRuleData)
	} else if typ == "zombie_template" {
		return encode.MustDeflateDeCompress(zombieTemplateData)
	} else if typ == "zombie_service" {
		return encode.MustDeflateDeCompress(zombieServiceData)
	} else if typ == "zombie_loot" {
		return encode.MustDeflateDeCompress(zombieLootData)
	} else if typ == "zombie_audit" {
		return encode.MustDeflateDeCompress(zombieAuditData)
	} else if typ == "port" {
		return encode.MustDeflateDeCompress(portData)
	} else if typ == "socket" {
		return encode.MustDeflateDeCompress(socketData)
	} else if typ == "http" {
		return encode.MustDeflateDeCompress(httpData)
	}
	return nil
}
