<script>
const baseUrl = "http://127.0.0.1:10001"
// const baseUrl = "/"

var secp256k1 = {
    // 从64位字符串中获取34位NodeID
    GetNodeIDFromPrivateKey : function(privateKey) {
        if (privateKey.length != 64) {
            return undefined
        }
        var Secp256k1 = window.Secp256k1
        // 把字符串转换为 Uint8Array
        var buf = new Uint8Array(32)
        for(var i=0;i<32;i++) {
            buf[i] = parseInt(privateKey.substring(i*2, i*2+2), 16)
        }

        var nodeID = Secp256k1.generatePublicKeyFromPrivateKeyData(window.Secp256k1.uint256(buf, 10))
        return "04" + nodeID.x.substring(0, 32)
    }
}

var ls = {
    set : function(key, value){
        localStorage.setItem(key, value)
    },
    get : function(key) {
        return localStorage.getItem(key)
    },
    remove : function(key) {
        localStorage.removeItem(key)
    }
}

export default {
    baseUrl,
    ls,
    secp256k1,
}
</script>