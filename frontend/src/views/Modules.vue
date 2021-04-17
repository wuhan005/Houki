<template>
  <div>
    <Nav/>
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold text-gray-900">
          Modules
        </h1>
      </div>
    </header>

    <main>
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex flex-col">
          <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
              <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-50">
                  <tr>
                    <th scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Module Name
                    </th>
                    <th scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Description
                    </th>
                    <th scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Author
                    </th>
                    <th scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status
                    </th>
                    <th scope="col" class="relative px-6 py-3"></th>
                  </tr>
                  </thead>

                  <tbody class="bg-white divide-y divide-gray-200">

                  <tr v-for="(mod, index) in module.list" v-bind:key="index">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div>
                          <div class="text-sm font-medium text-gray-900">{{ mod.title }}</div>
                          <div class="text-sm text-gray-500">{{ mod.id }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm text-gray-900">{{ mod.description }}</div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ mod.author }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span v-if="module.enabled_modules.includes(mod.id)"
                            class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        Active
                      </span>
                      <span v-else
                            class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800">
                        Inactive
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <span class="relative z-0 inline-flex shadow-sm rounded-md">
                        <button v-if="!module.enabled_modules.includes(mod.id)"
                                @click="enableModule(mod.id)"
                                type="button"
                                class="text-green-600 hover:text-indigo-900 relative inline-flex items-center px-4 py-2 rounded-l-md border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-700 hover:text-gray-500 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-700 transition ease-in-out duration-150">
                          Enable
                        </button>
                        <button v-else
                                @click="disableModule(mod.id)"
                                type="button"
                                class="text-indigo-600 hover:text-indigo-900 -ml-px relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-700 hover:text-gray-500 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-700 transition ease-in-out duration-150">
                          Disable
                        </button>
                        <button type="button"
                                class="text-red-600 hover:text-indigo-900 -ml-px relative inline-flex items-center px-4 py-2 rounded-r-md border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-700 hover:text-gray-500 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-700 transition ease-in-out duration-150">
                          Delete
                        </button>
                      </span>
                    </td>
                  </tr>

                  </tbody>
                </table>
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
  name: "Modules",
  data() {
    return {
      module: []
    }
  },

  mounted() {
    this.getModuleList()
  },

  methods: {
    getModuleList() {
      this.utils.GET('/modules').then(res => {
        this.module = res
      }).catch((err) => {

      })
    },

    enableModule(modID) {
      this.utils.POST(`/module/enable/${modID}`).then(res => {
        this.getModuleList()
      })
    },
    disableModule(modID) {
      this.utils.POST(`/module/disable/${modID}`).then(res => {
        this.getModuleList()
      })
    }
  },
  components: {Nav}
}
</script>

<style scoped>

</style>