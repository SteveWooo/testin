import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/register',
    name : 'register',
    component : ()=>import ('../views/Register.vue')
  },
  {
    path: '/task-publish',
    name : 'task-publish',
    component : ()=>import ('../views/TaskPublish.vue')
  },
  {
    path: '/task-list',
    name : 'task-list',
    component : ()=>import ('../views/TaskList.vue')
  },
  {
    path: '/my-task',
    name : 'my-task',
    component : ()=>import ('../views/MyTask.vue')
  },
  {
    path: '/task-trace',
    name : 'task-trace',
    component : ()=>import ('../views/TaskTrace.vue')
  },
  {
    path: '/task-review',
    name : 'task-review',
    component : ()=>import ('../views/TaskReview.vue')
  },
  {
    path: '/personal-center',
    name : 'personal-center',
    component : ()=>import ('../views/PersonalCenter.vue')
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
