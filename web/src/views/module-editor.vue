<template>
  <a-page-header
    :style="{ background: 'var(--color-bg-2)' }"
    :title="title"
    :subtitle="subTitle"
    :show-back="false"
  >
    <template #extra>
      <a-button type="primary" @click="handleSave">
        <template #icon>
          <icon-save />
        </template>
        Save
      </a-button>
    </template>
  </a-page-header>
  <a-form
    :model="formData"
    :label-col-props="{ span: 1 }"
    :wrapper-col-props="{ span: 22 }"
    style="margin: 16px"
  >
    <a-form-item field="name" label="Name">
      <a-input v-model="formData.name" style="width: 300px"></a-input>
    </a-form-item>
    <a-form-item field="body" label="Body">
      <MonacoEditor
        v-model:value="editorRaw"
        style="height: 500px; width: 100%"
        language="json"
        :theme="editorTheme"
        :options="editorOptions"
      ></MonacoEditor>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import MonacoEditor from 'monaco-editor-vue3';
  import { useAppStore } from '@/store';
  import { getModule, createModule, updateModule } from '@/api/modules';
  import type { CreateModuleData, UpdateModuleData } from '@/api/modules';
  import { Message } from '@arco-design/web-vue';

  const route = useRoute();
  const router = useRouter();
  const appStore = useAppStore();

  const mode = route.name === 'new-module' ? 'new' : 'update';
  const moduleID = ref<string>(route.params.id as string);
  const title = ref<string>(mode === 'new' ? 'New Module' : 'Update Module');
  const subTitle = ref<string>('');
  const editorRaw = ref<string>('');
  const editorTheme = computed(() =>
    appStore.theme === 'dark' ? 'vs-dark' : 'vs'
  );
  const editorOptions = {
    minimap: {
      enabled: false,
    },
  };

  const formData = ref<CreateModuleData | UpdateModuleData>({} as CreateModuleData);

  const fetchModule = () => {
    getModule(moduleID.value).then((res) => {
      subTitle.value = res.name;
      formData.value = {
        name: res.name,
        body: res.body,
      };
      editorRaw.value = JSON.stringify(res.body, null, 2);
    });
  };

  const handleSave = () => {
    try {
      formData.value.body = JSON.parse(editorRaw.value);
    } catch (e: any) {
      Message.error(`Failed to parse JSON: ${e.message}`);
      return;
    }

    if (mode === 'new') {
      createModule(formData.value as CreateModuleData).then((res) => {
        const url = router.resolve({
          name: 'update-module',
          params: {
            id: res.id,
          },
        });
        window.open(url.href, '_self');
      });
    } else {
      updateModule(moduleID.value, formData.value as UpdateModuleData)
        .then((res) => {
          Message.success(res);
        })
        .finally(() => {
          fetchModule();
        });
    }
  };

  onMounted(() => {
    if (mode === 'update') {
      fetchModule();
    }
  });
</script>

<style scoped lang="less"></style>
