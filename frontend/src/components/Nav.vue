<template>
  <nav class="bg-gray-800">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <span style="font-size: 32px">🧹</span>
          </div>
          <div class="md:block">
            <div class="ml-10 flex items-baseline space-x-4">
              <!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
              <router-link to="/">
                <a :class="activeCSS('Dashboard')"
                   class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Dashboard</a>
              </router-link>
              <router-link to="/modules">
                <a :class="activeCSS('Modules')"
                   class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Modules</a>
              </router-link>
              <router-link to="/store">
                <a :class="activeCSS('Store')"
                   class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Store</a>
              </router-link>
              <router-link to="/logs">
                <a :class="activeCSS('Logs')"
                   class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Logs</a>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: "Nav",

  methods: {
    activeCSS: function (name) {
      if (this.$route.name === name) {
        return 'bg-gray-900 text-white'
      }
    }
  }
}
</script>

<style scoped>

</style>