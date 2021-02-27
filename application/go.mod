module yiarce/application

go 1.15

require (
	github.com/gomodule/redigo v1.8.3
	yiarce/dragonnews v0.0.0
)

replace (
	yiarce/application => ../application
	yiarce/dragonnews => ../dragonnews
)
