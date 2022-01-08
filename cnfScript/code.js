// 整体状态的数据结构
var Status = {
    // 指定挖矿账号（后续应该由代表专家来做）
    // Miners : ["047204499d849948aaffdec7ce2703f5b3","0433cd50fa5977da115025e90cf5698c08",
    // "043abf9b64da3cf82a6833d827a6a60cb1","04c52654247aa39be86b5ce356ac7e24f8",
    // "0492ec813ab9ce7c94e49c84abcb0c7d64","049075a782f699fd18ca64cf7ccb0b7ef5",
    // "0429285759acca19681489804066c123fe","04b291af0ad8ed77f167d2d89da6dd310a",
    // "04ae7e5a2f9b426f0f18df1f4629e408ad","043f7acc95c1bf43ebd4bb7313979f427e",
    // "04961d37561a0cb5f8efaf95b555943b77","04bde666ba0e9078328897a8087cccc14a",
    // "041ba0e83f3e7962a388f5c0296ccacfe5","04d1f611569df79cff3d05a6aa8553bc7e",
    // "0475bb62e72d9fe2d92e542ee4f7aefd24","045ff90b3a6ea54f58178aeb4a6c60f81c",
    // "04f4b26de40eb5fc31bca10918bf414d41","04cb4cbd636d6694e4f992f54f65f1daa8",
    // "04cf8554facafcac058c70aa6ffc38a139","0458df0e6a93a464ca2de7fd4a13049b5b"],
    Miners : ["047204499d849948aaffdec7ce2703f5b3","0433cd50fa5977da115025e90cf5698c08",
    "043abf9b64da3cf82a6833d827a6a60cb1","04c52654247aa39be86b5ce356ac7e24f8",
    "0492ec813ab9ce7c94e49c84abcb0c7d64"],

    Hackers : [],
    Experts : [],
    Enterprises : [],
    Tasks : [],
    TaskHackers : [],

    Reputations : {}
}

