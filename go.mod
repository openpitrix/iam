module openpitrix.io/iam

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bitly/go-simplejson v0.5.0
	github.com/fatih/structs v1.1.0
	github.com/golang/protobuf v1.2.0
	github.com/google/gops v0.3.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/gorm v1.9.2
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/pkg/errors v0.8.1
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
	golang.org/x/net v0.7.0
	google.golang.org/genproto v0.0.0-20190215211957-bd968387e4aa
	google.golang.org/grpc v1.18.0
	gopkg.in/yaml.v2 v2.2.2
	kubesphere.io/im v0.0.0-20190216065048-6ad0bde84b60
	openpitrix.io/logger v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/jinzhu/now v1.0.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/appengine v1.1.0 // indirect
)

replace kubesphere.io/im v0.0.0-20190216065048-6ad0bde84b60 => github.com/kubesphere/im v0.0.0-20190216065048-6ad0bde84b60
