module dragonNews

go 1.21

toolchain go1.23.3

replace (
	yiarce/core => ./core
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
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-vgo/robotgo v0.110.1
	github.com/gorilla/websocket v1.5.3
	github.com/robotn/gohook v0.41.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	golang.org/x/sys v0.22.0
	google.golang.org/protobuf v1.35.2
	gopkg.in/yaml.v2 v2.4.0
	yiarce/core v0.0.0-00010101000000-000000000000
	yiarce/core/curl v0.0.0-00010101000000-000000000000
	yiarce/core/date v0.0.0-00010101000000-000000000000
	yiarce/core/dhttp v0.0.0-00010101000000-000000000000
	yiarce/core/encrypt v0.0.0-00010101000000-000000000000
	yiarce/core/file v0.0.0-00010101000000-000000000000
	yiarce/core/frame v0.0.0-00010101000000-000000000000
	yiarce/core/log v0.0.0-00010101000000-000000000000
	yiarce/core/timing v0.0.0-00010101000000-000000000000
	yiarce/core/yorm v0.0.0-00010101000000-000000000000
	yiarce/core/yorm/mysql v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gen2brain/shm v0.0.0-20230802011745-f2460f5984f7 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/kbinani/screenshot v0.0.0-20230812210009-b87d31814237 // indirect
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/otiai10/gosseract v2.2.1+incompatible // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/robotn/xgb v0.0.0-20190912153532-2cb92d044934 // indirect
	github.com/robotn/xgbutil v0.0.0-20190912154524-c861d6f87770 // indirect
	github.com/shirou/gopsutil/v3 v3.23.8 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/vcaesar/gops v0.30.2 // indirect
	github.com/vcaesar/imgo v0.40.0 // indirect
	github.com/vcaesar/keycode v0.10.1 // indirect
	github.com/vcaesar/tt v0.20.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	golang.org/x/image v0.12.0 // indirect
)
