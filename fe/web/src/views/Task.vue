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
                                        已报名数：<font>{{task.TaskHackers.length}} / {{task.MaxAuthorizationCount}}</font>
                                    </div>
                                </v-col>
                            </v-row>
                        </v-col>

                        <v-col cols="2" style="text-align:right">
                            <div v-if="permission.IsCreater == false && permission.IsJoinHacker == false">
                                <v-btn>立刻报名</v-btn>
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
                        <v-row style="text-align:left">
                            <v-col cols="12">
                                <h3>
                                    申请者列表：
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

            this.permission = res.Data.Permission
            this.task = res.Data.Task;
            this.taskCreaterEnterprise = res.Data.TaskCreaterEnterprise
        }
    }
}
</script>