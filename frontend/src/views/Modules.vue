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

                  <tr v-for="(module, index) in modules" v-bind:key="index">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div>
                          <div class="text-sm font-medium text-gray-900">{{ module.title }}</div>
                          <div class="text-sm text-gray-500">{{ module.id }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm text-gray-900">{{ module.description }}</div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ module.author }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                          class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        Active
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <a class="text-green-600 hover:text-indigo-900">Enable</a>
                      <a class="text-indigo-600 hover:text-indigo-900">Disable</a>
                      <a class="text-red-600 hover:text-indigo-900">Delete</a>
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
      modules: []
    }
  },

  mounted() {
    this.getModuleList()
  },

  methods: {
    getModuleList() {
      this.utils.GET('/modules').then(res => {
        this.modules = res
      }).catch((err) => {

      })
    }
  },
  components: {Nav}
}
</script>

<style scoped>

</style>