module github.com/lbryio/lbcwallet

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/davecgh/go-spew v1.1.1
	github.com/golang/protobuf v1.5.2
	github.com/jessevdk/go-flags v1.5.0
	github.com/jrick/logrotate v1.0.0
	github.com/lbryio/lbcd v0.22.100-beta-rc7
	github.com/lbryio/lbcutil v1.0.202-rc3
	github.com/lightningnetwork/lnd/clock v1.0.1
	github.com/stretchr/testify v1.7.0
	go.etcd.io/bbolt v1.3.6
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	google.golang.org/grpc v1.40.0
)

require (
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/aead/siphash v1.0.1 // indirect
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cockroachdb/errors v1.8.6 // indirect
	github.com/cockroachdb/logtags v0.0.0-20211118104740-dabe8e521a4f // indirect
	github.com/cockroachdb/pebble v0.0.0-20211124004043-0dc90bc41e62 // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/cockroachdb/sentry-go v0.6.1-cockroachdb.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/kkdai/bstream v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	golang.org/x/exp v0.0.0-20211123021643-48cbe7f80d7c // indirect
	golang.org/x/sys v0.0.0-20211123173158-ef496fb156ab // indirect
	golang.org/x/term v0.0.0-20210916214954-140adaaadfaf // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210921142501-181ce0d877f6 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// The old version of ginko that's used in btcd imports an ancient version of
// gopkg.in/fsnotify.v1 that isn't go mod compatible. We fix that import error
// by replacing ginko (which is only a test library anyway) with a more recent
// version.
replace github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.14.2

go 1.17
