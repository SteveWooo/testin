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
        function loadTrans(trans){
            if (trans.Type == "RegisterHacker") {
                Status.Hackers.push(trans.Hacker)
            }
            if (trans.Type == "RegisterEnterprise") {
                Status.Enterprises.push(trans.Enterprise)
            }
            if (trans.Type == "RegisterExpert") {
                Status.Experts.push(trans.Expert)
            }
            if (trans.Type == "PublishTaskByEnterprise") {
                Status.Tasks.push(trans.Task)
            }
        }

        var topBlock = JSON.parse(MC_GetTopBlock())
        var blocks = MC_GetBlockByRange(1, parseInt(topBlock.Number))
        for (var i=0;i<blocks.length;i++) {
            var block = JSON.parse(blocks[i])
            for (k=0;k<block.Transactions.length;k++) {
                // 构建黑客身份状态
                var trans = block.Transactions[k]
                loadTrans(trans)
            }
        }

        // 缓存中的合法交易也要加入其中，防止交易重复提交
        // 这里不需要考虑缓存交易中存在 注册企业+发布任务 这种交易组合
        // 因为发布任务的交易，不能在注册企业被上链之前认可
        transCache = MC_GetCacheByPrefix("transCache-" + (parseInt(topBlock.Number) + 1) + "-")
        transCache = JSON.parse(transCache)
        for (var i=0;i<transCache.length;i++) {
            var trans = JSON.parse(transCache[i])
            loadTrans(trans)
        }

        // console.log(JSON.stringify(Status))
    },
    Class : {
        Transaction : {
            // 从params中构造一个交易对象
            New : function(params){
                /**
                 * Interface:
                 * CheckSign() 交易签名检查
                 * CheckWorldStatus() 交易世界状态检查（如排重，业务逻辑梳理
                 */
                function Trans(_params){
                    this.Type = _params.Type
                    if (this.Type == "RegisterHacker") {
                        this.Hacker = Testin.Class.Hacker.New(_params.Hacker)
                        this.Hash = this.Hacker.Hash
                        this.CheckSign = function(){
                            return this.Hacker.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.Hacker.CheckWorldStatus()
                        }
                        return this
                    }

                    if (this.Type == "RegisterEnterprise") {
                        this.Enterprise = Testin.Class.Enterprise.New(_params.Enterprise)
                        this.Hash = this.Enterprise.Hash
                        this.CheckSign = function(){
                            return this.Enterprise.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.Enterprise.CheckWorldStatus()
                        }
                        return this
                    }

                    if (this.Type == "RegisterExpert") {
                        this.Expert = Testin.Class.Expert.New(_params.Expert)
                        this.Hash = this.Expert.Hash
                        this.CheckSign = function(){
                            return this.Expert.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.Expert.CheckWorldStatus()
                        }
                        return this
                    }

                    if (this.Type == "PublishTaskByEnterprise") {
                        this.Task = Testin.Class.Task.New(_params.Task)
                        this.Hash = this.Task.Hash
                        this.CheckSign = function(){
                            return this.Task.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.Task.CheckWorldStatus()
                        }
                        return this
                    }
                }

                return new Trans(params)
            },
        },

        /**
         * @params.transactions 本次交易的打包数据
         */
        Block : {
            New : function(params){
                function Block(_params) {
                    this.MerkleRoot = _params.MerkleRoot
                    this.Miner = _params.Miner
                    this.Number = _params.Number
                    this.PreviousHash = _params.PreviousHash
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Transactions = _params.Transactions
                    this.Signature = _params.Signature

                    for (var i=0;i<this.Transactions.length;i++) {
                        var trans = Testin.Class.Transaction.New(this.Transactions[i])
                        this.Transactions[i] = trans
                    }

                    return this
                }

                Block.prototype.CheckSign = function(){
                    // 检查merkle根
                    var transSource = ""
                    for (var i=0;i<this.Transactions.length;i++) {
                        transSource = transSource + this.Transactions[i].Hash
                        // 顺便把交易本身的签名检查一份
                        if (this.Transactions[i].CheckSign() == false) {
                            console.log("交易校验失败")
                            return false
                        }
                    }
                    if (this.MerkleRoot != MC_Sha256(transSource)) {
                        console.log("Merkle Root校验失败")
                        return false
                    }

                    var source = "Block" + this.MerkleRoot + this.Miner + this.Number + this.PreviousHash + this.Ts
                    var hash = MC_Sha256(source)
                    if(this.Hash != hash) {
                        console.log("区块Hash校验失败")
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.Miner)) { // 签名校验失败
                        console.log("区块签名校验失败")
                        return false
                    }

                    return true
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
                    return this
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

                // 业务查询调用接口：
                Hacker.prototype.CheckWorldStatus = function(){
                    for(var i=0;i<Status.Hackers.length;i++) {
                        if (this.Name == Status.Hackers[i].Name) {
                            return false
                        }

                        if (this.From == Status.Hackers[i].From) {
                            return false
                        }
                    }

                    return true
                }

                return new Hacker(params)
            }
        },
        Expert : {
            New : function(params){
                function Expert(_params){
                    this.From = _params.From
                    this.Name = _params.Name
                    this.Qualification = _params.Qualification
                    this.Resume = _params.Resume
                    this.Ts = _params.Ts
                    
                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                    return this
                }

                // 检查提交签名
                Expert.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "Expert" + this.From + this.Name + this.Qualification + this.Resume + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 业务查询调用接口：
                Expert.prototype.CheckWorldStatus = function(){
                    for(var i=0;i<Status.Experts.length;i++) {
                        if (this.Name == Status.Experts[i].Name) {
                            return false
                        }

                        if (this.From == Status.Experts[i].From) {
                            return false
                        }
                    }

                    return true
                }

                return new Expert(params)
            }
        },
        Enterprise : {
            New : function(params){
                function Enterprise(_params){
                    this.Connection = _params.Connection
                    this.From = _params.From
                    this.LogoPath = _params.LogoPath
                    this.Name = _params.Name
                    this.SocialCode = _params.SocialCode
                    this.Ts = _params.Ts
                    
                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                    return this
                }

                // 检查提交签名
                Enterprise.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "Enterprise" + this.Connection + this.From + this.LogoPath + this.Name + this.SocialCode + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 业务查询调用接口：
                Enterprise.prototype.CheckWorldStatus = function(){
                    for(var i=0;i<Status.Enterprises.length;i++) {
                        if (this.Name == Status.Enterprises[i].Name || this.SocialCode == Status.Enterprises[i].SocialCode) {
                            return false
                        }

                        if (this.From == Status.Enterprises[i].From) {
                            return false
                        }
                    }

                    return true
                }

                return new Enterprise(params)
            }
        },
        /**
         * 测试任务对象
         */
        Task : {
            New : function(params) {
                function Task(_params) {
                    // 初始化参数
                    this.Budget = _params.Budget
                    this.From = _params.From
                    this.MaxAuthorizationCount = _params.MaxAuthorizationCount
                    this.Name = _params.Name
                    this.Require = _params.Require
                    this.Resume = _params.Resume
                    this.Ts = _params.Ts

                    // 运作参数（不参与签名
                    this.Hackers = []
                    this.IsPublic = false // 默认任务为不公开任务

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                Task.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "Task" + this.Budget + this.From + this.MaxAuthorizationCount + this.Name + this.Require + this.Resume + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 检查创建人是否已经有认证的企业
                Task.prototype.CheckWorldStatus = function(){
                    var isEnterperiseAccount = false
                    for (var i=0;i<Status.Enterprises.length;i++) {
                        if (this.From == Status.Enterprises[i].From) {
                            isEnterperiseAccount = true
                        }
                    }
                    if (isEnterperiseAccount == false) {
                        return false
                    }

                    return true
                }

                return new Task(params)
            }
        },
        /**
         * 任务里面的测试员列表
         * @param params.Hacker 测试员对象
         */
        TaskHacker : {
            New : function(params){

            }
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

    // 前端不会给交易加类型，这里的接口都是把提交参数封装成一个
    // 个交易，所以需要在这里强制加上类型
    var transParam = {
        Type : "RegisterHacker",
        Hacker : hacker
    }
    var trans = Testin.Class.Transaction.New(transParam)

    // 检查世界状态
    if (trans.CheckWorldStatus() == false) {
        console.log("交易世界状态检查失败：" + trans.Type)
        return 
    }

    // 把交易缓存起来，等待矿工拉取
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    var thisBlockNumber = parseInt(topBlock.Number) + 1
    MC_SetCache("transCache-" + thisBlockNumber + "-" + trans.Hash, JSON.stringify(trans))
}

// 注册成为企业
exports.RegisterEnterprise = function(params) {
    var enterprise = Testin.Class.Enterprise.New(params)
    if (enterprise.CheckSign() == false) {
        console.log("提交数据签名校验失败: RegisterEnterprise")
        return 
    }

    Testin.BuildWorldStatus()
    var transParam = {
        Type : "RegisterEnterprise",
        Enterprise : enterprise
    }

    var trans = Testin.Class.Transaction.New(transParam)
    // 检查世界状态
    if (trans.CheckWorldStatus() == false) {
        console.log("交易世界状态检查失败：" + trans.Type)
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
    var expert = Testin.Class.Expert.New(params)
    if (expert.CheckSign() == false) {
        console.log("提交数据签名校验失败：RegisterExpert")
        return 
    }

    Testin.BuildWorldStatus()
    var transParam = {
        Type : "RegisterExpert",
        Expert : expert
    }

    var trans = Testin.Class.Transaction.New(transParam)
    // 检查世界状态
    if (trans.CheckWorldStatus() == false) {
        console.log("交易世界状态检查失败：" + trans.Type)
        return 
    }

    // 把交易缓存起来，等待矿工拉取
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    var thisBlockNumber = parseInt(topBlock.Number) + 1
    MC_SetCache("transCache-" + thisBlockNumber + "-" + trans.Hash, JSON.stringify(trans))
}

// 企业发布任务
exports.PublishTaskByEnterprise = function(params) {
    var task = Testin.Class.Task.New(params)
    if (task.CheckSign() == false ) {
        console.log("提交数据签名校验失败：PublishTaskByEnterprise");
        return 
    }

    Testin.BuildWorldStatus()
    var transParam = {
        Type : "PublishTaskByEnterprise",
        Task : task
    }
    var trans = Testin.Class.Transaction.New(transParam)
    // 检查世界状态
    if (trans.CheckWorldStatus() == false) {
        console.log("交易世界状态检查失败：" + trans.Type)
        return 
    }
    
    // 把交易缓存起来，等待矿工拉取
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    var thisBlockNumber = parseInt(topBlock.Number) + 1
    MC_SetCache("transCache-" + thisBlockNumber + "-" + trans.Hash, JSON.stringify(trans))
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
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    // 构建一个区块对象
    // console.log(JSON.stringify(params))
    var block = Testin.Class.Block.New(params.Block)

    // TODO 共识逻辑
    if (block.Miner != Status.Miners[0]) {
        console.log("非法打包")
        return 
    }

    if (block.PreviousHash != topBlock.Hash) {
        console.log("区块PreviousHash错误")
        return 
    }

    if (parseInt(block.Number) != parseInt(topBlock.Number) + 1) {
        console.log("区块编号错误")
        return 
    }

    if (block.CheckSign() == false ) {
        console.log("区块校验失败")
        return 
    }

    // 业务校验
    Testin.BuildWorldStatus()
    for (var i=0;i<block.Transactions.length;i++) {
        if (block.Transactions[i].CheckWorldStatus() == false) {
            console.log("交易世界状态检查失败")
            return 
        }
    }

    // 给交易打上Nonce
    for (var i=0;i<block.Transactions.length;i++) {
        block.Transactions[i].Nonce = i + ""
    }

    // 删除以往的所有相关缓存，防止缓存冗余
    MC_DeleteCacheByPrefix("transCache-" + topBlock.Number + "-")

    // 写入新区块
    MC_AddNewBlock(JSON.stringify(block))
    console.log("新区块写入完成：" + block.Number)

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