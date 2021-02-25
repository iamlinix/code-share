package conf

type SSLConf struct {
	Enable   bool   `yaml:"enable"`
	KeyFile  string `yaml:"key"`
	CertFile string `yaml:"cert"`
}

type DatabaseConf struct {
	Backend  int               `yaml:"backend"`
	Driver   string            `yaml:"driver"`
	Database string            `yaml:"database"`
	Username string            `yaml:"username"`
	Password string            `yaml:"password"`
	Extras   map[string]string `yaml:"extras"`
	Pool     map[string]string `yaml:"pool"`
}

type FileSystemConf struct {
	Backend int                    `yaml:"backend"`
	BaseDir string                 `yaml:"base_dir"`
	Extras  map[string]interface{} `yaml:"extras"`
}

type EncryptConf struct {
	Code string `yaml:"code"`
}

type LogDetailConf struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type LogConf struct {
	App    LogDetailConf `yaml:"app"`
	Access LogDetailConf `yaml:"access"`
}

type WebAuthConf struct {
	Enable bool                   `yaml:"enable"`
	Method int                    `yaml:"method"`
	Extras map[string]interface{} `yaml:"extras"`
}

type WebThrottleConf struct {
	Enable    bool     `yaml:"enable"`
	MaxPerSec int      `yaml:"max_per_sec"`
	MaxBurst  int      `yaml:"max_burst"`
	Urls      []string `yaml:"urls"`
}

type WebConf struct {
	WriteTimeout int                 `yaml:"wto"`
	ReadTimeout  int                 `yaml:"rto"`
	SSL          SSLConf             `yaml:"ssl"`
	Auth         WebAuthConf         `yaml:"auth"`
	Throttle     WebThrottleConf     `yaml:"throttle"`
	Statics      []map[string]string `yaml:"static"`
	Serves       []map[string]string `yaml:"serve"`
	Bind         string              `yaml:"bind"`
	Port         int                 `yaml:"port"`
	Cors         bool                `yaml:"cors"`
}

type WxConf struct {
	AppId  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

type Conf struct {
	Log        LogConf        `yaml:"log"`
	Mode       string         `yaml:"mode"`
	Database   DatabaseConf   `yaml:"database"`
	FileSystem FileSystemConf `yaml:"fs"`
	Encrypt    EncryptConf    `yaml:"encrypt"`
	Web        WebConf        `yaml:"web"`
	Wx         WxConf         `yaml:"wx"`
}
