<template>
  <div style="padding-left: 24px; padding-right: 24px">
    <a-typography-title :heading="5"> Forward Proxy</a-typography-title>
    <a-form layout="inline">
      <a-form-item field="Address">
        <a-input v-model="forwardAddress"></a-input>
      </a-form-item>
      <a-form-item>
        <a-button
          type="primary"
          :loading="forwardLoading"
          @click="handleSwitchForward"
        >
          {{ status.forward.enabled ? 'Disable' : 'Enable' }}
        </a-button>
      </a-form-item>
    </a-form>
    <a-typography-title :heading="5"> Reverse Proxy</a-typography-title>
    <a-form layout="inline">
      <a-form-item field="Address">
        <a-input v-model="reverseAddress"></a-input>
      </a-form-item>
      <a-form-item>
        <a-button
          type="primary"
          :loading="reverseLoading"
          @click="handleSwitchReverse"
        >
          {{ status.reverse.enabled ? 'Disable' : 'Enable' }}
        </a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import {
    proxyStatus,
    ProxyStatusResp,
    shutdownForwardProxy,
    shutdownReverseProxy,
    startForwardProxy,
    startReverseProxy,
  } from '@/api/proxy';
  import { Message } from '@arco-design/web-vue';

  const status = ref<ProxyStatusResp>({
    forward: {
      enabled: false,
    },
    reverse: {
      enabled: false,
    },
  });
  const fetchProxyStatus = () => {
    proxyStatus().then((res) => {
      status.value = res;
    });
  };

  const forwardLoading = ref<boolean>(false);
  const forwardAddress = ref<string>('0.0.0.0:9000');
  const handleSwitchForward = () => {
    forwardLoading.value = true;

    if (status.value.forward.enabled) {
      shutdownForwardProxy()
        .then((res) => {
          Message.success(res);
        })
        .finally(() => {
          forwardLoading.value = false;
          fetchProxyStatus();
        });
    } else {
      startForwardProxy(forwardAddress.value)
        .then((res) => {
          Message.success(res);
        })
        .finally(() => {
          forwardLoading.value = false;
          fetchProxyStatus();
        });
    }
  };

  const reverseLoading = ref<boolean>(false);
  const reverseAddress = ref<string>('0.0.0.0:443');
  const handleSwitchReverse = () => {
    reverseLoading.value = true;

    if (status.value.reverse.enabled) {
      shutdownReverseProxy()
        .then((res) => {
          Message.success(res);
        })
        .finally(() => {
          reverseLoading.value = false;
          fetchProxyStatus();
        });
    } else {
      startReverseProxy(reverseAddress.value)
        .then((res) => {
          Message.success(res);
        })
        .finally(() => {
          reverseLoading.value = false;
          fetchProxyStatus();
        });
    }
  };

  onMounted(() => {
    fetchProxyStatus();
  });
</script>

<style scoped lang="less"></style>
