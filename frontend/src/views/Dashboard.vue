<template>
  <div class="home">
    <Nav/>
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold text-gray-900">
          Dashboard
        </h1>
      </div>
    </header>

    <main>
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="bg-gray-50">
          <div v-if="proxyEnabled !== null"
               class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8 lg:flex lg:items-center lg:justify-between">
            <h2 class="text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl">
              <span v-if="proxyEnabled" class="block text-indigo-600">Proxy started.</span>
              <span v-else class="block text-black-600">Proxy stopped.</span>
              <dt v-if="proxyEnabled" class="text-sm leading-5 font-medium text-gray-500 truncate">
                Listening on ...
              </dt>
            </h2>
            <div class="mt-8 flex lg:mt-0 lg:flex-shrink-0">
              <div class="inline-flex rounded-md shadow" v-if="proxyEnabled">
                <a class="inline-flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-red-600"
                   @click="stopProxy">
                  Stop
                </a>
              </div>
              <div class="inline-flex rounded-md shadow" v-else>
                <a class="inline-flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600"
                   @click="startProxy">
                  Start
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import Nav from "../components/Nav";

export default {
  name: 'Dashboard',

  data() {
    return {
      proxyEnabled: null
    }
  },

  mounted() {
    this.getProxyStatus()
  },

  methods: {
    getProxyStatus() {
      this.utils.GET('/proxy/status').then(res => {
        this.proxyEnabled = res.enable
      })
    },
    startProxy() {
      this.utils.POST('/proxy/start', {
        'address': '0.0.0.0:8880',
      }).then(res => {
        this.getProxyStatus()
      })
    },
    stopProxy() {
      this.utils.POST('/proxy/stop').then(res => {
        this.getProxyStatus()
      })
    }
  },
  components: {
    Nav
  }
}
</script>
