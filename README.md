# lbcwallet

lbcwallet is a daemon, which provides lbry wallet functionality for a
single user.

Public and private keys are derived using the hierarchical
deterministic format described by
[BIP0032](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki).
Unencrypted private keys are not supported and are never written to disk.

lbcwallet uses the `m/44'/<coin type>'/<account>'/<branch>/<address index>`
HD path for all derived addresses, as described by
[BIP0044](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki).

Due to the sensitive nature of public data in a BIP0032 wallet,
lbcwallet provides the option of encrypting not just private keys, but
public data as well.  This is intended to thwart privacy risks where a
wallet file is compromised without exposing all current and future
addresses (public keys) managed by the wallet. While access to this
information would not allow an attacker to spend or steal coins, it
does mean they could track all transactions involving your addresses
and therefore know your exact balance.  In a future release, public data
encryption will extend to transactions as well.

Wallet clients can use one of two RPC servers:

  1. A legacy JSON-RPC server mostly compatible with Bitcoin Core

     The JSON-RPC server exists to ease the migration of wallet applications
     from Core, but complete compatibility is not guaranteed.  Some portions of
     the API (and especially accounts) have to work differently due to other
     design decisions (mostly due to BIP0044).  However, if you find a
     compatibility issue and feel that it could be reasonably supported, please
     report an issue.  This server is enabled by default.

  2. An experimental gRPC server

     The gRPC server uses a new API built for lbcwallet, but the API is not
     stabilized and the server is feature gated behind a config option
     (`--experimentalrpclisten`).  If you don't mind applications breaking due
     to API changes, don't want to deal with issues of the legacy API, or need
     notifications for changes to the wallet, this is the RPC server to use.
     The gRPC server is documented [here](./rpc/documentation/README.md).

## Security

We take security seriously. Please contact [security](mailto:security@lbry.com) regarding any security issues.
Our PGP key is [here](https://lbry.com/faq/pgp-key) if you need it.

## Requirements

- [Go](http://golang.org) 1.16 or newer.

- `lbcwallet` is not an SPV client and requires connecting to a local or remote
  `lbcd` for asynchronous blockchain queries and notifications over websockets.

  Full installation instructions can be found [here](https://github.com/lbryio/lbcd).

## To Build lbcwallet, lbcd, and lbcctl from Source

Install Go according to its [installation instructions](http://golang.org/doc/install).

Build `lbcwallet`

``` sh
git clone https://github.com/lbryio/lbcwallet
cd lbcwallet
go build .
```

To make the quick start guide self-contained, here's how we can build the `lbcd` and `lbcctl`

``` sh
git clone https://github.com/lbryio/lbcd
cd lbcd

# build lbcd
go build .

# build lbcctl
go build ./cmd/lbcctl
```

## Getting Started

The first time running the `lbcwallet` we need to create a new wallet.

``` sh
./lbcwallet --create
```

Start a local instance of `lbcd` and have the `lbcwallet` connecting to it.

``` sh
# Start a lbcd with its RPC credentials
./lbcd --txindex --rpcuser=lbcduser --rpcpass=lbcdpass

# Start a lbcwallet with its RPC credentials along with the lbcd's RPC credentials
# The default lbcd instance to conect to is already localhost:9245 so we don't need to specify it explicitly here.
./lbcwallet --username=rpcuser --password=rpcpass --lbcdusername=lbcduser --lbcdpassword=lbcdpass # --rpcconnect=localhost:9245

#
#             rpcuser/rpcpass                lbcduser/lbcdpass
# lbcctl  <-------------------> lbcwallet <--------------------> lbcd
#             RPC port 9244                   RPC port 9245
#
```

If the `lbcd` and `lbcwallet` use the same RPC credentials, we can skip the `--lbcdusername` and `--lbcdpassword`

``` sh
./lbcd --txindex --rpcuser=rpcuser --rpcpass=rpcpass

./lbcwallet --username=rpcuser --password=rpcpass

#
#             rpcuser/rpcpass                rpcuser/rpcpass
# lbcctl  <-------------------> lbcwallet <--------------------> lbcd
#             RPC port 9244                   RPC port 9245
#
```

Note:

- `lbcd` and `lbcwallet` implements two disjoint sets of RPCs.
- `lbcd` serves RPC on port 9245 while `lbcwallet` on port 9244.
- `lbcwallet` can proxy non-wallet RPCs to its associated `lbcd`.

Examples of using `lbcctl` to interact with the setup via RPCs:

1. Calling non-wallet RPC directly on lbcd:

   ``` sh
   ./lbcctl --rpcuser=rpcuser --rpcpass=rpcpass getblockcount

   #
   # lbcctl  <-- getblockcount() --> lbcd
   #             RPC port 9245      (handled)
   #
   ```

2. Calling wallet RPC on lbcwallet (using `--wallet`)

   ``` sh
   ./lbcctl --rpcuser=rpcuser --rpcpass=rpcpass --wallet getbalance

   #
   # lbcctl  <-- getbalance() --> lbcwallet
   #             RPC port 9244    (handled)
   #
   ```

3. Calling non-wallet RPC on lbcwallet, which proxies it to lbcd:

   ``` sh
   ./lbcctl --rpcuser=rpcuser --rpcpass=rpcpass --wallet getblockcount

   #
   # lbcctl  <-- getblockcount() --> lbcwallet <-- getblockcount() --> lbcd
   #             RPC port 9244       (proxied)     RPC port 9245
   #
   ```

## Default Network and RPC Ports

| Instance      | mainnet | testet | regtest |
| ------------- | ------- | ------ | ------- |
| lbcd Network  | 9246    | 19246  | 29246   |
| lbcd RPC      | 9245    | 19245  | 29245   |
| lbcwallet RPC | 9244    | 19244  | 29244   |

Examples

``` sh
./lbcctl                    getblockcount # port  9245
./lbcctl --wallet           getblockcount # port  9244
./lbcctl --testnet          getblockcount # port 19245
./lbcctl --wallet --regtest getblockcount # port 29244
```

## Contributing

Contributions to this project are welcome, encouraged, and compensated.
The [integrated github issue tracker](https://github.com/lbryio/lbcwallet/issues)
is used for this project. All pull requests will be considered.

<!-- ## Release Verification
Please see our [documentation on the current build/verification
process](https://github.com/lbryio/lbcwallet/tree/master/release) for all our
releases for information on how to verify the integrity of published releases
using our reproducible build system.
-->

## License

lbcwallet is licensed under the liberal ISC License.
