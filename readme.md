# polaris

**polaris** is a blockchain built using Cosmos SDK and Tendermint.

To start IBC, you should start two blockchain networks. Both blockchains use the same code. Each blockchain has a unique
chain ID.

## Commands

# Get started

```
ignite chain serve 
ignite chain serve -c polaris.yml
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

# Remove existing Relayer

```
rm -rf ~/.ignite/relayer
```

# Config Relayer

```
ignite relayer configure -a \
  --source-rpc "http://0.0.0.0:26657" \
  --source-faucet "http://0.0.0.0:4500" \
  --source-port "blog" \
  --source-version "blog-1" \
  --source-gasprice "0.0000025stake" \
  --source-prefix "cosmos" \
  --source-gaslimit 300000 \
  --target-rpc "http://0.0.0.0:26659" \
  --target-faucet "http://0.0.0.0:4501" \
  --target-port "blog" \
  --target-version "blog-1" \
  --target-gasprice "0.0000025stake" \
  --target-prefix "cosmos" \
  --target-gaslimit 300000
```

# Start Relayer

```
ignite relayer connect
```

### Configure

The Polaris blockchain development can be configured with `polaris.yml`.
The Meris blockchain development can be configured with `meris.yml`.



