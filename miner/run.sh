# go run main.go
go build -o ./bin/miner ./main.go
# sed -i 's/\r//' run.sh
trap "onCtrlC" SIGINT
function onCtrlC() {
    ps | grep ./bin/miner | grep -v grep | awk '{print $1}' | xargs kill -9
    echo "已关闭所有miner节点服务"
    exit
}

./bin/miner --privateKey 8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6 &
./bin/miner --privateKey 469ef6e06a15d66135732ffde307a63573529150d2e3cc1399f0d21285fba017 &
./bin/miner --privateKey a29e2136e7f5b6de2d8205195a819bd2bda3a38d6e5ebb07ff3ee80b20dbd486 &
./bin/miner --privateKey 98cf0979e72aabae9e192fae13f46a135c5fbad0ac719b979d007b27a7b85fbf &
./bin/miner --privateKey b7c03347692d8632cdc11aae9d458db1b13783477560c2f51fed424ac2a469ad &

while true; do
    sleep 1
done