// 导入 路由模块
import { createRouter, createWebHashHistory, createWebHistory } from "vue-router";

import indexVue from "./components/index.vue";
import keyVue from "./components/key.vue";
import idVue from "./components/id.vue";

const router = createRouter({
    history: createWebHashHistory(),
    routes :[
        {path:"/", component:indexVue},
        {path:"/:key", component:keyVue},
        {path:"/:key/:id", component:idVue},
    ]
});

export default router;