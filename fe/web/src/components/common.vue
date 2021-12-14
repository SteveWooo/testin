<script>
import * as secp from "@noble/secp256k1"
import crypto from 'crypto-js'
const BASE_URL = "http://127.0.0.1:10001"
const IPFS_API_HOST = "127.0.0.1"
const IPFS_API_PORT = 5001
const IPFS_GATEWAY = "http://127.0.0.1:8081"

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

var api = {
    GetWorldStatus : function(){
      return new Promise((resolve, reject)=>{
          window.axios({
            method : "get",
            url : BASE_URL + "/api/world_status/get",
        }).then(res=>{
            if (res.data.Status != 2000) {
                alert("获取世界状态失败")
            }
            resolve(res.data)
        }, rejectErr=>{
          reject(rejectErr)
        })
      })
    },

    // 获取企业发布的任务列表
    GetEnterprisePublishedTasks : function(page, itemPerPage){
        return new Promise((resolve, reject)=>{
          var now = +new Date()
          var salt = "salt" + now
          var saltHsah = crypto.SHA256(salt).toString()
          
          secp256k1.Sign(saltHsah).then(signature=>{
            if(signature == undefined) {
                reject("签名失败")
                return 
            }
            window.axios({
                method : "get",
                url : BASE_URL + "/api/enterprise/get_my_task?page=" + page
                + "&item_per_page=" + itemPerPage
                + "&node_id=" + secp256k1.GetNodeID()
                + "&salt=" + saltHsah
                + "&signature=" + signature,
            }).then(res=>{
                if (res.data.Status != 2000) {
                    alert("获取世界状态失败")
                }
                resolve(res.data)
            }, rejectErr=>{
            reject(rejectErr)
            })
          })
      })
    },

    // 获取任务详情
    GetTaskDetail : function(taskID){
        return new Promise((resolve, reject)=>{
          var now = +new Date()
          var salt = "salt" + now
          var saltHsah = crypto.SHA256(salt).toString()
          
          secp256k1.Sign(saltHsah).then(signature=>{
            if(signature == undefined) {
                reject("签名失败")
                return 
            }
            window.axios({
                method : "get",
                url : BASE_URL + "/api/common/get_task_detail?task_id=" + taskID
                + "&node_id=" + secp256k1.GetNodeID()
                + "&salt=" + saltHsah
                + "&signature=" + signature,
            }).then(res=>{
                if (res.data.Status != 2000) {
                    alert(res.data.Status + ":" + res.data.Message)
                }
                resolve(res.data)
            }, rejectErr=>{
            reject(rejectErr)
            })
          })
      })
    },


    CallTrans : function (params){
        return new Promise((resolve, reject)=>{
            window.axios({
                method : "post",
                url : BASE_URL + "/api/proxy",
                headers : {
                    "Content-Type" : "Application/json"
                },
                data : JSON.stringify({
                    "Params" : params
                })
            }).then(res=>{
                if (res.data.Status != 2000) {
                    alert("调用提交失败：" + res.data.Message)
                }
                resolve(res.data)
            }, rejectErr=>{
                reject(rejectErr)
            })
        })
    }
}

export default {
    BASE_URL,
    IPFS_API_HOST,
    IPFS_API_PORT,
    IPFS_GATEWAY,
    ls,
    secp256k1,
    api,
}
</script>