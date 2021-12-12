// 整体状态的数据结构
var Status = {
    // 指定挖矿账号（后续应该由代表专家来做）
    Miners : ["047204499d849948aaffdec7ce2703f5b3"],

    Hackers : [],
    Experts : [],
    Enterprises : [],
    Tasks : [],
}

var Testin = {
    // 全量读取区块，构造整体的世界状态。非常耗性能，上线时必须做缓存处理
    BuildWorldStatus : function(){
        var topBlock = JSON.parse(MC_GetTopBlock())
        var blocks = MC_GetBlockByRange(1, parseInt(topBlock.Number))
        for (var i=0;i<blocks.length;i++) {
            var block = JSON.parse(blocks[i])
            for (k=0;k<block.Transactions.length;k++) {
                // 构建黑客身份状态
                var trans = block.Transactions[k]
                if (trans.Type == "RegisterHacker") {
                    Status.Hackers.push(trans.Hacker)
                }
            }
        }
    },
    Class : {
        Transaction : {
            // 从params中构造一个交易对象
            New : function(params){
                function Trans(_params){
                    this.Type = _params.Type
                    if (this.Type == "RegisterHacker") {
                        this.Hacker = _params.Hacker
                        this.CheckSign = function(){
                            return this.Hacker.CheckSign()
                        }
                    }
                }

                return new Trans(params)
            },
        },

        /**
         * @params.topBlock 上一个区块数据
         * @params.transactions 本次交易的打包数据
         */
        Block : {
            New : function(params){
                function Block(_params) {
                    var topBlock = _params.topBlock
                }

                return new Block(params)
            }
        },

        Miner : {

        },
        Hacker : {
            New : function(params){
                function Hacker(_params){
                    this.From = _params.From
                    this.Name = _params.Name
                    this.Resume = _params.Resume
                    this.Qualification = _params.Qualification
                    this.Ts = _params.Ts
                    
                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                Hacker.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "Hacker" + this.From + this.Name + this.Qualification + this.Resume + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 排重，仅使用NodeID和Name来进行排重
                Hacker.prototype.IsExist = function() {
                    for(var i=0;i<Status.Hackers.length;i++) {
                        if (this.Name == Status.Hackers[i].Name) {
                            return true
                        }

                        if (this.From == Status.Hackers[i].From) {
                            return true
                        }
                    }

                    return false
                }

                return new Hacker(params)
            }
        },
        Expert : {

        },
        Enterprise : {

        },
        Task : {

        }
    }
}

// 注册成为测试员
exports.RegisterHacker = function(params) {
    var hacker = Testin.Class.Hacker.New(params)
    if (hacker.CheckSign() == false) {
        console.log("提交数据签名校验失败: RegisterHacker")
        return 
    }

    Testin.BuildWorldStatus()
    if (hacker.IsExist()) {
        console.log("该测试员名称已存在，或您已申请过测试员身份，无法重复申请")
        return 
    }

    // 前端不会给交易加类型，这里的接口都是把提交参数封装成一个
    // 个交易，所以需要在这里强制加上类型
    var transParam = {
        Type : "RegisterHacker",
        Hacker : hacker
    }
    var trans = Testin.Class.Transaction.New(transParam)
    if (trans.CheckSign() == false) {
        console.log("交易签名校验失败: RegisterHacker")
        return 
    }

    // 把交易缓存起来，等待矿工拉取
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    var thisBlockNumber = parseInt(topBlock.Number) + 1
    MC_SetCache("transCache-" + thisBlockNumber + "-" + trans.Hash, JSON.stringify(trans))
}

// 注册成为专家
exports.RegisterExpert = function(params) {

}

// 注册成为企业
exports.RegisterEnterprise = function(params) {

}

// 企业发布任务
exports.PublishTaskByEnterprise = function(params) {

}

// 专家审核任务
exports.ReviewTaskByExpert = function(params) {

}

// 测试员申请任务
exports.ApplyTaskByHacker = function(params){

}

// 企业授权测试员
exports.AuthorizationHackerToTaskByEnterprise = function(params){

}

// 测试员提交报告
exports.PublishReportByHacker = function(params) {

}

// 专家评审报告
exports.ReviewReportByExpert = function(params) {

}

// 企业认领报告
exports.ConfirmTaskByEnterprise = function(params) {

}

// 新区块
// @params.block 新区块的整体内容
exports.DoPackage = function(params) {
    // 把交易缓存起来，等待矿工拉取
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)
    console.log(JSON.stringify(params))

    // var block = {
    //     Hash : "",
    //     PreviousHash : "",
    //     Number : thisBlockNumber + "",
    //     Transactions : [],
    //     MerkleRoot : "",
    //     Miner : "",
    //     Ts : (+new Date()).toString(),
    //     Sign : "",
    // }
    // block.Transactions.push(trans)

    // MC_AddNewBlock(JSON.stringify(block))
}

/*
{
    "Hash" : "HiThisIsTestinProject",
    "PreviousHash" : "",
    "Number" : "1",
    "Transactions" : [],
    "MerkleRoot" : "",
    "Miner" : "",
    "Ts" : "",
    "Signature" : ""
}
*/