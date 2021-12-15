<template>
  <v-app style="width:100%;margin-left:0%;background-color:#fbfbfb">
    <div style="margin-top:20px;border:0px solid red">
      <v-row>
        <v-col cols="4"></v-col>
        <v-col cols="4" style="text-align:center">
          <h1>我发布的任务</h1>
        </v-col>
        <v-col cols="4" style="text-align:center">
          <v-btn color="primary" small @click="SwitchDialog('publishDialog')">
            新任务发布
          </v-btn>
        </v-col>
      </v-row>
    </div>

    <div style="margin-top:30px;color:#777;font-size:20px" v-if="myTasks.length == 0">
      暂无任务 ...<v-icon>mdi-bird</v-icon>
    </div>

    <div style="margin-top:20px;border:0px solid #eee;
      width:94%;margin-left:3%;border-radius:10px" v-if="myTasks.length != 0">
      <taskList :tasks="myTasks"></taskList>
      <v-row justify="center">
        <v-col cols="8">
          <v-container class="max-width">
            <v-pagination
              v-model="pageNumber"
              @input="Refresh()"
              class="my-4"
              :length="Math.ceil(taskCount / itemPerPage)"
            ></v-pagination>
          </v-container>
        </v-col>
      </v-row>
    </div>

    <v-dialog
      v-model="publishDialog"
      max-width="600px">
      <v-card>
        <v-card-title>
          任务发布
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row dense>
              <v-col cols="12">
                <v-text-field label="项目名称" v-model="forms.task.Name"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="项目基本介绍" v-model="forms.task.Resume"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="业务需求" v-model="forms.task.Require"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="项目预算" v-model="forms.task.Budget"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field label="最大可授权人数" v-model="forms.task.MaxAuthorizationCount"></v-text-field>
              </v-col>

              <v-col col="8"></v-col>
              <v-col cols="2">
                <v-btn @click="CloseAllDialog()">
                  取消
                </v-btn>
              </v-col>
              <v-col cols="2">
                <v-btn @click="DoPublishTask()" dark>
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
import taskList from '../components/taskList.vue'
import crypto from 'crypto-js'
export default {
  components : {
    taskList,
  },
  data : ()=>({
    publishDialog : false,

    myTasks : [],
    taskCount : 0,
    pageNumber : 1,
    itemPerPage : 10,

    forms : {
      task : {
        Budget : "",
        From : "",
        MaxAuthorizationCount : "",
        Name : "",
        Require : "",
        Resume : "",
        Ts : ""
      }
    }
  }),

  mounted : async function(){
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

    SwitchDialog : function(dialogName){
      this.initForm()
      this[dialogName] = !this[dialogName]
    },

    CloseAllDialog(){
      this.publishDialog = false
    },

    Refresh : async  function(){
      this.UpdateMyTask()
    },

    OpenTask (task){
      this.$router.push("/task/" + task.Hash)
    },

    // 获取任务
    UpdateMyTask : async function(){
      var res = await common.api.GetEnterprisePublishedTasks(this.pageNumber, this.itemPerPage)
      if(res.Status != 2000) {
        return 
      }

      this.taskCount = res.Data.Count;
      this.myTasks = res.Data.Tasks;
    },

    DoPublishTask : async function(){
      var form = this.forms.task;
      var params = {
        "MC_Call" : "PublishTaskByEnterprise" // 调用合约脚本
      }
      for (var i in form) {
        params[i] = form[i]
      }
      var now = +new Date()
      params["Ts"] = now + ""
      params.From = common.secp256k1.GetNodeID()

      // 签名
      var source = "Task" + params["Budget"] + params["From"] + params["MaxAuthorizationCount"] + params["Name"] + params["Require"] + params["Resume"] + params["Ts"]
      var hash = crypto.SHA256(source).toString()
      params.Hash = hash
      params.Signature = await common.secp256k1.Sign(hash)

      var res = await common.api.CallTrans(params)
      if (res.Status != 2000){
        return 
      }

      alert("提交成功")
      this.CloseAllDialog()
      this.Refresh()

    }
  }
}
</script>