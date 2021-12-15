<template >
    <v-list >
        <v-list-item :key="item.Hash" v-for="item in tasks">
            <v-list-item-content>
                <v-card style="padding:12px 12px 12px 12px">
                <v-row dense>
                    <v-col cols="10" style="text-align:left">
                    <h3>
                        {{item.Name}}
                    </h3>
                    </v-col>
                    <v-col cols="2">
                    <v-btn color="primary" @click="$router.push('/task/' + item.Hash)">打开任务</v-btn>
                    </v-col>

                    <v-col cols="12" style="text-align:left;font-size:14px;color:#777">
                    <div>
                        业务需求：{{item.Require}} | 发布日期：{{common.utils.GetDate(item.Ts)}}
                    </div>
                    </v-col>
                    <v-col cols="12" style="text-align:left;color:#777;font-size:14px">
                        项目介绍：{{item.Resume}}
                    </v-col>

                    <v-col cols="12" style="text-align:right;color:#999;font-size:10px">
                        发布者：{{item.From}}
                    </v-col>
                </v-row>
                </v-card>
            </v-list-item-content>
        </v-list-item>
    </v-list>
</template>

<script>
import common from '../components/common.vue'
export default {
    name : "taskList",
    data : ()=>({
        common: common
    }),
    props : {
        tasks : Array
    },
    mounted : function(){
        // console.log(this.tasks)
    }
}
</script>