var Testin = {
    // 把交易数据加载到全局变量Status中
    loadTrans : function(trans){
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
        if (trans.Type == "ApplyTaskByHacker") {
            Status.TaskHackers.push(trans.TaskHacker)
        }
        if (trans.Type == "AuthorizationHackerToTaskByEnterprise") {
            for (var i=0;i<Status.TaskHackers.length;i++) {
                if (Status.TaskHackers[i].TaskID == trans.AuthorizationHackerToTaskByEnterprise.TaskID && Status.TaskHackers[i].From == trans.AuthorizationHackerToTaskByEnterprise.HackerID) {
                    Status.TaskHackers[i].IsPermission = "true"
                    // Status.TaskHackers[i].ExpertList = ["047204499d849948aaffdec7ce2703f5b3"] // hard code
                    Status.TaskHackers[i].ExpertList = [] // 改成任何人都可以评价
                    Status.TaskHackers[i].PermissionInformation = trans.AuthorizationHackerToTaskByEnterprise.PermissionInformation
                }
            }
        }
        if (trans.Type == "PublishReportByHacker" ) {
            for (var i=0;i<Status.TaskHackers.length;i++) {
                if (Status.TaskHackers[i].TaskID == trans.TaskHackerReport.TaskID && Status.TaskHackers[i].From == trans.TaskHackerReport.From) {
                    Status.TaskHackers[i].ReportPath = trans.TaskHackerReport.ReportPath
                    // 同时清空专家评审信息
                    Status.TaskHackers[i].ExpertReviewReports = []
                    break
                }
            }
        }

        // 调整信誉值
        if (trans.Type == "ReviewReportByExpert" ) {
            for (var i=0;i<Status.TaskHackers.length;i++) {
                if (Status.TaskHackers[i].TaskID == trans.ExpertReviewReport.TaskID && Status.TaskHackers[i].From == trans.ExpertReviewReport.HackerID) {
                    Status.TaskHackers[i].ExpertReviewReports.push(trans.ExpertReviewReport)
                    if(Status.Reputations[trans.ExpertReviewReport.From] == undefined) {
                        Status.Reputations[trans.ExpertReviewReport.From] = 0.0
                    }
                    Status.Reputations[trans.ExpertReviewReport.From] = Status.Reputations[trans.ExpertReviewReport.From] + 1.0
                    
                    break
                }
            }
        }
    },

    // 加载一个区块，把区块内的交易信息加载到全局中去
    loadBlock : function(blockJSON) {
        var block = JSON.parse(blockJSON)
        // 加载区块矿工，调整信誉值
        if(block.Miner.length == 34) { // 防止加载第一个区块的空矿工
            Status.Reputations[block.Miner] = 0.0 // 成功出块的矿工信誉值归零。TODO : 这里需要兑换成虚拟货币
        }
        for (k=0;k<block.Transactions.length;k++) {
            // 构建黑客身份状态
            var trans = block.Transactions[k]
            Testin.loadTrans(trans)
        }
    },

    // 读取所有区块，重新构建一个Status出来
    reloadAllBlock : function(){
        for(var i=0;i<Status.Miners.length;i++) {
            if (Status.Reputations[Status.Miners[i]] == undefined) {
                Status.Reputations[Status.Miners[i]] = 0.0
            }
        }
        var topBlock = JSON.parse(MC_GetTopBlock())
        var blocks = MC_GetBlockByRange(1, parseInt(topBlock.Number))
        for (var i=0;i<blocks.length;i++) {
            Testin.loadBlock(blocks[i])
        }
    },

    // 全量读取区块，构造整体的世界状态。非常耗性能，上线时必须做缓存处理
    BuildWorldStatus : function(param){ 
        var worldStatus = MC_GetCache("worldStatus")
        if(worldStatus == "") { // 这种就是网络启动的时候会调用一次
            // 初始化矿工信誉值
            Testin.reloadAllBlock()
            MC_SetCache("worldStatus", JSON.stringify(Status))
        } else { // 这里比较常被调用
            Status = JSON.parse(worldStatus)
        }

        // 有时候需要重新构造一次区块数据，比如更新算法后
        // if (param != undefined && param.ReloadAllBlock == true) {
        //     Testin.ReloadAllBlock()
        // }

        // 缓存中的合法交易也要加入其中，防止交易重复提交
        // 这里不需要考虑缓存交易中存在 注册企业+发布任务 这种交易组合
        // 因为发布任务的交易，不能在注册企业被上链之前认可
        if (param != undefined && param.LoadCache == true) {
            var topBlock = JSON.parse(MC_GetTopBlock())
            transCache = MC_GetCacheByPrefix("transCache-" + (parseInt(topBlock.Number) + 1) + "-")
            transCache = JSON.parse(transCache)
            for (var i=0;i<transCache.length;i++) {
                var trans = JSON.parse(transCache[i])
                Testin.loadTrans(trans)
            }
        }
        // console.log(JSON.stringify(Status))
    },

    // 动态规划最新区块数据
    AppendNewBlockToWorldStatus : function(blockJSON){
        MC_DeleteCacheByPrefix("worldStatus")
        Testin.loadBlock(blockJSON)
        MC_SetCache("worldStatus", JSON.stringify(Status))
    },

    // 共识过程中出现的操作类
    Consensus : {
        // 打包意向发布
        PackageIntention : {
            New : function(params) {
                function PackageIntention(_params) {
                    // 初始化参数
                    this.From = _params.From
                    this.Intention = _params.Intention
                    this.Term = _params.Term
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                PackageIntention.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "DoPackageIntention" + this.From + this.Intention + this.Term + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 查询投票参与者是否属于矿工列表
                PackageIntention.prototype.CheckWorldStatus = function(){
                    for(var i=0;i<Status.Miners.length;i++) {
                        if (Status.Miners[i] == this.From) {
                            return true
                        }
                    }
                    return false                    
                }

                return new PackageIntention(params)
            }
        },

        // 打包权威排行分布
        IntentionRank : {
            New : function(params) {
                function IntentionRank(_params) {
                    // 初始化参数
                    this.From = _params.From
                    this.Rank_1 = _params.Rank_1
                    this.Term = _params.Term
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                IntentionRank.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "ShareIntentionRank" + this.From + this.Rank_1 + this.Term + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 查询投票参与者是否属于矿工列表
                IntentionRank.prototype.CheckWorldStatus = function(){
                    for(var i=0;i<Status.Miners.length;i++) {
                        if (Status.Miners[i] == this.From) {
                            return true
                        }
                    }
                    return false                    
                }

                return new IntentionRank(params)
            }
        }
    },

    // 业务操作类，比如授权测试任务
    Operation : {
        // TODO 分配专家
        AuthorizationHackerToTaskByEnterprise : {
            New : function(params) {
                function Auth(_params) {
                    // 初始化参数
                    this.From = _params.From
                    this.HackerID = _params.HackerID
                    this.PermissionInformation = _params.PermissionInformation
                    this.TaskID = _params.TaskID
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                Auth.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "AuthorizationHackerToTaskByEnterprise" + this.From + this.HackerID + this.PermissionInformation + this.TaskID + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 检查任务 & 测试员ID的组合是否存在
                // 检查企业是否创建人
                // 检查是否已经授权
                Auth.prototype.CheckWorldStatus = function(){
                    var isItemExist = false
                    var isEnterpriseSelf = false // 检查这个任务发起者是否本人，不是的话不能授权
                    var isAlreadyAuth = false
                    for (var i=0;i<Status.TaskHackers.length;i++) {
                        if (Status.TaskHackers[i].TaskID == this.TaskID && Status.TaskHackers[i].From == this.HackerID) {
                            isItemExist = true
                            if (Status.TaskHackers[i].IsPermission == "true") {
                                isAlreadyAuth = true
                            }
                            for (var k=0;k<Status.Tasks.length;k++) {
                                if (Status.Tasks[k].Hash == this.TaskID && Status.Tasks[k].From == this.From) {
                                    isEnterpriseSelf = true 
                                    break
                                }
                            }
                            break
                        }
                    }
                    if (isAlreadyAuth == true) {
                        console.log("无法重复授权")
                        return false
                    }
                    if (isItemExist == false) {
                        console.log("测试员申请不存在")
                        return false
                    }
                    if (isEnterpriseSelf == false) {
                        console.log("该任务不是授权人发起的")
                        return false
                    }

                    return true
                }

                return new Auth(params)
            }
        },

        // 测试报告提交
        /**
         * @param.From hackerId 
         * @param.TaskID
         * @param.ReportPath 报告ipfs地址
         */
         TaskHackerReport : {
            New : function(params) {
                function TaskHackerReport(_params) {
                    // 初始化参数
                    this.From = _params.From
                    this.ReportPath = _params.ReportPath
                    this.TaskID = _params.TaskID
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                TaskHackerReport.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "TaskHackerReport" + this.From + this.ReportPath + this.TaskID + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 检查任务 & 测试员ID的组合是否存在
                TaskHackerReport.prototype.CheckWorldStatus = function(){
                    var isItemExist = false
                    for (var i=0;i<Status.TaskHackers.length;i++) {
                        if (Status.TaskHackers[i].TaskID == this.TaskID && Status.TaskHackers[i].From == this.From) {
                            isItemExist = true
                            break
                        }
                    }
                    if (isItemExist == false) {
                        console.log("测试员申请不存在")
                        return false
                    }

                    return true
                }

                return new TaskHackerReport(params)
            }
         },
         /**
          * 专家评审报告：
          * @params From 专家ID
          * @params TaskID 任务ID
          * @parmas HackerID 测试员id，和任务id一起定位到具体的报告位置
          * @params Score 评分
          * @params Memo 评语
          */
         ExpertReviewReport : {
            New : function(params) {
                function ExpertReviewReport(_params) {
                    // 初始化参数
                    this.From = _params.From
                    this.HackerID = _params.HackerID
                    this.Memo = _params.Memo
                    this.Score = _params.Score
                    this.TaskID = _params.TaskID
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature
                }

                // 检查提交签名
                ExpertReviewReport.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "ExpertReviewReport" + this.From + this.HackerID + this.Memo + this.Score + this.TaskID + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                ExpertReviewReport.prototype.CheckWorldStatus = function(){
                    // 检查任务 & 测试员ID的组合是否存在
                    var isItemExist = false
                    // 检查自己是否已经评分
                    var isAlreadyReview = false

                    for (var i=0;i<Status.TaskHackers.length;i++) {
                        if (Status.TaskHackers[i].TaskID == this.TaskID && Status.TaskHackers[i].From == this.HackerID) {
                            isItemExist = true
                            for (var k=0;k<Status.TaskHackers[i].ExpertReviewReports.length;k++) {
                                if (this.From == Status.TaskHackers[i].ExpertReviewReports[k].From) {
                                    isAlreadyReview = true
                                }
                            }

                            break
                        }
                    }
                    if (isItemExist == false) {
                        console.log("测试员申请不存在")
                        return false
                    }
                    if (isAlreadyReview == true) {
                        console.log("专家重复评价")
                        return false
                    }

                    return true
                }

                return new ExpertReviewReport(params)
            }
         }
    },

    // 对象类
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

                    if (this.Type == "ApplyTaskByHacker") {
                        this.TaskHacker = Testin.Class.TaskHacker.New(_params.TaskHacker)
                        this.Hash = this.TaskHacker.Hash
                        this.CheckSign = function(){
                            return this.TaskHacker.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.TaskHacker.CheckWorldStatus()
                        }
                    }

                    if (this.Type == "AuthorizationHackerToTaskByEnterprise") {
                        this.AuthorizationHackerToTaskByEnterprise = Testin.Operation.AuthorizationHackerToTaskByEnterprise.New(_params.AuthorizationHackerToTaskByEnterprise)
                        this.Hash = this.AuthorizationHackerToTaskByEnterprise.Hash
                        this.CheckSign = function(){
                            return this.AuthorizationHackerToTaskByEnterprise.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.AuthorizationHackerToTaskByEnterprise.CheckWorldStatus()
                        }
                    }
                    if (this.Type == "PublishReportByHacker" ) {
                        this.TaskHackerReport = Testin.Operation.TaskHackerReport.New(_params.TaskHackerReport)
                        this.Hash = this.TaskHackerReport.Hash
                        this.CheckSign = function(){
                            return this.TaskHackerReport.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.TaskHackerReport.CheckWorldStatus()
                        }
                    }

                    if (this.Type == "ReviewReportByExpert") {
                        this.ExpertReviewReport = Testin.Operation.ExpertReviewReport.New(_params.ExpertReviewReport)
                        this.Hash = this.ExpertReviewReport.Hash
                        this.CheckSign = function(){
                            return this.ExpertReviewReport.CheckSign()
                        }
                        this.CheckWorldStatus = function(){
                            return this.ExpertReviewReport.CheckWorldStatus()
                        }
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
                            console.log("测试员重名错误")
                            return false
                        }

                        if (this.From == Status.Hackers[i].From) {
                            console.log("一个账号不可重复申请测试员")
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
                    this.TaskHackers = []
                    this.IsPublic = "false" // 默认任务为不公开任务

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
                            break 
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
         * 任务里面的测试员列表（重要的报告对象）
         * @param params.Hacker 测试员对象
         */
        TaskHacker : {
            New : function(params){
                function TaskHacker(_params) {
                    // 初始化参数
                    this.From = _params.From // 与HackerID绑定
                    this.TaskID = _params.TaskID
                    this.Ts = _params.Ts

                    this.Hash = _params.Hash
                    this.Signature = _params.Signature

                    this.IsPermission = "false"
                    this.PermissionInformation = ""
                    this.ExpertList = []
                    this.ReportPath = ""

                    this.ExpertReviewReports = []
                    this.thisNegotiations = []
                }

                // 检查提交签名
                TaskHacker.prototype.CheckSign = function(){
                    // 校验Hash
                    var source = "TaskHacker" + this.From + this.TaskID + this.Ts
                    var hash = MC_Sha256(source)
                    if (this.Hash != hash) { // 哈希校验失败
                        return false
                    }

                    if (!MC_Secp256k1_Check(hash, this.Signature, this.From)) { // 签名校验失败
                        return false
                    }

                    return true
                }

                // 检查任务是否存在，
                // 检查是否已经提交过参与信息，
                // 检查已经通过的人数是否达到上限
                TaskHacker.prototype.CheckWorldStatus = function(){
                    var isTaskExist = false
                    var neverJoin = true
                    var isEnoughPermissionHacker = false
                    var isHacker = false

                    for (var i=0;i<Status.Hackers.length;i++) {
                        if (this.From == Status.Hackers[i].From) {
                            isHacker = true
                        }
                    }
                    if (isHacker == false) {
                        console.log("申请人未有测试员资格")
                        return false
                    }
                    for (var i=0;i<Status.Tasks.length;i++) {
                        if (this.TaskID == Status.Tasks[i].Hash) {
                            isTaskExist = true

                            // 检查已经被授权的人数是否达到上限
                            var permissionCount = 0
                            for(var k=0;k<Status.TaskHackers.length;k++) {
                                if (Status.TaskHackers[k].TaskID != this.TaskID) {
                                    continue 
                                }

                                if (Status.TaskHackers[k].IsPermission == "true") {
                                    permissionCount = permissionCount + 1
                                }

                                if (Status.TaskHackers[k].From == this.From) {
                                    neverJoin = false
                                }
                            }
                            if(permissionCount >= Status.Tasks[i].MaxAuthorizationCount) {
                                isEnoughPermissionHacker = true
                            }

                            break 
                        }
                    }
                    // 任务必须存在
                    if (isTaskExist == false) {
                        console.log("任务不存在")
                        return false
                    }
                    // 不可重复报名参与
                    if (neverJoin == false) {
                        console.log("申请人不可重复申请")
                        return false
                    }
                    // 不可超过限定人数
                    if (isEnoughPermissionHacker == true) {
                        console.log("任务人数达到上限")
                        return false
                    }

                    return true
                }

                return new TaskHacker(params)
            }
        },
    }
}

// 注册成为测试员
exports.RegisterHacker = function(params) {
    var hacker = Testin.Class.Hacker.New(params)
    if (hacker.CheckSign() == false) {
        console.log("提交数据签名校验失败: RegisterHacker")
        return 
    }
    Testin.BuildWorldStatus({
        LoadCache : true
    })

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

    Testin.BuildWorldStatus({
        LoadCache : true
    })
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

    Testin.BuildWorldStatus({
        LoadCache : true
    })
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

    Testin.BuildWorldStatus({
        LoadCache : true
    })
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

// 测试员申请任务
exports.ApplyTaskByHacker = function(params){
    var taskHacker = Testin.Class.TaskHacker.New(params)
    if (taskHacker.CheckSign() == false ) {
        console.log("提交数据签名校验失败：PublishTaskByEnterprise");
        return 
    }

    Testin.BuildWorldStatus({
        LoadCache : true
    })
    var transParam = {
        Type : "ApplyTaskByHacker",
        TaskHacker : taskHacker
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

// 企业授权测试员
exports.AuthorizationHackerToTaskByEnterprise = function(params){
    var auth = Testin.Operation.AuthorizationHackerToTaskByEnterprise.New(params)
    if (auth.CheckSign() == false) {
        console.log("提交数据签名校验失败：AuthorizationHackerToTaskByEnterprise");
        return 
    }
    Testin.BuildWorldStatus({
        LoadCache : true
    })

    var transParam = {
        Type : "AuthorizationHackerToTaskByEnterprise",
        AuthorizationHackerToTaskByEnterprise : auth
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

// 测试员提交报告
exports.PublishReportByHacker = function(params) {
    var taskHackerReport = Testin.Operation.TaskHackerReport.New(params)
    if (taskHackerReport.CheckSign() == false) {
        console.log("提交数据签名校验失败：PublishReportByHacker");
        return 
    }
    Testin.BuildWorldStatus({
        LoadCache : true
    })

    var transParam = {
        Type : "PublishReportByHacker",
        TaskHackerReport : taskHackerReport
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

// 专家评审报告
exports.ReviewReportByExpert = function(params) {
    var review = Testin.Operation.ExpertReviewReport.New(params)
    if (review.CheckSign() == false) {
        console.log("提交数据签名校验失败：ReviewReportByExpert");
        return 
    }
    Testin.BuildWorldStatus({
        LoadCache : true
    })

    var transParam = {
        Type : "ReviewReportByExpert",
        ExpertReviewReport : review
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

// 企业认领报告
exports.ConfirmTaskByEnterprise = function(params) {

}

// 出块意向，调用该接口代表宣布本人是否愿意参与到下一轮的出块环节中
// 出块意向缓存到缓存中，等待矿工拉取使用
// @params.From
// @params.Intention true / false 是否参与出块
// @params.Term 轮次，用于标识重试次数
// @params.Ts 时间戳
// @parmas.Hash 
// @parmas.Signature
exports.DoPackageIntention = function(params) {
    var packageIntention = Testin.Consensus.PackageIntention.New(params.PackageIntention)
    if (packageIntention.CheckSign() == false ) {
        console.log("提交数据签名校验失败：DoPackageIntention");
        return 
    }
    Testin.BuildWorldStatus({
        LoadCache : true
    })
    // 直接检查打包意向的世界状态即可，共识机制不需要缝成一个交易
    if (packageIntention.CheckWorldStatus() == false ) {
        console.log("打包意向检查世界状态失败：DoPackageIntention");
        return 
    }

    // 缓存打包意向，等待矿工拉取
    // *这个打包意向缓存在被打包后，要清除掉
    // 把交易缓存起来，等待矿工拉取
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)
    var thisBlockNumber = parseInt(topBlock.Number) + 1
    // MC_DeleteCacheByPrefix("packageIntentionCache-" + thisBlockNumber + "-")
    MC_SetCache("packageIntentionCache-" + thisBlockNumber + "-" + packageIntention.Term + "-" + packageIntention.From, 
        JSON.stringify(packageIntention))
}

// 矿工节点收集到2/3个打包意向后，对有意向的出块者进行排名，广播第一名出来
// @params.From
// @params.Rank_1 string,string 第一名的NodeID，可用逗号分隔并列，并列NodeID通过哈希环随机算法进行排序
exports.ShareIntentionRank = function(params){
    var intentionRank = Testin.Consensus.IntentionRank.New(params.IntentionRank)
    if (intentionRank.CheckSign() == false ) {
        console.log("提交数据签名校验失败：ShareIntentionRank");
        return 
    }
    Testin.BuildWorldStatus({
        LoadCache : true
    })
    // 直接检查打包意向的世界状态即可，共识机制不需要缝成一个交易
    if (intentionRank.CheckWorldStatus() == false ) {
        console.log("打包意向检查世界状态失败：ShareIntentionRank");
        return 
    }

    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)
    var thisBlockNumber = parseInt(topBlock.Number) + 1

    // 后续打包的时候，严格检查这个缓存
    // 如果打包者确实达到了2/3的投票数量，就ok
    // 否则打包失败
    MC_SetCache("packageIntentionRankCache-" + thisBlockNumber + "-" + intentionRank.Term + "-" + intentionRank.From, 
        JSON.stringify(intentionRank))
}

// 核心共识机制，选主逻辑
function getMinerOfThisRound(term){
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)
    var thisBlockNumber = parseInt(topBlock.Number) + 1

    var intentionRanks = MC_GetCacheByPrefix("packageIntentionRankCache-" + thisBlockNumber + "-" + term + "-")
    intentionRanks = JSON.parse(intentionRanks)

    // 每个节点获得第一名次数
    var ranker1List = {}
    for(var i=0;i<intentionRanks.length;i++) {
        var ir = JSON.parse(intentionRanks[i])
        var maxRankers = ir["Rank_1"].split(",")
        for (var k=0;k<maxRankers.length;k++) {
            if(ranker1List[maxRankers[k]] == undefined) {
                ranker1List[maxRankers[k]] = {
                    count : 0
                }
            }
            ranker1List[maxRankers[k]].count = ranker1List[maxRankers[k]].count + 1
        }
    }

    // 获得第一名的最大次数
    var rank1MaxCount = 0
    for (var nid in ranker1List) {
        if (ranker1List[nid].count > rank1MaxCount) {
            rank1MaxCount = ranker1List[nid].count
        }
    }

    // 2/3 * n + 1
    // 最少的认可次数需要是2/3 TODO，其实应该1/3就ok，后面有平票逻辑处理。
    // 强行2/3的话，可能会导致重试率很高
    var minIntentionCount = Math.floor((2/3) * Status.Miners.length) + 1
    if (rank1MaxCount < minIntentionCount) {
        console.log("票数第一节点未达到最低要求区块")
        return undefined
    }

    // 出现平票情况的节点：
    rank1MaxMiner = []
    for (var nid in ranker1List) {
        if (ranker1List[nid].count == rank1MaxCount) {
            rank1MaxMiner.push(nid)
        }
    }

    // 对平票节点进行排序（哈希环伪随机算法
    // 03d40bf3 -> 64228339
    var startHash = topBlock.Hash
    var hashRangeStartNum = parseInt("0x" + startHash.substring(0, 8))
    // 冒泡
    for(var i=0;i<rank1MaxMiner.length;i++) {
        for (var k=i+1;k<rank1MaxMiner.length;k++) {
            var nodeIDNumForI = parseInt("0x" + rank1MaxMiner[i].substring(2, 10))
            var nodeIDNumForK = parseInt("0x" + rank1MaxMiner[k].substring(2, 10))
            // nodeIDNumForI = nodeIDNumForI - hashRangeStartNum
            // nodeIDNumForK = nodeIDNumForK - hashRangeStartNum

            // if (nodeIDNumForI < 0) {
            //     nodeIDNumForI = 4294967295 + nodeIDNumForI
            // }
            // if (nodeIDNumForK < 0) {
            //     nodeIDNumForK = 4294967295 + nodeIDNumForK
            // }

            if (nodeIDNumForI < nodeIDNumForK) {
                var temp = rank1MaxMiner[i]
                rank1MaxMiner[i] = rank1MaxMiner[k]
                rank1MaxMiner[k] = temp
            }
        }
    }
    // console.log(JSON.stringify(rank1MaxMiner))
    if (rank1MaxMiner.length == 0) {
        return undefined
    }
    
    var randomMinerIndex = MC_Mod(hashRangeStartNum, rank1MaxMiner.length)

    return rank1MaxMiner[randomMinerIndex]
}

// 1 master端收到request，就广播preprepare，让每个节点都缓存
// 一个preprepare包
exports.DoPBFTPreprepare = function(params) {
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    // 简单做算了
    var prepreParePack = params.PrepreParePack
    MC_SetCache("PBFTPrepreParePack-" + (parseInt(topBlock.Number) + 1) + "-" + prepreParePack.Hash, JSON.stringify(prepreParePack))

}

// 2 矿工持续监听preprepare包，如果有，就拉下来签名，并发布prepare包
exports.DoPBFTPrepare = function(params) {
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    // 简单做算了
    var preParePack = params.PreparePack
    MC_SetCache("PBFTPreparePack-" + (parseInt(topBlock.Number) + 1) + "-" + preParePack.Hash + "-" + preParePack.From, JSON.stringify(preParePack))
}

// 3 矿工发布prepare包后，持续监听prepare包，当prepare包达到2/3个数量时
// 就发布commit包，然后监听最新区块
// DoPBFTCommit 每次被调用，就查看缓存中commit是否达到2/3，如果达到2/3，就addNewBlock
exports.DoPBFTCommit = function(params) {
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    // 简单做算了
    var commit = params.CommitPack
    MC_SetCache("PBFTCommitPack-" + (parseInt(topBlock.Number) + 1) + "-" + commit.Hash + "-" + commit.From, JSON.stringify(commit))

    // 检查是否达到 2/3，到达就直接写块。矿工直接读块

    // 删除缓存的时候，要把commit以往的所有区块的缓存都删除掉
    // 因为由于这里读到2/3个包后就出块，肯定会冗余一些的缓存包的
}

// 新区块
// @params.block 新区块的整体内容
exports.DoPackage = function(params) {
    var topBlock = MC_GetTopBlock()
    topBlock = JSON.parse(topBlock)

    // 构建一个区块对象
    // console.log(JSON.stringify(params))
    var block = Testin.Class.Block.New(params.Block)
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
    Testin.BuildWorldStatus({
        LoadCache : false
    })

    // 核心：共识逻辑 - 获取当前轮次的Master
    var roundMiner = getMinerOfThisRound(params.Term)
    if(roundMiner == undefined) {
        console.log("非法打包: roundMiner为空")
        return 
    }
    if (block.Miner != roundMiner) {
        console.log("非法打包: "+block.Miner+" 非roundMiner:" + roundMiner)
        return 
    }

    for (var i=0;i<block.Transactions.length;i++) {
        if (block.Transactions[i].CheckWorldStatus() == false) {
            MC_DeleteCacheByPrefix("transCache-" + block.Number + "-" + block.Transactions[i].Hash)
            console.log("交易世界状态检查失败:" + JSON.stringify(block.Transactions[i]))
            return 
        }
    }

    // 给交易打上Nonce
    for (var i=0;i<block.Transactions.length;i++) {
        block.Transactions[i].Nonce = i + ""
    }

    // 删除以往的所有相关缓存，防止缓存冗余
    MC_DeleteCacheByPrefix("transCache-") // 删除交易缓存
    MC_DeleteCacheByPrefix("packageIntentionCache-") // 删除Intention缓存
    MC_DeleteCacheByPrefix("packageIntentionRankCache-") // 删除Intention排行缓存

    // 写入新区块
    MC_AddNewBlock(JSON.stringify(block))
    // 刷新世界状态缓存数据
    Testin.AppendNewBlockToWorldStatus(JSON.stringify(block))

    console.log("新区块写入完成：" + block.Number)
}

/*
{
    "Hash" : "e6ce8e4fcd93226d57aab1da604556e82847f5cada574a554eddbbd4618ea577",
    "PreviousHash" : "",
    "Number" : "1",
    "Transactions" : [],
    "MerkleRoot" : "",
    "Miner" : "",
    "Ts" : "",
    "Signature" : ""
}
*/