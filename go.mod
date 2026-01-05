module github.com/toutaio/toutago

go 1.24.0

toolchain go1.24.11

require (
	github.com/adrg/frontmatter v0.2.0
	github.com/spf13/cobra v1.10.2
	github.com/toutaio/toutago-cosan-router v1.0.5
	github.com/toutaio/toutago-datamapper v1.0.8
	github.com/toutaio/toutago-fith-renderer v1.0.6
	github.com/toutaio/toutago-nasc-dependency-injector v1.0.9
	github.com/toutaio/toutago-ritual-grove v0.2.2
	github.com/toutaio/toutago-scela-bus v1.5.4
	github.com/toutaio/toutago-sil-migrator v1.0.5
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/text v0.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/toutaio/toutago-ritual-grove => ../toutago-ritual-grove
