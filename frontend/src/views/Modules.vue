<template>
  <div>
    <section class="hero">
      <div class="hero-body">
        <p class="title">Modules</p>
        <p class="subtitle"></p>
      </div>
    </section>

    <b-table
        :data="module.list"
        :loading="isLoading"
        :row-class="(row) => module.enabled_modules.includes(row.id) && 'enabled-row'"
    >
      <b-table-column label="Module Name" v-slot="props">{{ props.row.title }}{{ props.row.id }}</b-table-column>
      <b-table-column label="Description" v-slot="props">{{ props.row.description }}</b-table-column>
      <b-table-column label="Author" v-slot="props">{{ props.row.author }}</b-table-column>
      <b-table-column v-slot="props">
        <div class="buttons">
          <b-button type="is-warning is-light" v-if="module.enabled_modules.includes(props.row.id)"
                    @click="disableModule(props.row.id)">Disable
          </b-button>
          <b-button type="is-success is-light" v-else
                    @click="enableModule(props.row.id)">Enable
          </b-button>
          <b-button type="is-danger is-light">Delete</b-button>
        </div>
      </b-table-column>
    </b-table>
  </div>
</template>

<script>

export default {
  name: "Modules",
  data() {
    return {
      isLoading: true,
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
        this.isLoading = false
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
  }
}
</script>

<style>
tr.enabled-row {
  background: #f2effb;
}
</style>