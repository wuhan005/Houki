<template>
  <div class="home">
    <Nav/>
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold text-gray-900">
          Logs
        </h1>
      </div>
    </header>

    <main>
      <div class="bg-white">
        <div class="max-w-7xl mx-auto">
          <div>
            <ul class="divide-y divide-gray-200" x-max="1" v-if="logs.length !== 0">
              <li v-for="(request, index) in logs" v-bind:key="index" class="py-4">
                <div class="flex space-x-3">
                  <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                      <h3 class="text-sm font-medium">
                        <span
                            :class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium leading-4 bg-${methodColor[request.method]}-100 text-${methodColor[request.method]}-800`">
                          {{ request.method }}</span>

                        {{ request.host }}<span class="text-sm text-gray-500">{{ request.path.slice(0, 250) }}</span>
                      </h3>
                      <p class="text-sm text-gray-500">{{ new Date(request.time * 1000).toLocaleString() }}</p>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
            <div v-else>
              <p>No logs</p>
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
  name: "Logs",

  data: () => ({
    methodColor: {
      'GET': 'green',
      'POST': 'blue',
      'PUT': 'yellow',
      'DELETE': 'pink',
      'OPTIONS': 'purple',
      'CONNECT': 'gray'
    },
    logs: [],
    stream: null,
  }),

  mounted() {
    this.getStream()
  },

  beforeDestroy() {
    this.stream.close()
  },

  methods: {
    getStream() {
      this.stream = new EventSource(`${this.utils.baseURL}/logs`)
      this.stream.onmessage = (event) => {
        let data = JSON.parse(event.data)
        switch (data['type']) {
          case 'request':
            this.push(data.message)
        }
      };
      this.stream.onerror = function (err) {
        console.log(err)
      };
    },

    push(data) {
      if (this.logs.length > 13) {
        this.logs = this.logs.slice(0, 13)
      }
      this.logs.unshift(data)
    }
  },
  components: {
    Nav
  }
}
</script>

<style scoped>

</style>