import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import CustomTheme from '@/theme.js'
import 'primeflex/primeflex.css'
import 'primeicons/primeicons.css'

import '@/assets/base.css'
import '@/assets/styles.scss'

import PrimeVue from 'primevue/config'
const app = createApp(App)

app.use(createPinia())
app.use(PrimeVue, {
    // Default theme configuration
    theme: {
        preset: CustomTheme,
        options: {
            prefix: 'c',
            darkModeSelector: 'system',
            cssLayer: false
        }
    }
})

app.use(router)

import FocusTrap from 'primevue/focustrap'
app.directive('focustrap', FocusTrap)

app.mount('#app')
