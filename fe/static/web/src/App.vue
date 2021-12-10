<template>
  <v-app>
    <v-app-bar app color="white">
      <v-app-bar-nav-icon @click.stop="drawer=!drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>
        <h3>
          众测平台
        </h3>
      </v-toolbar-title>

      <v-spacer></v-spacer>
    </v-app-bar>

    <v-navigation-drawer 
    app 
    v-model="drawer">
    </v-navigation-drawer>

    <v-main>
        <v-container>
          <Main/>
          <Login/>
        </v-container>
    </v-main>

    <v-footer app></v-footer>
  </v-app>
</template>

<script>

import Main from './components/Main'
import Login from './components/Login'
export default {
  name: 'App',
  components: {
    Main,
    Login,
  },

  data: function(){
    return {
      drawer : false
    }
  },

  mounted : function(){
    window.config = {
      baseUrl : "http://127.0.0.1:10001"
    }
  }
};
</script>
