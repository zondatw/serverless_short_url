module github.com/zondatw/serverless_short_url/shorturlfunction

go 1.15

require (
	cloud.google.com/go/firestore v1.5.0
	cloud.google.com/go/pubsub v1.10.3
	cloud.google.com/go/storage v1.15.0 // indirect
	firebase.google.com/go v3.13.0+incompatible
	github.com/alicebob/miniredis/v2 v2.14.4
	github.com/gomodule/redigo v1.8.4
	github.com/linkedin/goavro v2.1.0+incompatible
	github.com/linkedin/goavro/v2 v2.10.0 // indirect
	golang.org/x/exp/errors v0.0.0-20210514180818-737f94c0881e
	google.golang.org/api v0.47.0
)
