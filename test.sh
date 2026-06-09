#!/usr/bin/env bash
BIN="$PWD/fireeth"
for BS in 1 10; do
  rm -rf /tmp/fhtest; mkdir -p /tmp/fhtest/one /tmp/fhtest/merged /tmp/fhtest/forked /tmp/fhtest/poller
  cat > /tmp/rn.yml <<YAML
start:
  args: [reader-node]
  flags:
    log-to-file: false
    common-one-block-store-url: file:///tmp/fhtest/one
    common-merged-blocks-store-url: file:///tmp/fhtest/merged
    common-forked-blocks-store-url: file:///tmp/fhtest/forked
    reader-node-grpc-listen-addr: :9111
    reader-node-one-block-suffix: test1
    reader-node-readiness-max-latency: 1200s
    reader-node-path: $BIN
    reader-node-arguments: "tools poller generic-evm http://megaeth-arch40.kan.eosn.io:8545 6600000 --data-dir=/tmp/fhtest/poller --max-block-fetch-duration=120s --parallel-workers=10 --block-fetch-batch-size=$BS"
YAML
  "$BIN" start -c /tmp/rn.yml >/tmp/rn.log 2>&1 &
  PID=$!
  sleep 45
  kill $PID 2>/dev/null; pkill -P $PID 2>/dev/null; pkill -f 'tools poller generic-evm' 2>/dev/null; wait $PID 2>/dev/null
  N=$(ls /tmp/fhtest/one 2>/dev/null | wc -l | tr -d ' ')
  echo "block-fetch-batch-size=$BS -> $N one-block files in 45s"
  sleep 2
done
