<template>
    <v-app>
        <h1>
            任务详情
        </h1>
    </v-app>
</template>

<script>
import common from '../components/common.vue'
import crypto from 'crypto-js'
export default {
    data : ()=>({
        task : {}
    }),
    mounted : function(){
        this.Refresh()
    },

    methods : {
        Refresh : function(){
            this.UpdateTaskDetail()
        },

        UpdateTaskDetail: async function(){
            var res = await common.api.GetTaskDetail(this.$route.params.task_id)
            console.log(res)
        }
    }
}
</script>