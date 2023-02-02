/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	opts "xtravisions.com/xmgo/options"
)

// Credential mongodb 鉴权
//
//	参见 https://docs.mongodb.com/drivers/go/current/fundamentals/auth/
type Credential struct {
	AuthMechanism string `json:"authMechanism"` // 鉴权加密机制
	AuthSource    string `json:"authSource"`    // 鉴权数据库
	Username      string `json:"username"`      // 用户名
	Password      string `json:"password"`      // 密码
	PasswordSet   bool   `json:"passwordSet"`   // 使用 GSSAPI 鉴权加密时设置为 `true`
}

// ReadPref mongodb 只读操作服务器选择策略
type ReadPref struct {
	// 允许服务器被认为有资格选择的最长时间
	MaxStalenessMS int64 `json:"maxStalenessMS"`
	// 读取操作偏好设置
	// 	默认为 PrimaryMode
	Mode readpref.Mode `json:"mode"`
}

// Config mongodb 连接配置
type Config struct {
	// mongodb 连接地址
	// 	参见 https://docs.mongodb.com/manual/reference/connection-string/
	Uri string `json:"uri"`
	// mongodb 鉴权
	Auth *Credential `json:"auth"`
	// 用于创建连接时的超时设置
	//	设置为 0 意味不使用超时设置
	//	默认为 30 秒
	ConnectTimeoutMS *int64 `json:"connectTimeoutMS"`
	// 连接池最大值
	//	如果设置为 0 则使用 math.MaxInt64
	//	默认为 100
	MaxPoolSize *uint64 `json:"maxPoolSize"`
	// 连接池最小值，即初始化连接时创建的默认连接数量
	//	默认为 0
	MinPoolSize *uint64 `json:"minPoolSize"`
	// 数据读写操作等待超时时间
	//	默认为 300 秒
	SocketTimeoutMS *int64 `json:"socketTimeoutMS"`
	// 只读操作服务器选择策略
	ReadPreference *ReadPref `bson:"readPreference"`
}

// Client mongodb 连接
type Client struct {
	client   *mongo.Client
	registry *bsoncodec.Registry
	conf     Config
}

// NewClient 创建 mongodb 连接
//
//	@param conf 连接配置
//	@param opts 源生连接参数
func NewClient(conf *Config, o ...opts.ClientOptions) (cli *Client) {
	opt, err := newConnectOpts(conf, o...)
	if err != nil {
		fmt.Println("配置 MongoDB 连接参数失败")
		return nil
	}

	client, err := client(context.Background(), opt)
	if err != nil {
		fmt.Println("创建 MongoDB 连接失败")
		return
	}

	cli = &Client{
		client:   client,
		conf:     *conf,
		registry: opt.Registry,
	}

	if actions, ok := onConnected[conf.Uri]; ok {
		for _, cb := range actions {
			if err := cb.Fn(cli); err != nil {
				fmt.Println("Mongo 执行 OnConnect 钩子失败")
			}
		}
	}

	return
}

// Close 关闭 mongodb 连接
func (c *Client) Close() error {
	err := c.client.Disconnect(context.TODO())
	return err
}

// Database 连接到指定名称的数据库
//
//	@param name 数据库名称
//	@param options 数据库连接参数
func (c *Client) Database(name string, o ...*opts.DatabaseOptions) *Database {
	opt := options.Database()
	if len(o) > 0 {
		if o[0].DatabaseOptions != nil {
			opt = o[0].DatabaseOptions
		}
	}

	database := &Database{database: c.client.Database(name, opt), registry: c.registry}

	if cli, ok := onOpened[c.conf.Uri]; ok {
		if actions, ok := cli[name]; ok {
			for _, cb := range actions {
				if err := cb.Fn(database); err != nil {
					fmt.Println("Mongo 执行 OnOpen 钩子失败")
				}
			}
		}
	}

	return database
}

// Ping 确认连接是否可用
//
//	@param timeout 超时时间
func (c *Client) Ping(timeout int64) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	defer cancel()

	err = c.client.Ping(ctx, readpref.Primary())

	return
}

func (c *Client) Session(opt ...*opts.SessionOptions) (*Session, error) {
	sessionOpts := options.Session()
	if len(opt) > 0 && opt[0].SessionOptions != nil {
		sessionOpts = opt[0].SessionOptions
	}
	s, err := c.client.StartSession(sessionOpts)
	return &Session{session: s}, err
}

func (c *Client) DoTransaction(callback func(sessCtx context.Context) (interface{}, error), opts ...*opts.TransactionOptions) (interface{}, error) {
	return c.DoTransactionWithCtx(context.TODO(), callback, opts...)
}

