package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chainreactors/zombie/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSession struct {
	service string
	client  *mongo.Client
	ctx     context.Context
}

func (s *mongoSession) Service() string  { return s.service }
func (s *mongoSession) Raw() interface{} { return s.client }

func (s *mongoSession) Close() error {
	if s.client != nil {
		return s.client.Disconnect(s.ctx)
	}
	return nil
}

func (s *mongoSession) Query(query string, args ...any) ([][]string, error) {
	query = strings.TrimSpace(query)

	parts := strings.SplitN(query, " ", 2)
	cmd := parts[0]

	switch strings.ToLower(cmd) {
	case "show":
		if len(parts) > 1 && strings.HasPrefix(strings.ToLower(parts[1]), "db") {
			dbs, err := s.client.ListDatabaseNames(s.ctx, bson.D{})
			if err != nil {
				return nil, err
			}
			rows := [][]string{{"database"}}
			for _, db := range dbs {
				rows = append(rows, []string{db})
			}
			return rows, nil
		}
		if len(parts) > 1 && strings.HasPrefix(strings.ToLower(parts[1]), "collection") {
			dbAndRest := strings.TrimPrefix(strings.ToLower(parts[1]), "collections ")
			colls, err := s.client.Database(dbAndRest).ListCollectionNames(s.ctx, bson.D{})
			if err != nil {
				return nil, err
			}
			rows := [][]string{{"collection"}}
			for _, c := range colls {
				rows = append(rows, []string{c})
			}
			return rows, nil
		}
	}

	cmdDoc := bson.D{{Key: cmd, Value: 1}}
	if len(parts) > 1 {
		cmdDoc = append(cmdDoc, bson.E{Key: "arg", Value: parts[1]})
	}
	result := s.client.Database("admin").RunCommand(s.ctx, cmdDoc)
	if result.Err() != nil {
		return nil, result.Err()
	}
	raw, err := result.DecodeBytes()
	if err != nil {
		return nil, err
	}
	return [][]string{{"result"}, {raw.String()}}, nil
}

func init() {
	pkg.RegisterPlugin("mongo", &MongoPlugin{})
	pkg.Services.Register(&pkg.Service{Name: "mongo", DefaultPort: "27017", Alias: []string{"mongodb"}, Source: pkg.PluginSource})
}

type MongoPlugin struct{}

func (p *MongoPlugin) Name() string { return "mongo" }

func (p *MongoPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	var url string
	if task.Password == "" {
		url = fmt.Sprintf("mongodb://%v:%v", task.IP, task.Port)
	} else {
		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", task.Username, task.Password, task.IP, task.Port)
	}
	// SetServerSelectionTimeout 限制 server selection / 命令等待,否则只 SetConnectTimeout
	// 在过滤端口上仍会按驱动默认 30s 阻塞,超出 task.Timeout。
	timeout := time.Duration(task.Timeout) * time.Second
	clientOptions := options.Client().ApplyURI(url).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout)

	client, err := mongo.Connect(task.Context, clientOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(task.Context, nil); err != nil {
		client.Disconnect(task.Context)
		return nil, err
	}
	return &mongoSession{service: task.Service, client: client, ctx: task.Context}, nil
}

// Unauth 以无凭据连接,并执行需要鉴权的 listDatabases 来证明“未授权可访问”。
// 只用 Ping 不够:mongo 的 ping 命令在开启鉴权的实例上同样放行,会把需鉴权实例
// 误报为未授权;listDatabases 在开启鉴权时返回 Unauthorized 错误,从而正确区分
// 真正的无认证实例(返回会话=命中)与需鉴权实例(返回错误=未命中)。
func (p *MongoPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	url := fmt.Sprintf("mongodb://%v:%v", task.IP, task.Port)
	timeout := time.Duration(task.Timeout) * time.Second
	clientOptions := options.Client().ApplyURI(url).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout)

	client, err := mongo.Connect(task.Context, clientOptions)
	if err != nil {
		return nil, err
	}
	if _, err := client.ListDatabaseNames(task.Context, bson.D{}); err != nil {
		client.Disconnect(task.Context)
		return nil, err
	}
	return &mongoSession{service: task.Service, client: client, ctx: task.Context}, nil
}
