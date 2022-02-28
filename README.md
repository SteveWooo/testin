# testin
基于CnF的区块链众测平台

## 启动

##### 1. 部署共识脚本到meta consensus上

共识脚本位置：/cnfScript/code.js 

1. 文件底部有创世区块代码
2. 算法名为test
3. 参与成员设置为："" （空字符串）

##### 2. 启动服务

1. 启动脚本位置：/miner/run_demo.mac.sh
```sh
cd miner
bash run_demo.mac.sh
```

2. 该脚本启动了前端服务，service服务，矿工worker

3. 入口：http://localhost:10001

##### 3. 安装、启动ipfs

启动脚本位置：/ipfs/run.sh

需要先把IPFS代码拉到 $GOPATH/src/github.com/ipfs 文件夹中

```sh
cd $GOPATH/src/github.com/ipfs
git clone https://github.com/ipfs/go-ipfs 
```

拿来玩玩的话，不开也行，就是拿来存图的。

##### 4. 清空数据

停止所有cnf节点，进入cnf项目，删除runtime下所有04开头的文件夹。他们是每个节点的levelDB数据库

```sh
cd cnf_2/runtime
rm -rf 04*
```