func (c *Client) DoTransactionWithCtx(ctx context.Context, callback func(sessCtx context.Context) (interface{}, error), opts ...*opts.TransactionOptions) (interface{}, error) {
	if !c.transactionAllowed() {
		return nil, ErrTransactionNotSupported
	}

	s, err := c.Session()
	if err != nil {
		return nil, err
	}

	defer s.EndSession(ctx)
	return s.StartTransaction(ctx, callback, opts...)
}

func (c *Client) ServerVersion() string {
	var buildInfo bson.Raw
	err := c.client.Database("admin").RunCommand(
		context.Background(),
		bson.D{{"buildInfo", 1}},
	).Decode(&buildInfo)

	if err != nil {
		fmt.Println("尝试执行获取 mongodb 版本信息时出错", err)
		return ""
	}

	v, err := buildInfo.LookupErr("version")
	if err != nil {
		fmt.Println("获取 mongodb 版本信息出错", err)
		return ""
	}

	return v.StringValue()
}

func client(ctx context.Context, opt *options.ClientOptions) (client *mongo.Client, err error) {
	client, err = mongo.Connect(ctx, opt)
	if err != nil {
		fmt.Println("MongoDB Connect 失败", err)
		return
	}

	// half of default connect timeout
	pCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err = client.Ping(pCtx, readpref.Primary()); err != nil {
		fmt.Println("MongoDB Ping 失败", err)
		return
	}

	return
}

func newConnectOpts(conf *Config, o ...opts.ClientOptions) (*options.ClientOptions, error) {
	option := options.Client()
	for _, apply := range o {
		option = options.MergeClientOptions(apply.ClientOptions)
	}
	if conf.ConnectTimeoutMS != nil {
		timeoutDur := time.Duration(*conf.ConnectTimeoutMS) * time.Millisecond
		option.SetConnectTimeout(timeoutDur)

	}
	if conf.SocketTimeoutMS != nil {
		timeoutDur := time.Duration(*conf.SocketTimeoutMS) * time.Millisecond
		option.SetSocketTimeout(timeoutDur)
	} else {
		option.SetSocketTimeout(300 * time.Second)
	}
	if conf.MaxPoolSize != nil {
		option.SetMaxPoolSize(*conf.MaxPoolSize)
	}
	if conf.MinPoolSize != nil {
		option.SetMinPoolSize(*conf.MinPoolSize)
	}
	if conf.ReadPreference != nil {
		readPreference, err := newReadPref(*conf.ReadPreference)
		if err != nil {
			return nil, err
		}
		option.SetReadPreference(readPreference)
	}
	if conf.Auth != nil {
		auth, err := newAuth(*conf.Auth)
		if err != nil {
			return nil, err
		}
		option.SetAuth(auth)
	}
	option.ApplyURI(conf.Uri)

	return option, nil
}

func newAuth(auth Credential) (credential options.Credential, err error) {
	if auth.AuthMechanism != "" {
		credential.AuthMechanism = auth.AuthMechanism
	}
	if auth.AuthSource != "" {
		credential.AuthSource = auth.AuthSource
	}
	if auth.Username != "" {
		// Validate and process the username.
		if strings.Contains(auth.Username, "/") {
			err = ErrNotSupportedUsername
			return
		}
		credential.Username, err = url.QueryUnescape(auth.Username)
		if err != nil {
			err = ErrNotSupportedUsername
			return
		}
	}
	credential.PasswordSet = auth.PasswordSet
	if auth.Password != "" {
		if strings.Contains(auth.Password, ":") {
			err = ErrNotSupportedPassword
			return
		}
		if strings.Contains(auth.Password, "/") {
			err = ErrNotSupportedPassword
			return
		}
		credential.Password, err = url.QueryUnescape(auth.Password)
		if err != nil {
			err = ErrNotSupportedPassword
			return
		}
		credential.Password = auth.Password
	}
	return
}

func newReadPref(pref ReadPref) (*readpref.ReadPref, error) {
	readPrefOpts := make([]readpref.Option, 0, 1)
	if pref.MaxStalenessMS != 0 {
		readPrefOpts = append(readPrefOpts, readpref.WithMaxStaleness(time.Duration(pref.MaxStalenessMS)*time.Millisecond))
	}
	mode := readpref.PrimaryMode
	if pref.Mode != 0 {
		mode = pref.Mode
	}
	readPreference, err := readpref.New(mode, readPrefOpts...)
	return readPreference, err
}

func (c *Client) transactionAllowed() bool {
	vr, err := compareVersions("4.0", c.ServerVersion())
	if err != nil {
		return false
	}
	if vr > 0 {
		fmt.Println("transaction is not supported because mongo server version is below 4.0")
		return false
	}

	return true
}
