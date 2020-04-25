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

type Config struct {
	Database struct {
		Host string `name:"DB_HOST" required:"true"`
		Port int    `name:"DB_PORT" required:"true"`
	}
	Addresses  []string      `name:"ADDRESSES" required:"true"`
	Interval   time.Duration `name:"INTERVAL" required:"true"`
	EnableLogs string        `name:"ENABLE_LOGS"`
}

func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
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
    "Database": {
        "Host": "localhost",
        "Port": 5432
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

- [ ] Prefixes for nested structs
- [ ] Simple validation
