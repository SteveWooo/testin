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
                                        {{task.Name}}
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
            <div style="margin-top:15px" v-if="permission.IsPremissionHacker && myTaskHackerInfo != undefined">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue">mdi-email</v-icon> 提交报告
                </h3>
                <v-card >
                    <v-card-text>
                        <v-row style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    授权信息：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                {{myTaskHackerInfo.PermissionInformation}}
                            </v-col>
                            <v-col cols="12">
                                <h3>
                                    报告情况：
                                </h3>
                            </v-col>
                            <v-col cols="12" v-if="myTaskHackerInfo.ReportPath == ''">
                                尚未提交
                            </v-col>
                            <v-col cols="12" v-if="myTaskHackerInfo.ReportPath != ''">
                                <div>
                                    最新报告地址：<a :href="common.IPFS_GATEWAY + '/ipfs/' + myTaskHackerInfo.ReportPath">点击下载</a>
                                </div>
                                <div>
                                    提交日期：{{common.utils.GetDate(myTaskHackerInfo.Ts)}}
                                </div>
                            </v-col>
                            <v-col cols="12">
                                <v-btn small color="primary" @click="SwitchDialog('publishHackerReportDialog')">提交/修改报告</v-btn>
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
                                        <v-row dense style="text-align:left;border:0px solid red;font-weight:bold">
                                            <v-col cols="2">
                                                申请者昵称
                                            </v-col>
                                            <v-col cols="4">
                                                简介
                                            </v-col>
                                            <v-col cols="2">
                                                资历
                                            </v-col>
                                            <v-col cols="2">
                                                授权操作
                                            </v-col>
                                            <v-col cols="2">
                                                授权信息
                                            </v-col>
                                        </v-row>
                                    </v-list-item>
                                    <template v-for="item in task.TaskHackers">
                                        <v-list-item :key="item.From">
                                            <v-row dense style="border:0px solid red;width:100%;">
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

            <!-- 企业报告管理模块 -->
            <div v-if="permission.IsCreater" style="margin-top:15px">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue-grey">mdi-call-split</v-icon> 报告管理
                </h3>
                <v-card>
                    <v-card-text>
                        <v-row dense style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    报告总表：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <v-list style="border:0px solid red" dense>
                                    <v-list-item>
                                        <v-row dense style="text-align:left;border:0px solid red;font-weight:bold">
                                            <v-col cols="2">
                                                测试员昵称
                                            </v-col>
                                            <v-col cols="2">
                                                资历
                                            </v-col>
                                            <v-col cols="4">
                                                报告情况
                                            </v-col>
                                            <v-col cols="2">
                                                评审情况
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
                                                <v-col cols="2">
                                                    {{item.Hacker.Qualification}}
                                                </v-col>
                                                <v-col cols="4" >
                                                    <div v-if="item.ReportPath != ''">
                                                        <div>
                                                            最新报告地址：<a :href="common.IPFS_GATEWAY + '/ipfs/' + item.ReportPath">点击下载</a>
                                                        </div>
                                                        <div>
                                                            提交日期：{{common.utils.GetDate(item.Ts)}}
                                                        </div>
                                                    </div>
                                                    <div v-if="item.ReportPath == ''">
                                                        <font style="color:#777">
                                                            报告未提交
                                                        </font>
                                                    </div>
                                                </v-col>
                                                <v-col cols="4">
                                                    <v-row dense>
                                                        <v-col cols="12" v-for="review in item.ExpertReviewReports" :key="review.Hash">
                                                            {{review.Score}} : {{review.Memo}}
                                                        </v-col>
                                                    </v-row>
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

            <!-- 专家报告评审 -->
            <div v-if="permission.IsTaskExpert" style="margin-top:15px">
                <h3 style="text-align:left;line-height:50px;color:#666;font-size:20px">
                    <v-icon color="blue-grey">mdi-dog</v-icon> 待审报告
                </h3>
                <v-card>
                    <v-card-text>
                        <v-row dense style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    报告总表：
                                </h3>
                            </v-col>
                            <v-col cols="12">
                                <v-list style="border:0px solid red" dense>
                                    <v-list-item>
                                        <v-row dense style="text-align:left;border:0px solid red;font-weight:bold">
                                            <v-col cols="2">
                                                测试员昵称
                                            </v-col>
                                            <v-col cols="2">
                                                资历
                                            </v-col>
                                            <v-col cols="4">
                                                报告情况
                                            </v-col>
                                            <v-col cols="2">
                                                评审情况
                                            </v-col>
                                        </v-row>
                                    </v-list-item>
                                    <template v-for="item in expertNeedToReview">
                                        <v-list-item :key="item.From">
                                            <v-row dense style="border:0px solid red;width:100%">
                                                <v-col cols="2">
                                                    <h4>
                                                        {{item.Hacker.Name}}
                                                    </h4>
                                                </v-col>
                                                <v-col cols="2">
                                                    {{item.Hacker.Qualification}}
                                                </v-col>
                                                <v-col cols="4" >
                                                    <div v-if="item.ReportPath != ''">
                                                        <div>
                                                            最新报告地址：<a :href="common.IPFS_GATEWAY + '/ipfs/' + item.ReportPath">点击下载</a>
                                                        </div>
                                                        <div>
                                                            提交日期：{{common.utils.GetDate(item.Ts)}}
                                                        </div>
                                                    </div>
                                                </v-col>
                                                <v-col cols="2">
                                                    <v-btn small color="primary" @click="
                                                        SwitchDialog('expertReviewDialog');
                                                        forms.expertReview.HackerID=item.From;
                                                        forms.expertReview.TaskID=item.TaskID;">
                                                        评审
                                                    </v-btn>
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

            <div style="height:100px"></div>
        </div>

        <v-dialog
            v-model="publishHackerReportDialog"
            max-width="600px">
            <v-card>
                <v-card-title>
                任务发布
                </v-card-title>

                <v-card-text>
                <v-container>
                    <v-row dense>
                    <v-col cols="12">
                        <v-text-field disabled label="文件有效IPFS编码（选择文件后生成）：" v-model="forms.taskHackerReport.ReportPath"></v-text-field>
                        <v-file-input 
                            v-model="doms.hacker.reportFile"
                            label="点击选择报告文件"
                            @change="UploadReportFile()"></v-file-input>
                    </v-col>

                    <v-col col="8"></v-col>
                    <v-col cols="2">
                        <v-btn @click="SwitchDialog('publishHackerReportDialog')">
                        取消
                        </v-btn>
                    </v-col>
                    <v-col cols="2">
                        <v-btn @click="DoPublishTaskHackerReport()" dark>
                        提交
                        </v-btn>
                    </v-col>
                    </v-row>
                </v-container>
                </v-card-text>
            </v-card>
        </v-dialog>

        <v-dialog
            v-model="expertReviewDialog"
            max-width="600px">
            <v-card>
                <v-card-title>
                提交评审意见
                </v-card-title>

                <v-card-text>
                    <v-container>
                        <v-row dense>
                            <v-col cols="6">
                                <v-text-field disabled label="测试员ID：" v-model="forms.expertReview.HackerID"></v-text-field>
                            </v-col>
                            <v-col cols="6">
                                <v-text-field disabled label="任务ID：" v-model="forms.expertReview.TaskID"></v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-text-field label="评分（0到100）：" v-model="forms.expertReview.Score"></v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-text-field label="评语：" v-model="forms.expertReview.Memo"></v-text-field>
                            </v-col>

                            <v-col col="8"></v-col>
                            <v-col cols="2">
                                <v-btn @click="SwitchDialog('expertReviewDialog')">
                                取消
                                </v-btn>
                            </v-col>
                            <v-col cols="2">
                                <v-btn @click="DoPublishExpertReview()" dark>
                                提交
                                </v-btn>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-card-text>
            </v-card>
        </v-dialog>
    </v-app>
