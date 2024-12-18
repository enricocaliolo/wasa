import './shared/assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './shared/router'
import LoadingSpinner from './shared/components/LoadingSpinner.vue'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.component('LoadingSpinner', LoadingSpinner)

app.mount('#app')
