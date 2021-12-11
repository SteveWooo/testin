<script>
import * as secp from "@noble/secp256k1"

const baseUrl = "http://127.0.0.1:10001"
// const baseUrl = "/"

var secp256k1 = {
    getPrivateKeyFromString : function(privateKey){
        var buf = new Uint8Array(32)
        for(var i=0;i<32;i++) {
            buf[i] = parseInt(privateKey.substring(i*2, i*2+2), 16)
        }
        return window.Secp256k1.uint256(buf, 10)
    },
    // 从64位字符串中获取34位NodeID
    GetNodeIDFromPrivateKey : function(privateKey) {
        return secp.getPublicKey(privateKey).substring(0, 34)
    },

    GetNodeID : function(){
        var pk = ls.get("privateKey")
        if (pk == undefined) {
            return undefined
        }
        return secp256k1.GetNodeIDFromPrivateKey(pk)
    },
    Sign : async function(message) {
        var pk = ls.get("privateKey")
        if (pk == undefined) {
            return undefined
        }

        // var privateKey = secp256k1.getPrivateKeyFromString(pk)
        // var source = window.Secp256k1.uint256(message)
        // var signature = window.Secp256k1.ecsign(privateKey, source)
        // console.log(signature)

        var signature = await secp.sign(message, pk, {
            recovered : true,
            der : false,
            canonical : true
        })
        var resultSign = ""
        if (signature[1] == 1) {
            resultSign = "1c" + signature[0]
        }
        if (signature[1] == 0) {
            resultSign = "1b" + signature[0]
        }
        return resultSign
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