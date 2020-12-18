package system

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"time"
)

type config struct {
	/** 数据源 */
	Databases []*DB

	/** 访问配置 */
	Server *serverConfig

	/** tokenConfig 配置 */
	Token *tokenConfig

	/** gin 配置 */
	Gin *ginConfig

	/** wechat 配置*/
	Wechat *wechatConfig
}

type wechatConfig struct {
	PayApiClientCert string `yaml:"wechat-pay-api-client-cert"`
}

// 生成token的配置
type tokenConfig struct {
	/** 发行人 */
	Issuer string `yaml:"issuer"`

	/** 签名秘钥 */
	SignKey string `yaml:"sign-key"`

	/** 签名生效时间  单位秒（配置负数表示立即生效，正数表示延迟生效*/
	ActiveTime int64 `yaml:"active-time"`

	/** token有效期 单位小时 */
	ExpiredTime int64 `yaml:"expired-time"`
}

type ginConfig struct {
	/** 运行模式 */
	//RunMode string `yaml:"run-mode"`

	/** 是否开启请求方式检测 */
	HandleMethodNotAllowed bool `yaml:"handle-method-not-allowed"`

	/** 设置请求占用最大空间 */
	MaxMultipartMemory int64 `yaml:"max-multipart-memory"`
}

type serverConfig struct {
	/** 访问地址 */
	Addr string `yaml:"addr"`

	/** 读取超时时间 */
	ReadTimeout time.Duration `yaml:"read-timeout"`

	/** 写超时时间 */
	WriteTimeout time.Duration `yaml:"write-timeout"`

	/** 最长闲置时间 */
	IdleTimeout time.Duration `yaml:"idle-timeout"`

	/** 请求头最大字节数 */
	MaxHeaderBytes int `yaml:"max-header-bytes"`
}

type DB struct {
	Driver string
	Nodes  datasource
}

// 数据源
type datasource struct {
	Database string `yaml:"database"`
	/** 驱动 */
	//Driver string `yaml:"driver"`

	/** 联接ip*/
	Host string `yaml:"host"`

	/**端口号*/
	Port string `yaml:"port"`

	/** 连接地址 */
	Url string `yaml:"url"`

	/** 用户名 */
	Username string `yaml:"username"`

	/** 密码 */
	Password string `yaml:"password"`

	/** 最大打开连接数 */
	MaxOpenConns int `yaml:"max-open-conns"`

	/** 最大闲置连接数 */
	MaxIdleConns int `yaml:"max-idle-conns"`

	/** 是否开启打印sql的日志 */
	ShowSql bool `yaml:"show-sql"`

	/** 自动创建表的时候 禁止表名为复数形式 */
	SingularTable bool `yaml:"singular-table"`
}

var configuration *config

func LoadDatasourceConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return err
	}
	return err
}

// 获取数据库连接配置
func GetDatasource() (datasource []*DB) {
	return configuration.Databases
}

func LoadServerConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return err
	}
	return err
}

// 获取项目启动配置
func GetServerConfig() (serverConfig *serverConfig) {
	return configuration.Server
}

// 获取 gin 配置
func GetGinConfig() (ginConfig *ginConfig) {
	return configuration.Gin
}

// 获取 wechat 配置
func GetWechatConfig() (wechatConfig *wechatConfig) {
	return configuration.Wechat
}

func LoadTokenConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return err
	}
	return err
}

// 获取Token配置
func GetTokenConfig() (tokenConfig *tokenConfig) {
	return configuration.Token
}

// 如果配置文件写成一个文件，调下面的方法
func LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return err
	}
	return err
}

// 获取配置文件
func GetConfig() *config {
	return configuration
}
