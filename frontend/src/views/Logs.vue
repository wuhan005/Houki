<template>
  <div>
    <section class="hero">
      <div class="hero-body">
        <p class="title">Logs</p>
        <p class="subtitle"></p>
      </div>
    </section>

    <b-table
        :data="logs"
        hoverable
    >
      <b-table-column v-slot="props">{{ props.row.method }}</b-table-column>
      <b-table-column v-slot="props">{{ props.row.host }}</b-table-column>
      <b-table-column v-slot="props">{{ props.row.path.slice(0, 250) }}</b-table-column>
      <b-table-column v-slot="props">{{ new Date(props.row.time * 1000).toLocaleString() }}</b-table-column>
    </b-table>
  </div>
</template>

<script>

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
  }
}
</script>

<style scoped>

</style>