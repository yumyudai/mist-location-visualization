package locapiserver

// Config defines the configuration structure for the location API server
type Config struct {
	Mist struct {
		Endpoint        string    `mapstructure:"endpoint"`
		Apikey          string    `mapstructure:"apikey"`
		LocationTimeout int       `mapstructure:"location_timeout"`
		RefreshTime     int       `mapstructure:"refresh_time"`
		Secret          string    `mapstructure:"secret"`
		Debug           bool      `mapstructure:"debug"`
	}                             `mapstructure:"mist"`
	Db struct {
		Driver  string `mapstructure:"driver"`
		Debug   bool   `mapstructure:"debug"`
		Mysql   struct {
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Host     string `mapstructure:"host"`
			Database string `mapstructure:"database"`
		} `mapstructure:"mysql"`
	} `mapstructure:"db"`
	Http struct {
		ServerName string `mapstructure:"server_name"`
		Listen     string `mapstructure:"listen"`
		BasicAuth  bool   `mapstructure:"basic_auth"`
		Debug      bool   `mapstructure:"debug"`
		Users      []struct {
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
		} `mapstructure:"users"`
	} `mapstructure:"http"`
}