</template>

<script>
import common from '../components/common.vue'
import crypto from 'crypto-js'
import * as ipfs from 'ipfs-http-client'
export default {
    data : ()=>({
        task : undefined,
        myTaskHackerInfo : undefined, // 作为测试员，需要看到自己的测试报告的信息
        expertNeedToReview : [], // 专家待评审的任务
        taskPermissionCount : 0, // 已经授权的测试员数量
        taskCreaterEnterprise : undefined, // 创建任务的企业信息
        permission : undefined, // 打开这歌页面的用户能获取到的情报

        publishHackerReportDialog : false,
        expertReviewDialog : false,

        forms : {
            taskHackerReport : {
                ReportPath : ""
            },
            expertReview : {
                HackerID : "",
                TaskID : "",
                Score : "",
                Memo : "",
            }
        },

        doms : {
            hacker : {
                reportFile : undefined
            }
        },

        common : common
    }),
    mounted : function(){
        this.Refresh()
    },

    methods : {
        Refresh : function(){
            this.UpdateTaskDetail()
        },

        initForm : function(){
            for(let f in this.forms) {
                for(let i in this.forms[f]) {
                    this.forms[f][i] = ""
                }
            }
            for(let d in this.doms) {
                for(let i in this.doms[d]) {
                    this.doms[d][i] = undefined
                }
            }
        },

        SwitchDialog : function(dialogName){
            this.initForm()
            this[dialogName] = !this[dialogName]
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
            this.myTaskHackerInfo = res.Data.MyTaskHackerInfo
            this.expertNeedToReview = res.Data.ExpertNeedToReview
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
        },
        // 上传报告文件到IPFS，然后获取返回的连接
        UploadReportFile : async function(){
            var that = this;
            var ipfsClient = ipfs.create({
                host : common.IPFS_API_HOST,
                port : common.IPFS_API_PORT
            })

            var uploader = ipfsClient.add(this.doms.hacker.reportFile, {
                progress : function(prog){

                }
            })
            uploader.then(function(r) {
                that.forms.taskHackerReport.ReportPath = r.path
            })
        },

        // 测试员提交测试报告
        DoPublishTaskHackerReport : async function(){
            var form = this.forms.taskHackerReport;
            if (form.ReportPath == "") {
                alert("文件上传地址获取失败");
                return 
            }

            var params = {
                "MC_Call" : "PublishReportByHacker", // 调用合约脚本
                "From" : common.secp256k1.GetNodeID(),
                "TaskID" : this.myTaskHackerInfo.TaskID,
                "ReportPath" : form.ReportPath
            }
            var now = +new Date()
            params["Ts"] = now + ""

            // 签名
            var source = "TaskHackerReport" + params["From"] + params["ReportPath"] + params["TaskID"] + params["Ts"]
            var hash = crypto.SHA256(source).toString()
            params.Hash = hash
            params.Signature = await common.secp256k1.Sign(hash)

            var res = await common.api.CallTrans(params)
            if (res.Status != 2000){
                return 
            }

            alert("提交成功")
            this.SwitchDialog("publishHackerReportDialog")
        },

        // 专家提交评审报告
        DoPublishExpertReview : async function(){
            var form = this.forms.expertReview;
            if (parseInt(form.Score).toString() == "NaN" || parseInt(form.Score) < 0 || parseInt(form.Score) > 100) {
                alert("分数参数填写错误，合法范围：0～100");
                return 
            }

            var params = {
                "MC_Call" : "ReviewReportByExpert", // 调用合约脚本
                "From" : common.secp256k1.GetNodeID(), // 专家id
                "TaskID" : form.TaskID,
                "HackerID" : form.HackerID,
                "Score" : form.Score,
                "Memo" : form.Memo
            }
            var now = +new Date()
            params["Ts"] = now + ""

            // 签名
            var source = "ExpertReviewReport" + params["From"] + params["HackerID"] + params["Memo"] + params["Score"] + params["TaskID"] + params["Ts"]
            var hash = crypto.SHA256(source).toString()
            params.Hash = hash
            params.Signature = await common.secp256k1.Sign(hash)

            var res = await common.api.CallTrans(params)
            if (res.Status != 2000){
                return 
            }

            alert("提交成功")
            this.SwitchDialog("expertReviewDialog")
        }

    }
}
</script>