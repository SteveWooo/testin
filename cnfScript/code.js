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

                }
                // 从缓存中提取出旧交易，把最新交易中的签名合并起来
                Trans.prototype.MergeSign = function(cacheTrans){

                }

                // 检查交易签名是否足够
                Trans.prototype.CheckSign = function(miners){

                }

                return Trans
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

        },
        Expert : {

        },
        Enterprise : {

        },
        Task : {

        }
    }
}

// 新交易。
exports.DoTransaction = function(params){
    // 1、构造交易
    var transaction = Testin.Class.Transaction.New(params)
    // 2、从缓存中获取对应的交易
    var oldTransaction = MC_GetCache("")
    // 3、合并最新交易的签名
    transaction.MergeSign(oldTransaction)

    // 4、检查签名数量
    var checkRes = transaction.CheckSign(Status.Miners)
    if(checkRes == false) {
        MC_SetCache("") // 不够的话就放回交易缓存区中
    } else {
        MC_SetCache("") // 如果足够，就放入打包缓冲区中
    }
}

// 新区块
exports.DoPackage = function(params) {

}