<template>
    <v-app v-cloak style="width:90%;margin-left:5%">
        <div style="margin-top:30px;border:0px solid red" v-if="task != undefined">
            <v-card>
                <v-card-title primary-title>
                    <v-row style="padding:12px 0 0 12px">
                        <v-col cols="2" style="border:0px solid red">
                            <v-img :src="taskCreaterEnterprise.LogoPath"
                                width="100%"></v-img>
                        </v-col>
                        <v-col cols="8">
                            <v-row dense style="text-align:left">
                                <v-col cols="12">
                                    <h2 style="font-size:25px">
                                        {{taskCreaterEnterprise.Name}}
                                    </h2>
                                </v-col>

                                <v-col cols="6" style="text-align:left;line-height:30px;color:#666;font-size:13px">
                                    <div >
                                        发布日期：{{(new Date(parseInt(task.Ts)).getFullYear()) + '/' + ((new Date(parseInt(task.Ts)).getMonth()) + 1) + '/' + (new Date(parseInt(task.Ts)).getDay())}}
                                    </div>
                                    <div>
                                        联系方式：{{taskCreaterEnterprise.Connection}}
                                    </div>
                                </v-col>
                                <v-col cols="6" style="text-align:left;color:#666;font-size:13px">
                                    <div> </div>
                                    <div>
                                        项目预算：<font style="color:red;font-weight:bold">{{task.Budget}}</font>
                                    </div>
                                    <div v-if="task.TaskHackers != undefined">
                                        已授权数：<font>{{taskPermissionCount}} / {{task.MaxAuthorizationCount}}
                                            <font style="color:grey">{{taskPermissionCount >= task.MaxAuthorizationCount?'已满':''}}</font>
                                        </font>
                                    </div>
                                </v-col>
                            </v-row>
                        </v-col>

                        <v-col cols="2" style="text-align:right">
                            <div v-if="permission.IsCreater == false && permission.IsJoinHacker == false && taskPermissionCount < task.MaxAuthorizationCount">
                                <v-btn color="primary" @click="DoApply()">立刻报名</v-btn>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-title>
            </v-card>
            
            <div style="margin-top:15px">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue">mdi-message-text</v-icon> 基本信息
                </h3>
                <v-card >
                    <v-card-text>
                        <v-row style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    项目简介：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <div>
                                    {{task.Resume}}
                                </div>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>

                <v-card style="margin-top:15px">
                    <v-card-text>
                        <v-row style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    业务需求：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <div>
                                    {{task.Require}}
                                </div>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </div>

            <!-- 已被授权的黑客报告提交模块 -->
            <div style="margin-top:15px" v-if="permission.IsPremissionHacker">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue">mdi-email</v-icon> 提交报告
                </h3>
                <v-card >
                    <v-card-text>
                        <v-row style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    报告提交：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <div>
                                    {{task.Require}}
                                </div>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </div>

            <!-- 企业授权模块 -->
            <div v-if="permission.IsCreater" style="margin-top:15px">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue">mdi-domain</v-icon> 测试授权
                </h3>
                <v-card>
                    <v-card-text>
                        <v-row dense style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    申请者列表：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <v-list style="border:0px solid red" dense>
                                    <v-list-item>
                                        <v-row dense style="text-align:left;border:0px solid red">
                                            <v-col cols="2">
                                                申请者昵称
                                            </v-col>
                                            <v-col cols="4">
                                                简介
                                            </v-col>
                                            <v-col cols="2">
                                                资历
                                            </v-col>
                                            <v-col cols="4">
                                                授权操作
                                            </v-col>
                                        </v-row>
                                    </v-list-item>
                                    <template v-for="item in task.TaskHackers">
                                        <v-list-item :key="item.From">
                                            <v-row dense style="border:0px solid red;width:100%">
                                                <v-col cols="2">
                                                    <h4>
                                                        {{item.Hacker.Name}}
                                                    </h4>
                                                </v-col>
                                                <v-col cols="4">
                                                    {{item.Hacker.Resume}}
                                                </v-col>
                                                <v-col cols="2">
                                                    {{item.Hacker.Qualification}}
                                                </v-col>
                                                <v-col cols="2" >
                                                    <v-btn small v-if="item.IsPermission=='false'" @click="DoAuth(item.Hacker.From)" color="primary">授权</v-btn>
                                                    <v-btn small v-if="item.IsPermission=='true'" disabled color="success">已授权</v-btn>
                                                </v-col>
                                                <v-col cols="2">
                                                    {{item.PermissionInformation}}
                                                </v-col>

                                                <v-col cols="12">
                                                    <v-divider :key="item.From"></v-divider>
                                                </v-col>
                                            </v-row>
                                        </v-list-item>
                                    </template>
                                </v-list>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </div>

            <!-- 报告管理模块 -->
            <div v-if="permission.IsCreater || permission.isTaskExpert" style="margin-top:15px">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue-grey">mdi-call-split</v-icon> 报告管理
                </h3>
                <v-card>
                    <v-card-text>
                        <v-row style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    测试报告：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <div>
                                    {{task.Require}}
                                </div>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </div>
            
            <div style="height:100px"></div>
        </div>
    </v-app>
