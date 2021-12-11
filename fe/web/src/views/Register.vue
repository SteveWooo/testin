<template>
  <v-app>
    <h1>资质认证</h1>
    <v-row>
      <v-col cols="4">
        <v-btn @click="switchDialog('hacker')">
          测试员认证
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
  </v-app>
</template>

<script>
import common from '../components/common.vue'
import crypto from 'crypto-js'
export default {
  name : "register",
  data : ()=>({
    hackerDialog : true,
    expertDialog : false,
    enterpriseDialog : false,

    forms : {
      hacker : {
        Name : "",
        Resume : "",
        Qualification : "",
      }
    }
  }),
  mounted : function(){

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
      params["Ts"] = now
      params.From = common.secp256k1.GetNodeID()

      // 签名
      var source = "Hacker" + params["From"] + params["Name"] + params["Qualification"] + params["Resume"] + now
      var hash = crypto.SHA256(source).toString()
      params.Hash = hash
      params.Signature = await common.secp256k1.Sign(hash)
      

      this.axios({
          method : "post",
          url : common.baseUrl + "/api/proxy",
          headers : {
              "Content-Type" : "Application/json"
          },
          data : JSON.stringify({
              "Params" : params
          })
      }).then(res=>{
          console.log(res.data)
      }, reject=>{
          console.log(reject)
      })
    }
  }
}
</script>