// 整体状态的数据结构
var Status = {
    // 指定挖矿账号（后续应该由代表专家来做）
    Miners : [],

    Hackers : [],
    Experts : [],
    Enterprises : [],
    Tasks : [],
}

// 类与方法
var Testin = {
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

        Block : {
            New : function(params){
                function Block(_params) {

                }

                return Block
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
        console.log("签名校验失败: RegisterHacker")
        return 
    }

    // 加上类型，就是
    var transParam = {
        Type : "RegisterHacker",
        Hacker : hacker
    }
    var trans = Testin.Class.Transaction.New(transParam)
    console.log(trans.CheckSign())
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
exports.DoPackage = function(params) {

}