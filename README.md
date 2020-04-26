# twelvenv

![Go](https://github.com/izbudki/twelvenv/workflows/Go/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/izbudki/twelvenv/branch/master/graph/badge.svg)](https://codecov.io/gh/izbudki/twelvenv)

## Usage example

```
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/izbudki/twelvenv"
)

type DBConfig struct {
	Host string `name:"DB_HOST" required:"true"`
	Port int    `name:"DB_PORT" required:"true"`
}

type Config struct {
	FooDatabase DBConfig      `env_prefix:"FOO_"`
	BarDatabase DBConfig      `env_prefix:"BAR_"`
	Addresses   []string      `name:"ADDRESSES" required:"true"`
	Interval    time.Duration `name:"INTERVAL" required:"true"`
	EnableLogs  string        `name:"ENABLE_LOGS"`
}

func main() {
	os.Setenv("FOO_DB_HOST", "localhost")
	os.Setenv("BAR_DB_HOST", "localhost")
	os.Setenv("FOO_DB_PORT", "5432")
	os.Setenv("BAR_DB_PORT", "5433")
	os.Setenv("ENABLE_LOGS", "true")
	os.Setenv("ADDRESSES", "localhost:8000,localhost:8001,localhost:8002")
	os.Setenv("INTERVAL", "15s")

	var cfg Config
	err := twelvenv.FromEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	configJSON, _ := json.MarshalIndent(cfg, "", "    ")
	fmt.Println(string(configJSON))
}
```

Result:
```
{
    "FooDatabase": {
        "Host": "localhost",
        "Port": 5432
    },
    "BarDatabase": {
        "Host": "localhost",
        "Port": 5433
    },
    "Addresses": [
        "localhost:8000",
        "localhost:8001",
        "localhost:8002"
    ],
    "Interval": 15000000000,
    "EnableLogs": "true"
}
```

## TODO

- [x] Prefixes for nested structs
- [ ] Simple validation
- [ ] Verbose error messages
- [ ] Clever tags analysis
