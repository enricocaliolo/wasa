import {createApp} from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './shared/router/router.js';
import BaseModal from './shared/components/BaseModal.vue';

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.component('BaseModal', BaseModal)

app.mount('#app')
