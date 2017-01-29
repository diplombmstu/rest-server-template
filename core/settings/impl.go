package settings

import (
    "time"
    "gopkg.in/ini.v1"
)

const (
    KeyRsaPrivKey = "rsa_priv_key"
    KeyRsaPubKey  = "rsa_pub_key"

    KeyDbHost = "db_host"
    KeyDbName = "dm_name"

    KeyServePort = "serve_port"
    KeyTimeout   = "timeout"

    KeyLocalUsers = "user_repo"
)

// impls ISettings
type settings struct {
    cfg *ini.File
}

func NewSettings(configFile string) (ISettings, error) {
    cfg, err := ini.Load(configFile)
    if err != nil {
        return nil, err
    }

    result := &settings{cfg:cfg}

    return result, nil
}

func (settings *settings) Validate() error {
    return nil
}

func (settings *settings) GetRsaPrivateKey() string {
    return settings.cfg.Section("").Key(KeyRsaPrivKey).String()
}

func (settings *settings) GetRsaPublicKey() string {
    return settings.cfg.Section("").Key(KeyRsaPubKey).String()
}

func (settings *settings) GetDataBaseServerName() string {
    return settings.cfg.Section("").Key(KeyDbHost).String()
}

func (settings *settings) GetDataBaseName() string {
    return settings.cfg.Section("").Key(KeyDbName).String()
}

func (settings *settings) GetPortToServe() int {
    result, _ := settings.cfg.Section("").Key(KeyServePort).Int()
    return result
}

func (settings *settings) GetTimeout() time.Duration {
    result, _ := settings.cfg.Section("").Key(KeyTimeout).Int64()
    return time.Duration(result)
}
