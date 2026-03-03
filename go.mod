module dragonNews

go 1.24.0

replace (
	yiarce/core => ./core
	yiarce/core/cache => ./core/cache
	yiarce/core/config => ./core/config
	yiarce/core/curl => ./core/curl
	yiarce/core/date => ./core/date
	yiarce/core/dhttp => ./core/dhttp
	yiarce/core/encrypt => ./core/encrypt
	yiarce/core/file => ./core/file
	yiarce/core/frame => ./core/frame
	yiarce/core/log => ./core/log
	yiarce/core/timing => ./core/timing
	yiarce/core/yorm => ./core/yorm
	yiarce/core/yorm/mysql => ./core/yorm/mysql
)

require (
	github.com/gorilla/websocket v1.5.3
	golang.org/x/sys v0.41.0
	yiarce/core v0.0.0-00010101000000-000000000000
	yiarce/core/cache v0.0.0-00010101000000-000000000000
	yiarce/core/config v0.0.0-00010101000000-000000000000
	yiarce/core/date v0.0.0-00010101000000-000000000000
	yiarce/core/dhttp v0.0.0-00010101000000-000000000000
	yiarce/core/encrypt v0.0.0-00010101000000-000000000000
	yiarce/core/file v0.0.0-00010101000000-000000000000
	yiarce/core/frame v0.0.0-00010101000000-000000000000
	yiarce/core/log v0.0.0-00010101000000-000000000000
	yiarce/core/timing v0.0.0-00010101000000-000000000000
	yiarce/core/yorm v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	yiarce/core/yorm/mysql v0.0.0-00010101000000-000000000000 // indirect
)
