<template>
  <div>
    <section class="hero">
      <div class="hero-body" v-if="proxyEnabled !== null">
        <div v-if="proxyEnabled">
          <p class="title">Proxy started.</p>
          <p class="subtitle">Listening on ...</p>
          <div class="buttons">
            <b-button type="is-danger is-light" @click="stopProxy">Stop proxy</b-button>
            <b-button type="is-info is-light" @click="downloadCA">Download CA</b-button>
            <b-button type="is-info is-light" @click="generateCA">Generate a new CA</b-button>
          </div>
        </div>
        <div v-else>
          <p class="title">Proxy stopped.</p>
          <p class="subtitle"></p>
          <div class="buttons">
            <b-button type="is-primary is-light" @click="startProxy">Start proxy</b-button>
            <b-button type="is-info is-light" @click="downloadCA">Download CA</b-button>
            <b-button type="is-info is-light" @click="generateCA">Generate a new CA</b-button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import {saveAs} from 'file-saver'

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
    },
    generateCA() {
      this.utils.POST('/proxy/ca/generate').then(res => {
        saveAs(new Blob([res]), 'houki_ca.crt')
      }).catch(err => this.$buefy.toast.open({message: err.response.data.msg, type: 'is-danger'}))
    },
    downloadCA() {
      this.utils.GET('/proxy/ca').then(res => {
        saveAs(new Blob([res]), 'houki_ca.crt')
      })
    }
  }
}
</script>
