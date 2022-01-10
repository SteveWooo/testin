# go run main.go
go build -o ./bin/miner ./main.go
# sed -i 's/\r//' run.sh
trap "onCtrlC" SIGINT
function onCtrlC() {
    ps -ef | grep ./bin/miner | grep -v grep | awk '{print $2}' | xargs kill -9
    ps -ef | grep ../fe/service/main.go | grep -v grep | awk '{print $2}' | xargs kill -9
    echo "已关闭所有miner节点服务"
    exit
}

go run ../fe/service/main.go & 
sleep 5
# 激活 worldStatue
curl -H "Content-Type: application/json" -X POST -d '{"Params":{"MC_Call":"RegisterHacker","Name":"1","Resume":"1","Qualification":"1","Ts":"1641787886002","From":"0402e7ebab98d26359b235dd0405c6fb86","Hash":"bf423d096b388e1d35f7c4e00fd7e0f50cbf7423cf2ca50b48bb9a1cd76a7635","Signature":"1ba6857fb24dea7b74352332c39e926d848d378942cb88f31eac4c49698de2d2ce61970cd7a075d2bfa16733824f8adedec83d143b8159e97aa4f736f27de878ef"}}' http://localhost:10001/api/proxy


# ./bin/miner --privateKey 8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6 &
# ./bin/miner --privateKey 469ef6e06a15d66135732ffde307a63573529150d2e3cc1399f0d21285fba017 &
# ./bin/miner --privateKey a29e2136e7f5b6de2d8205195a819bd2bda3a38d6e5ebb07ff3ee80b20dbd486 &
./bin/miner --privateKey 98cf0979e72aabae9e192fae13f46a135c5fbad0ac719b979d007b27a7b85fbf --consensusMode PBFT &
# ./bin/miner --privateKey b7c03347692d8632cdc11aae9d458db1b13783477560c2f51fed424ac2a469ad &

# ./bin/miner --privateKey 797d861b693d1292e0af37861aa997afe96e1118c8c98f3f60155e797924698d &
# ./bin/miner --privateKey 6be56c87f8ac72722dbb8af867d1004f0b4a29af39fdcede44d43a7893d89519 &
# ./bin/miner --privateKey 149a34fd946f645d5f7dd4ba47cb1988e9047636e6236e7d0473eb3064fa41e3 &
# ./bin/miner --privateKey 7f171ab6fa3c909b1bcd7e9ff01decfd731f90ddb4846fadc56319be0218335b &
# ./bin/miner --privateKey cd86fae64530f60b75ec9f29db1464ecab5ea8f109765d24f755b17da9995328 &
# ./bin/miner --privateKey 6e880ebcee430fc1195819212a20f6a5f2b7ace77ec8b7262db6a89058b87261 &
# ./bin/miner --privateKey 6291917be72c3300677f6b52226cdd6483bf1eebd365c89441aac10e3dcf7ea2 &
# ./bin/miner --privateKey f12ce2ac3927ac1fba15e6db167aaf94dbbb4d6f242d55bc0d67b6bf7b895f43 &
# ./bin/miner --privateKey 1a9b01bf6c47dc0f5d3774137cca8908af0278b0ac545ce3652d495028de4b98 &
# ./bin/miner --privateKey ca6bda359d321e49ba942ae2df07d962294d8dae402b25415630f306ba06bc59 &
# ./bin/miner --privateKey 88d679ac2bb7ce86e30fbb3e19975642f1314d9353b3c0c8e90dcb3d0cfa2f61 &
# ./bin/miner --privateKey 31dea1c146f9260b20cc086c36c8377ff53f27aff96dd486519f5b1c52225df1 &
# ./bin/miner --privateKey 9f9c6e6ba61a07fe8a7046aa556a1cc98c78990db5540c1a4b1d3f1a067761ba &
# ./bin/miner --privateKey 759987048840864177c6e5ee1ed791e51460647f9db411498a5a100cc7cc227f &
# ./bin/miner --privateKey 7eaea8abe7c8bf5fa78798bce9d053c3eae27b75a42ae571f1275b89be6b388f &

while true; do
    sleep 1
done