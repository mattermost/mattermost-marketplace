module github.com/mattermost/mattermost-marketplace

go 1.12

require (
	github.com/akrylysov/algnhsa v0.0.0-20190319020909-05b3d192e9a7
	github.com/aws/aws-lambda-go v1.13.0 // indirect
	github.com/go-ldap/ldap v3.0.3+incompatible // indirect
	github.com/google/go-github/v28 v28.0.0
	github.com/gorilla/mux v1.7.3
	github.com/mattermost/mattermost-server v5.12.3+incompatible
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.3.0
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/oauth2 v0.0.0-20190319182350-c85d3e98c914
	golang.org/x/sys v0.0.0-20190902133755-9109b7679e13 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