</template>

<script>
import common from '../components/common.vue'
import crypto from 'crypto-js'
export default {
    data : ()=>({
        task : undefined,
        taskPermissionCount : 0,
        taskCreaterEnterprise : undefined,
        permission : undefined
    }),
    mounted : function(){
        this.Refresh()
    },

    methods : {
        Refresh : function(){
            this.UpdateTaskDetail()
        },

        UpdateTaskDetail: async function(){
            var res = await common.api.GetTaskDetail(this.$route.params.task_id)
            if (res.Status != 2000) {
                return 
            }
            console.log(res)

            if(res.Data.Task.TaskHackers == undefined) {
                res.Data.Task.TaskHackers = []
            }
            res.Data.TaskCreaterEnterprise.LogoPath = common.IPFS_GATEWAY + "/ipfs/" + res.Data.TaskCreaterEnterprise.LogoPath

            this.taskPermissionCount = 0;
            for (var i=0;i<res.Data.Task.TaskHackers.length;i++) {
                if (res.Data.Task.TaskHackers[i].IsPermission == "true") [
                    this.taskPermissionCount ++ 
                ]
            }

            this.permission = res.Data.Permission
            this.task = res.Data.Task;
            this.taskCreaterEnterprise = res.Data.TaskCreaterEnterprise
        },

        // 提交任务参与申请
        DoApply: async function(){
            if (this.taskPermissionCount >= this.task.MaxAuthorizationCount){
                alert("报名人数已到达上限，无法申请")
                return 
            }
            if (!confirm("确定提交报名吗？")) {
                return 
            }
            var params = {
                "MC_Call" : "ApplyTaskByHacker", // 调用合约脚本
                "From" : common.secp256k1.GetNodeID(),
                "TaskID" : this.task.Hash
            }
            var now = +new Date()
            params["Ts"] = now + ""

            // 签名
            var source = "TaskHacker" + params["From"] + params["TaskID"] + params["Ts"]
            var hash = crypto.SHA256(source).toString()
            params.Hash = hash
            params.Signature = await common.secp256k1.Sign(hash)

            var res = await common.api.CallTrans(params)
            if (res.Status != 2000){
                return 
            }

            alert("提交成功")
        },

        // 企业授权开发者
        DoAuth : async function(hackerID){
            var permissionInformation = prompt("请填写您的授权范围：", "")
            if (permissionInformation == null) {
                return 
            }
            if (!confirm("确定提交您的授权吗？")) {
                return 
            }

            var params = {
                "MC_Call" : "AuthorizationHackerToTaskByEnterprise", // 调用合约脚本
                "From" : common.secp256k1.GetNodeID(),
                "HackerID" : hackerID,
                "PermissionInformation" : permissionInformation,
                "TaskID" : this.task.Hash,
            }
            var now = +new Date()
            params["Ts"] = now + ""

            // 签名
            var source = "AuthorizationHackerToTaskByEnterprise" + params["From"] + params["HackerID"] + params["PermissionInformation"] + params["TaskID"] + params["Ts"]
            var hash = crypto.SHA256(source).toString()
            params.Hash = hash
            params.Signature = await common.secp256k1.Sign(hash)

            var res = await common.api.CallTrans(params)
            if (res.Status != 2000){
                return 
            }

            alert("提交成功")
        }
    }
}
</script>