<template>
  <v-app style="width:90%;margin-left:5%">
    <div style="margin-top:20px">
        <h1>个人中心</h1>
    </div>
    <div v-if='nodeID==""'>
      <v-container>
        <v-row>
          <v-col xs cols="12">
            <v-text-field v-model="privateKey" label="请输入您的私钥"></v-text-field>
            8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6
          </v-col>
          <v-col cols="6" xs>
          </v-col>
          <v-col xs cols="6" style="text-align:center">
            <v-btn dark
              @click=Login()>
              加密录入
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </div>

    <div v-if='nodeID!=""' style="margin-top:30px">
      <v-row>
        <v-col cols="12">
          已登陆NodeID：{{nodeID}}
        </v-col>
        <v-col cols="12" v-if="PersonalJobs.length != 0">
          您的身份：{{PersonalJobs.join(", ")}}
        </v-col>
        <v-col cols="12">
          <v-btn color="error" @click="Logout()">
            登出
          </v-btn>
        </v-col>
      </v-row>
    </div>

  </v-app>
</template>

<script>
import common from '../components/common.vue'
export default {
  data : ()=>({
    nodeID : "",
    privateKey : "",
    PersonalJobs : []
  }),

  mounted : function(){
    var privateKey = common.ls.get("privateKey")
    if(privateKey == undefined) {
      return 
    }
    var nodeID = common.secp256k1.GetNodeIDFromPrivateKey(privateKey)
    this.nodeID = nodeID
    this.Refresh()
  },

  methods : {
    Refresh : function(){
      this.UpdatePersonalJob()
    },
    UpdatePersonalJob : async function(){
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
    Login : function(){
      if (this.privateKey.length != 64) {
        alert("密钥长度不正确，请重新输入")
        return 
      }
      common.ls.set("privateKey", this.privateKey)
      // 8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6
      // nodeID : 047204499d849948aaffdec7ce2703f5b3
      var nodeID = common.secp256k1.GetNodeIDFromPrivateKey(this.privateKey)
      // this.nodeID = nodeID
      // this.privateKey = ""
      // this.Refresh()
      window.location.reload()
    },
    Logout : function(){
      if (confirm("确认登出并清除密钥信息吗？") == false) {
        return 
      }
      common.ls.remove("privateKey")
      this.nodeID = ""
      this.privateKey = ""
      this.PersonalJobs = []
    }
  }
}
</script>