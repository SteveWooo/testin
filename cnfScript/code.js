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
    // 3、如果缓存中没有这个交易，则检查交易合法性

    // 4、把合法的交易存入缓存中，等待矿工拉取并打包

}

// 新区块
exports.DoPackage = function(params) {

}