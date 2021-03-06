import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import utils from './utils'
import Buefy from 'buefy'
import 'buefy/dist/buefy.css'

Vue.prototype.utils = utils
Vue.config.productionTip = false

Vue.use(Buefy)

new Vue({
    router,
    render: function (h) {
        return h(App)
    }
}).$mount('#app')
