<template>
  <div id="app">
    <v-app-bar
      app
      dark
      dense
    >
      <v-toolbar-title>
        众测平台
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-app-bar>
    
    <v-navigation-drawer
      app
      permanent
      expand-on-hover
    >
      <v-list>
        <v-list-item class="px-2">
          <v-icon>mdi-account</v-icon>
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
    navList : [{
      name : "首页",
      path : "/",
      icon : "mdi-star"
    },{
      name : "资质认证",
      path : "/register",
      icon : "mdi-star"
    },{
      name : "测试任务发布",
      path : "/task-publish",
      icon : "mdi-folder"
    },{
      name : "任务列表",
      path : "/task-list",
      icon : "mdi-star"
    },{
      name : "我的任务",
      path : "/my-task",
      icon : "mdi-star"
    },{
      name : "任务溯源",
      path : "/task-trace",
      icon : "mdi-star"
    },{
      name : "待评报告",
      path : "/task-review",
      icon : "mdi-star"
    },{
      name : "个人中心",
      path : "/personal-center",
      icon : "mdi-star"
    }],

  }),

  mounted : function(){
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
