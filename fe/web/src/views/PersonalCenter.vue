<template>
  <div>
    <h1>个人中心</h1>
    <div v-if='nodeID==""'>
      <v-container>
        <v-row>
          <v-col xs cols="12">
            <v-text-field v-model="privateKey" label="请输入您的私钥"></v-text-field>
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
        <v-col>
          <v-btn red @click="Logout()">
            登出
          </v-btn>
        </v-col>
      </v-row>
    </div>

  </div>
</template>

<script>
import common from '../components/common.vue'
export default {
  data : ()=>({
    nodeID : "",
    privateKey : "",
  }),

  mounted : function(){
    var privateKey = common.ls.get("privateKey")
    if(privateKey == undefined) {
      return 
    }
    var nodeID = common.secp256k1.GetNodeIDFromPrivateKey(privateKey)
    this.nodeID = nodeID
  },

  methods : {
    Login : function(){
      if (this.privateKey.length != 64) {
        alert("密钥长度不正确，请重新输入")
        return 
      }
      common.ls.set("privateKey", this.privateKey)
      // 8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6
      // nodeID : 047204499d849948aaffdec7ce2703f5b3
      var nodeID = common.secp256k1.GetNodeIDFromPrivateKey(this.privateKey)
      this.nodeID = nodeID
      console.log(nodeID)
      this.privateKey = ""
    },
    Logout : function(){
      if (confirm("确认登出并清除密钥信息吗？") == false) {
        return 
      }
      common.ls.remove("privateKey")
      this.nodeID = ""
      this.privateKey = ""
    }
  }
}
</script>