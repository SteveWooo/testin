const http = require("http")
const secp = require("../fe/web/node_modules/@noble/secp256k1")
const crypto = require("../fe/web/node_modules/crypto-js")

var config = {
    // PrivateKey : "8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad4",
    PrivateKey : "8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad3",

    TaskID : "69ed39489b89d104bd429de0d1202ddc060cc06ddc7d0f4f901676d218012329"
}
var utils = {
    Sign : async function(message, pk) {
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

function doPost(index){
    return new Promise(async function(resolve){
        var now = +new Date()
        var data = {
            "Params" : {
                "MC_Call" : "PublishReportByHacker", // 调用合约脚本
                From: "0402e7ebab98d26359b235dd0405c6fb86",
                ReportPath: "QmNnGUHwagNAzhBUwGvsfCMvuVKfL294eKPu1r2juXVN3J",
                // TaskID: "d66a43d5f4165769a277d0f01fbf2af5faf9999ce25ce04316a42a27b6d3014e",
                TaskID : config.TaskID,
                Ts: now.toString()
            },
            "ParamBatch" : []
        }

        // 混入随机因子
        data["Params"]["Ts"] = (parseInt(data["Params"]["Ts"]) + index).toString()
        var source = "TaskHackerReport" + data["Params"]["From"] + data["Params"]["ReportPath"] + data["Params"]["TaskID"] + data["Params"]["Ts"]
        var hash = crypto.SHA256(source).toString()
        data["Params"]["Hash"] = hash
        var signature = await utils.Sign(data["Params"]["Hash"], config.PrivateKey)
        data["Params"]["Signature"] = signature

        // 批量生成
        for(var i=0;i<10;i++) {
            var tempParam = {
                "MC_Call" : "PublishReportByHacker", // 调用合约脚本
                From: "0402e7ebab98d26359b235dd0405c6fb86",
                ReportPath: "QmNnGUHwagNAzhBUwGvsfCMvuVKfL294eKPu1r2juXVN3J",
                // TaskID: "56d8c9fc346bec5fdea30f5d4c8811a5f9ea0d6838cb8d6a9950c2afdac55874",
                TaskID : config.TaskID,
                Ts: now.toString()
            }
            tempParam["Ts"] = (parseInt(now) + i).toString()
            var source = "TaskHackerReport" + tempParam["From"] + tempParam["ReportPath"] + tempParam["TaskID"] + tempParam["Ts"]
            var hash = crypto.SHA256(source).toString()
            tempParam["Hash"] = hash
            var signature = await utils.Sign(tempParam["Hash"], config.PrivateKey)
            tempParam["Signature"] = signature
            data["ParamBatch"].push(tempParam)
        }
        var option = {
            // hostname : "localhost",
            hostname : "192.168.10.202",
            path : "/api/proxy",
            port : 10001,
            method : "POST",
            headers : {
                "Content-Type" : "application/json"
            }
        }
    
        var req = http.request(option, function(res) {
            res.on("data", function(data){
                resolve(data.toString())
            })
        })
    
        req.on("error", function(err) {
            console.log(err)
        })
    
        req.write(JSON.stringify(data))
        req.end()
    })
}

async function Sleep(time){
    return new Promise(resolve=>{
        setTimeout(function(){
            resolve()
        }, time)
    })
}

async function main(){
    for(var i=0;true;i++) {
        var promises = []
        for(var k=0;k<1;k++) {
            promises.push(doPost(k))
        }
        await Promise.all(promises)
        await Sleep(1000)
        console.log(`done : ${i}`)
        // break
    }
}

main()