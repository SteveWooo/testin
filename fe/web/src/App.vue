<template>
  <div id="app">
    <v-app-bar
      app
      fix
      color="primary"
      dense
    >
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>
        <h3>
          Testin 元众测平台
        </h3>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-app-bar>
    
    <v-navigation-drawer
      
      temporary
      absolute
      v-model="drawer"
    >
      <v-list>
        <v-list-item class="px-2">
          <v-icon>mdi-tune</v-icon>
          <v-list-item-title>
            服务列表
          </v-list-item-title>
        </v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list
        nav
        dense
      >
        <router-link
          style="text-decoration: none; color: inherit;"
          v-for="item in navList"
          :key="item.name"
          :to="item.path">
          <v-list-item link>
            <v-list-item-icon>
              <v-icon>{{item.icon}}</v-icon>
            </v-list-item-icon>
            <v-list-item-title>
              {{item.name}}
            </v-list-item-title>
          </v-list-item>
        </router-link>
      </v-list>
    </v-navigation-drawer>

    <v-main app>
      <router-view/>
    </v-main>

    <v-footer
      app
      padless
    >
      <v-card
        tile
        flat
        class="flex">
        <v-card-text class="py-2 white--text text-center">
        {{ new Date().getFullYear() }} — <strong>DeadFish</strong>
      </v-card-text>
      </v-card>
    </v-footer>
  </div>
</template>

<script>
import common from './components/common.vue'
export default {
  data : ()=>({
    drawer : false, // 左边菜单
    navList : [{
      name : "首页",
      path : "/",
      icon : "mdi-home"
    },
    {
      name : "个人中心",
      path : "/personal-center",
      icon : "mdi-account"
    },{
      name : "资质认证",
      path : "/register",
      icon : "mdi-bat"
    }
    // {
    //   name : "我的任务",
    //   path : "/my-task",
    //   icon : "mdi-star"
    // },{
    //   name : "任务溯源",
    //   path : "/task-trace",
    //   icon : "mdi-star"
    // },{
    //   name : "待评报告",
    //   path : "/task-review",
    //   icon : "mdi-star"
    // }],
  ]}),

  mounted : async function(){
    var personalJob = await common.utils.GetMyJobs()
    if (personalJob.Hacker == true || personalJob.Expert == true || personalJob.Enterprise == true) {
      this.navList.push({
        name : "任务列表",
        path : "/task-list",
        icon : "mdi-folder"
      })
    }

    if(personalJob.Enterprise == true) {
      this.navList.push({
        name : "测试任务发布",
        path : "/task-publish",
        icon : "mdi-bug"
      })
    }
  },

  methods : {

  }
}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
