#!/bin/sh

set -e

KOMERCO_EDGE_BIN=./komerco-chain
CHAIN_CUSTOM_OPTIONS=$(tr "\n" " " << EOL
--block-gas-limit 10000000
--epoch-size 10
--chain-id 51001
--name komerco-chain-docker
--premine 0x228466F2C715CbEC05dEAbfAc040ce3619d7CF0B:0xD3C21BCECCEDA1000000
--premine 0xca48694ebcB2548dF5030372BE4dAad694ef174e:0xD3C21BCECCEDA1000000
EOL
)

case "$1" in

   "init")
      case "$2" in 
         "ibft")
         if [ -f "$GENESIS_PATH" ]; then
              echo "Secrets have already been generated."
         else
              echo "Generating secrets..."
              secrets=$("$KOMERCO_EDGE_BIN" secrets init --insecure --num 4 --data-dir /data/data- --json)
              echo "Secrets have been successfully generated"
              echo "Generating IBFT Genesis file..."
              cd /data && /komerco-chain/komerco-chain genesis $CHAIN_CUSTOM_OPTIONS \
                --dir genesis.json \
                --consensus ibft \
                --ibft-validators-prefix-path data- \
                --validator-set-size=4 \
                --bootnode "/dns4/node-1/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[0] | .node_id')" \
                --bootnode "/dns4/node-2/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[1] | .node_id')"
         fi
              ;;
          "komerbft")
              echo "Generating KomerBFT secrets..."
              secrets=$("$KOMERCO_EDGE_BIN" komerbft-secrets init --insecure --num 4 --data-dir /data/data- --json)
              echo "Secrets have been successfully generated"

              echo "Generating manifest..."
              "$KOMERCO_EDGE_BIN" manifest --path /data/manifest.json --validators-path /data --validators-prefix data-

              echo "Generating KomerBFT Genesis file..."
              "$KOMERCO_EDGE_BIN" genesis $CHAIN_CUSTOM_OPTIONS \
                --dir /data/genesis.json \
                --consensus komerbft \
                --manifest /data/manifest.json \
                --validator-set-size=4 \
                --bootnode "/dns4/node-1/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[0] | .node_id')" \
                --bootnode "/dns4/node-2/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[1] | .node_id')"
              ;;
      esac
      ;;

   *)
      echo "Executing komerco-chain..."
      exec "$KOMERCO_EDGE_BIN" "$@"
      ;;

esac
