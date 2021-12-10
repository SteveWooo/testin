<template>
    <v-form>
        <v-text-field v-model="privateKey"></v-text-field>
        <v-btn small color="white" @click="test()">text</v-btn>
    </v-form>
</template>

<script>
export default {
    name : "Login",
    data : ()=>({
      privateKey : ""  
    }),
    mounted : ()=>{

    },

    methods : {
        test(){
            this.axios({
                method : "post",
                url : window.config.baseUrl + "/api/hacker/register",
                headers : {
                    "Content-Type" : "Application/json"
                },
                data : JSON.stringify({
                    "name" : "hahahahahha"
                })
            }).then(res=>{
                console.log(res)
            }, reject=>{
                console.log(reject)
            })
        }
    }
}
</script>