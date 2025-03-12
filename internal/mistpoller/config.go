package mistpoller

type Config struct {
	Db struct {
		Driver			string	  `mapstructure:"driver"`
		Debug			bool	  `mapstructure:"debug"`
		Mysql			struct {
			User		string	  `mapstructure:"user"`
			Password	string	  `mapstructure:"password"`
			Host		string	  `mapstructure:"host"`
			Database	string	  `mapstructure:"database"`
		}                                 `mapstructure:"mysql"`
	}                                         `mapstructure:"db"`
	Mist struct {
		Endpoint		string	  `mapstructure:"endpoint"`
		Apikey			string	  `mapstructure:"apikey"`
		Debug			bool	  `mapstructure:"debug"`
	}                                         `mapstructure:"mist"`
	Datasource []struct {
		Uri			string	  `mapstructure:"uri"`
		Datalayout		string	  `mapstructure:"data_layout"`
		Interval		int       `mapstructure:"interval"`
	}                                         `mapstructure:"datasource"`
}
