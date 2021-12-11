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
                url : window.config.baseUrl + "/api/proxy",
                headers : {
                    "Content-Type" : "Application/json"
                },
                data : JSON.stringify({
                    "Params" : {
                        "MC_Call" : "RegisterHacker",
                        "name" : "SteveWooo"
                    }
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