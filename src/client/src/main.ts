import { createApp } from 'vue'
import App from './App.vue'
import router from './routes';

import TDesign from 'tdesign-mobile-vue';

createApp(App)
    .use(router)
    .use(TDesign)
    .mount('#app')
