package settings

import (
    "time"
)

type ISettings interface {
    Validate() error

    GetRsaPrivateKey() string
    GetRsaPublicKey() string

    GetDataBaseServerName() string
    GetDataBaseName() string

    GetPortToServe() int
    GetTimeout() time.Duration
}

