# ../../../ipfs/go-ipfs/cmd/ipfs/ipfs init
# 12D3KooWS8NvbeF2spMBcEnsNum7jq3iPD2gsEXT3ETnjsD7BaLk

trap "onCtrlC" SIGINT
function onCtrlC() {
    ../../../ipfs/go-ipfs/cmd/ipfs/ipfs shutdown
    echo "已关闭所有ipfs服务"
    exit
}

../../../ipfs/go-ipfs/cmd/ipfs/ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
../../../ipfs/go-ipfs/cmd/ipfs/ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
../../../ipfs/go-ipfs/cmd/ipfs/ipfs daemon

while true;do
    sleep 10
done 