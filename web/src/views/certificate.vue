<template>
  <a-row :gutter="32" style="padding-left: 24px; padding-right: 24px">
    <a-col :span="12">
      <a-typography-title :heading="5"> Certificate</a-typography-title>
      <a-form :model="formData" layout="vertical">
        <a-form-item field="certificate" label="Certificate">
          <a-textarea
            v-model="formData.certificate"
            :auto-size="true"
          ></a-textarea>
        </a-form-item>
        <a-form-item field="privateKey" label="Private Key">
          <a-textarea
            v-model="formData.privateKey"
            :auto-size="true"
          ></a-textarea>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="onUpdateCertificate">
            Update
          </a-button>
        </a-form-item>
      </a-form>
    </a-col>

    <a-col :span="12">
      <a-typography-title :heading="5">
        Certificate Metadata
      </a-typography-title>
      <a-descriptions :column="1">
        <a-descriptions-item label="Issuer">
          {{ metadata.issuer }}
        </a-descriptions-item>
        <a-descriptions-item label="Date">
          {{ dayjs(metadata.validFrom).format('YYYY-MM-DD HH:mm:ss') }} -
          {{ dayjs(metadata.validTo).format('YYYY-MM-DD HH:mm:ss') }}
        </a-descriptions-item>
        <a-descriptions-item label="Public Key Algorithm">
          {{ metadata.publicKeyAlgorithm }}
        </a-descriptions-item>
        <a-descriptions-item label="Signature Algorithm">
          {{ metadata.signatureAlgorithm }}
        </a-descriptions-item>
        <a-descriptions-item label="Serial Number">
          {{ metadata.serialNumber }}
        </a-descriptions-item>
      </a-descriptions>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import {
    CertificateMetadata,
    getCertificate,
    updateCertificate,
    UpdateCertificateData,
  } from '@/api/certificate';
  import { Message } from '@arco-design/web-vue';
  import dayjs from 'dayjs';

  const metadata = ref<CertificateMetadata>({});
  const formData = ref<UpdateCertificateData>({});
  const fetchCertificate = () => {
    getCertificate().then((res) => {
      metadata.value = res.metadata;
      formData.value = {
        certificate: res.certificate,
        privateKey: res.privateKey,
      };
    });
  };

  const onUpdateCertificate = () => {
    updateCertificate(formData.value)
      .then((res) => {
        Message.success(res);
      })
      .finally(() => {
        fetchCertificate();
      });
  };
  onMounted(() => {
    fetchCertificate();
  });
</script>

<style scoped lang="less"></style>
