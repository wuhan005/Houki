import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import utils from './utils'
import '@/assets/css/tailwind.css'

Vue.prototype.utils = utils
Vue.config.productionTip = false

new Vue({
    router,
    render: function (h) {
        return h(App)
    }
}).$mount('#app')
