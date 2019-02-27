module openpitrix.io/iam

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bitly/go-simplejson v0.5.0
	github.com/fatih/structs v1.1.0
	github.com/golang/protobuf v1.2.0
	github.com/google/gops v0.3.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/now v1.0.0 // indirect
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/pkg/errors v0.8.1
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	google.golang.org/genproto v0.0.0-20190215211957-bd968387e4aa
	google.golang.org/grpc v1.18.0
	gopkg.in/yaml.v2 v2.2.2
	kubesphere.io/im v0.0.0-20190216065048-6ad0bde84b60
	openpitrix.io/logger v0.1.0
)

replace kubesphere.io/im v0.0.0-20190216065048-6ad0bde84b60 => github.com/kubesphere/im v0.0.0-20190216065048-6ad0bde84b60
