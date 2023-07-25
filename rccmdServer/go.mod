module github.com/jmbenlloch/rccmd/rccmdServer

go 1.18

require (
	github.com/kardianos/service v1.2.2
	github.com/magefile/mage v1.15.0
	github.com/sirupsen/logrus v1.9.0
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace github.com/jmbenlloch/rccmd/rccmdServer/pkg => ./pkg
