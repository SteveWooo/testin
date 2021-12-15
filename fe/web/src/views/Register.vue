<template>
  <v-app>
    <div style="margin-top:20px">
      <h1>资质认证</h1>
    </div>
    <div v-if="PersonalJobs.length != 0" style="height:50px;line-height:50px;text-align:center;padding-left:20px;font-size:20px">
      <v-icon color="blue">mdi-account</v-icon> 您已拥有的身份：{{PersonalJobs.join("，")}}
    </div>

    <v-row style="margin-top:30px" v-if="PersonalJobs.length == 0">
      <v-col cols="4">
        <v-btn @click="switchDialog('hacker')">
          测试员认证
        </v-btn>
      </v-col>

      <v-col cols="4">
        <v-btn @click="switchDialog('enterprise')">
          企业注册登记
        </v-btn>
      </v-col>

      <v-col cols="4">
        <v-btn @click="switchDialog('expert')">
          专家认证
        </v-btn>
      </v-col>
    </v-row>

    <v-dialog
      v-model="hackerDialog"
      max-width="600px">
      <v-card>
        <v-card-title>
          测试员资质申请
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row dense>
              <v-col cols="12">
                <v-text-field label="姓名" v-model="forms.hacker.Name"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="简历介绍" v-model="forms.hacker.Resume"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="资历" v-model="forms.hacker.Qualification"></v-text-field>
              </v-col>
              <v-col col="8"></v-col>
              <v-col cols="2">
                <v-btn @click="switchDialog('hacker')">
                  取消
                </v-btn>
              </v-col>
              <v-col cols="2">
                <v-btn @click="DoHackerRegister()" dark>
                  提交
                </v-btn>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="enterpriseDialog"
      max-width="600px">
      <v-card>
        <v-card-title>
          企业注册登记
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row dense>
              <v-col cols="12">
                <v-text-field label="姓名" v-model="forms.enterprise.Name"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="社会统一编码" v-model="forms.enterprise.SocialCode"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="联系方式" v-model="forms.enterprise.Connection"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-row>
                  <v-col cols="4">
                    <v-file-input 
                      v-model="doms.enterprise.logoFile"
                      accept="image/*"
                      label="点击选择 Logo 文件"
                      @change="UploadEnterpriseLogo()"></v-file-input>
                  </v-col>
                  <v-col cols="8">
                    <v-text-field style="font-size:10px" disabled label="Logo IPFS 位置" v-model="forms.enterprise.LogoPath"></v-text-field>
                  </v-col>
                </v-row>

              </v-col>
              <v-col col="8"></v-col>
              <v-col cols="2">
                <v-btn @click="switchDialog('enterprise')">
                  取消
                </v-btn>
              </v-col>
              <v-col cols="2">
                <v-btn @click="DoEnterpriseRegister()" dark>
                  提交
                </v-btn>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="expertDialog"
      max-width="600px">
      <v-card>
        <v-card-title>
          专家身份注册
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row dense>
              <v-col cols="12">
                <v-text-field label="姓名" v-model="forms.expert.Name"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="简介" v-model="forms.expert.Resume"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="资历" v-model="forms.expert.Qualification"></v-text-field>
              </v-col>
              <v-col col="8"></v-col>
              <v-col cols="2">
                <v-btn @click="switchDialog('expert')">
                  取消
                </v-btn>
              </v-col>
              <v-col cols="2">
                <v-btn @click="DoExpertRegister()" dark>
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
  name : "register",
  data : ()=>({

    // 个人身份
    PersonalJobs : [],

    hackerDialog : false,
    enterpriseDialog : false,
    expertDialog : false,

    doms : {
      enterprise : {
        logoFile : undefined
      }
    },

    forms : {
      hacker : {
        Name : "",
        Resume : "",
        Qualification : "",
      },
      enterprise : {
        Name : "",
        SocialCode : "",
        LogoPath : "",
        Connection : "",
      },
      expert : {
        Name : "",
        Resume : "",
        Qualification : "",
      },
    }
  }),
  mounted : function(){
    this.Refresh()
  },

  methods : {
    initForm : function(){
      for(var f in this.forms) {
        for(var i in this.forms[f]) {
          this.forms[f][i] = ""
        }
      }
    },
    // 切换资质认证提交
    switchDialog : function(dialogName){
      this.initForm()
      if (dialogName == "hacker") {
        this.hackerDialog = !this.hackerDialog
      }

      if (dialogName == "expert") {
        this.expertDialog = !this.expertDialog
      }

      if (dialogName == "enterprise") {
        this.enterpriseDialog = !this.enterpriseDialog
      }
    },
    CloseAllDialog(){
      this.hackerDialog = false
      this.expertDialog = false
      this.enterpriseDialog = false
    },

    Refresh : async function(){
      var worldStatusRes = await common.api.GetWorldStatus()
      if (worldStatusRes.Status != 2000) {
        return 
      }
      var worldStatus = worldStatusRes.Data

      // 个人信息初始化
      this.PersonalJobs = [];
      var myNodeID = common.secp256k1.GetNodeID()
      for(let i=0;i<worldStatus.Hackers.length;i++) {
        if (myNodeID == worldStatus.Hackers[i].From) {
          this.PersonalJobs.push("测试员")
        }
      }

      for(let i=0;i<worldStatus.Enterprises.length;i++) {
        if (myNodeID == worldStatus.Enterprises[i].From) {
          this.PersonalJobs.push("企业")
        }
      }

      for(let i=0;i<worldStatus.Experts.length;i++) {
        if (myNodeID == worldStatus.Experts[i].From) {
          this.PersonalJobs.push("专家")
        }
      }
    },

    // TODO 常用函数封装起来
    // GetWorldStatus : function(){
    //   var that = this;
    //   return new Promise((resolve, reject)=>{
    //       that.axios({
    //         method : "get",
    //         url : common.BASE_URL + "/api/world_status/get",
    //     }).then(res=>{
    //         resolve(res.data)
    //     }, rejectErr=>{
    //       reject(rejectErr)
    //     })
    //   })
    // },

    UploadEnterpriseLogo : function(){
      var that = this;
      var ipfsClient = ipfs.create({
        host : common.IPFS_API_HOST,
        port : common.IPFS_API_PORT
      })

      var uploader = ipfsClient.add(this.doms.enterprise.logoFile, {
        progress : function(prog){

        }
      })
      uploader.then(function(r) {
        that.forms.enterprise.LogoPath = r.path
      })
    },

    // 统一提交接口
    submit : function(params){
      var that = this;
      this.axios({
          method : "post",
          url : common.BASE_URL + "/api/proxy",
          headers : {
              "Content-Type" : "Application/json"
          },
          data : JSON.stringify({
              "Params" : params
          })
      }).then(res=>{
        if (res.data.Status != 2000) {
          alert(res.Data.Message)
          return 
        }
        alert("您的申请已提交，共审通过后即可获得身份")
        that.CloseAllDialog()
        that.Refresh()
      }, reject=>{
          that.Refresh()
          console.log(reject)
      })
    },

    // 提交黑客资质认证材料到区块链上
    DoHackerRegister : async function(){
      var form = this.forms.hacker;
      var params = {
        "MC_Call" : "RegisterHacker" // 调用合约脚本
      }
      for (var i in form) {
        params[i] = form[i]
      }
      var now = +new Date()
      params["Ts"] = now + ""
      params.From = common.secp256k1.GetNodeID()

      // 签名
      var source = "Hacker" + params["From"] + params["Name"] + params["Qualification"] + params["Resume"] + params["Ts"]
      var hash = crypto.SHA256(source).toString()
      params.Hash = hash
      params.Signature = await common.secp256k1.Sign(hash)

      this.submit(params)
    },

    // 提交企业认证材料到区块链上
    DoEnterpriseRegister : async function(){
      var form = this.forms.enterprise;
      var params = {
        "MC_Call" : "RegisterEnterprise" // 调用合约脚本
      }
      for (var i in form) {
        params[i] = form[i]
      }
      var now = +new Date()
      params["Ts"] = now + ""
      params.From = common.secp256k1.GetNodeID()

      // 签名
      var source = "Enterprise" + params["Connection"] + params["From"] + params["LogoPath"] + params["Name"] + params["SocialCode"] + params["Ts"]
      var hash = crypto.SHA256(source).toString()
      params.Hash = hash
      params.Signature = await common.secp256k1.Sign(hash)

      this.submit(params)
    },

    DoExpertRegister : async function(){
      var form = this.forms.expert;
      var params = {
        "MC_Call" : "RegisterExpert" // 调用合约脚本
      }
      for (var i in form) {
        params[i] = form[i]
      }
      var now = +new Date()
      params["Ts"] = now + ""
      params.From = common.secp256k1.GetNodeID()

      // 签名
      var source = "Expert" + params["From"] + params["Name"] + params["Qualification"] + params["Resume"] + params["Ts"]
      var hash = crypto.SHA256(source).toString()
      params.Hash = hash
      params.Signature = await common.secp256k1.Sign(hash)
      
      this.submit(params)
    }
  }
}
</script>