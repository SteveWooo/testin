<template>
  <v-app style="width:100%;margin-left:0%;margin-top:00px;background-color:#fbfbfb">
    <div style=";border:0px solid red">
      <v-row>
        <v-img height="300px" src="../assets/banner_1.jpg"></v-img>
      </v-row>
      <v-row>
        <v-col cols="4" style="text-align:left">
          
        </v-col>
        <v-col cols="4">
          <h1>任务列表</h1>
        </v-col>
      </v-row>
    </div>

    <div style="width:94%;margin-left:3%;border-radius:10px;">
      <div style="margin-top:30px;color:#777;font-size:20px" v-if="tasks.length == 0">
        暂无任务 ...<v-icon>mdi-bird</v-icon>
      </div>

      <div style="margin-top:20px;border:0px solid #eee;background-color:#fff;
        padding:10px 10px 10px 10px" v-if="tasks.length != 0">
        <taskList :tasks="tasks"></taskList>
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
    </div>

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

    tasks : [],
    taskCount : 0,
    pageNumber : 1,
    itemPerPage : 10,

    forms : {
      
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
      this[dialogName] = !this[dialogName]
    },

    CloseAllDialog(){
      this.publishDialog = false
    },

    Refresh : async  function(){
      this.UpdateTask()
    },

    OpenTask (task){
      this.$router.push("/task/" + task.Hash)
    },

    // 获取任务
    UpdateTask : async function(){
      var res = await common.api.GetTask(this.pageNumber, this.itemPerPage)
      if(res.Status != 2000) {
        return 
      }

      this.taskCount = res.Data.Count;
      this.tasks = res.Data.Tasks;
    },
  }
}
</script>