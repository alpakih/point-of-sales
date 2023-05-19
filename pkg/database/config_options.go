package database

import (
	"strconv"
)

type ConfigOption func(*Config)

func ConfigDriverName(driverName string) ConfigOption {
	return func(cfg *Config) {
		cfg.Driver = driverName
		cfg.TemplateDsn = templateDsn[driverName]
	}
}

func ConfigHost(host string) ConfigOption {
	return func(cfg *Config) { cfg.Host = host }
}

func ConfigPort(port string) ConfigOption {
	return func(cfg *Config) { cfg.Port = port }
}

func ConfigUsername(username string) ConfigOption {
	return func(cfg *Config) { cfg.Username = username }
}

func ConfigPassword(password string) ConfigOption {
	return func(cfg *Config) { cfg.Password = password }
}

func ConfigDebugEnabled(enabled bool) ConfigOption {
	return func(cfg *Config) { cfg.Debug = enabled }
}

func ConfigMaxOpenConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxOpenConnection = value }
}

func ConfigMaxIdleConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxIdleConnection = value }
}

func ConfigMaxLifeTimeConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxLifeTimeConnection = value }
}

func ConfigMaxIdleTimeConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxIdleTimeConnection = value }
}

func ConfigNewrelicIntegration(enabled bool) ConfigOption {
	return func(cfg *Config) { cfg.NewrelicIntegration = enabled }
}

func ConfigFromEnvironment(dbConfigEnv map[string]string) ConfigOption {
	return configFromEnvironment(dbConfigEnv)
}

func configFromEnvironment(getEnv map[string]string) ConfigOption {

	return func(config *Config) {
		config.Driver = getEnv["driver"]
		config.TemplateDsn = templateDsn[getEnv["driver"]]
		config.Host = getEnv["host"]
		config.Port = getEnv["port"]
		config.Username = getEnv["username"]
		config.Password = getEnv["password"]
		config.Name = getEnv["name"]
		config.Options = getEnv["options"]
		if parse, err := strconv.ParseBool(getEnv["debug"]); err != nil {
			config.Debug = true
		} else {
			config.Debug = parse
		}
		if parse, err := strconv.Atoi(getEnv["maxopenconn"]); err == nil {
			config.MaxOpenConnection = parse
		}
		if parse, err := strconv.Atoi(getEnv["maxidleconn"]); err == nil {
			config.MaxIdleConnection = parse
		}
		if parse, err := strconv.Atoi(getEnv["maxlifetimeconn"]); err == nil {
			config.MaxLifeTimeConnection = parse
		}
		if parse, err := strconv.Atoi(getEnv["maxidletimeconn"]); err == nil {
			config.MaxIdleTimeConnection = parse
		}
		if parse, err := strconv.ParseBool(getEnv["newrelicintegration"]); err != nil {
			config.NewrelicIntegration = false
		} else {
			config.NewrelicIntegration = parse
		}
	}
}
