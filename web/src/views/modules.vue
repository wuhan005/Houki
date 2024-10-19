<template>
  <a-row>
    <a-col :flex="1">
      <a-form
        :model="listModulesParams"
        :label-col-props="{ span: 6 }"
        :wrapper-col-props="{ span: 18 }"
        label-align="left"
      >
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item field="enabledOnly" label="Only Enabled">
              <a-switch v-model="listModulesParams.enabledOnly" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-col>
    <a-divider style="height: 84px" direction="vertical" />
    <a-col :flex="'86px'" style="text-align: right; margin-bottom: 18px">
      <a-space direction="vertical" :size="18">
        <a-button type="primary" @click="handleSearch" style="width: 100px">
          <template #icon>
            <icon-search />
          </template>
          Search
        </a-button>
        <a-button style="width: 100px" @click="handleReset">
          <template #icon>
            <icon-refresh />
          </template>
          Reset
        </a-button>
      </a-space>
    </a-col>
  </a-row>
  <a-divider style="margin-top: 0" />
  <a-row style="margin-bottom: 16px">
    <a-col :span="12">
      <a-space>
        <a-button type="primary" @click="handleNewModule">
          <template #icon>
            <icon-plus />
          </template>
          New Module
        </a-button>
      </a-space>
    </a-col>
    <a-col
      :span="12"
      style="display: flex; align-items: center; justify-content: end"
    >
      <a-tooltip content="Refresh">
        <div class="action-icon" @click="handleSearch">
          <icon-refresh size="18" />
        </div>
      </a-tooltip>
      <a-dropdown @select="handleSelectDensity">
        <a-tooltip content="Density">
          <div class="action-icon">
            <icon-line-height size="18" />
          </div>
        </a-tooltip>
        <template #content>
          <a-doption
            v-for="item in DENSITIES"
            :key="item.value"
            :value="item.value"
            :class="{ active: item.value === size }"
          >
            <span>{{ item.name }}</span>
          </a-doption>
        </template>
      </a-dropdown>
      <a-tooltip content="Columns">
        <a-popover
          trigger="click"
          position="bl"
          @popup-visible-change="popupVisibleChange"
        >
          <div class="action-icon">
            <icon-settings size="18" />
          </div>
          <template #content>
            <div id="tableSetting">
              <div
                v-for="(item, index) in showColumns"
                :key="item.dataIndex"
                class="setting"
              >
                <div style="margin-right: 4px; cursor: move">
                  <icon-drag-arrow />
                </div>
                <div>
                  <a-checkbox
                    v-model="item.checked"
                    @change="
                      handleChangeColumns(
                        $event,
                        item as TableColumnData,
                        index
                      )
                    "
                  >
                  </a-checkbox>
                </div>
                <div class="title">
                  {{ item.title === '#' ? '序列号' : item.title }}
                </div>
              </div>
            </div>
          </template>
        </a-popover>
      </a-tooltip>
    </a-col>
  </a-row>
  <a-table
    row-key="id"
    :loading="isLoading"
    :pagination="pagination"
    :columns="COLUMNS"
    :data="modulesData"
    :bordered="false"
    :size="size"
    @page-change="handlePageChange"
  >
  </a-table>
</template>

<script setup lang="ts">
  import { ref, onMounted, nextTick, watch } from 'vue';
  import { Pagination } from '@/types';
  import { listModules, ListModulesParams } from '@/api/modules';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import cloneDeep from 'lodash/cloneDeep';
  import Sortable from 'sortablejs';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  type Column = TableColumnData & { checked?: true };
  type SizeProps = 'mini' | 'small' | 'medium' | 'large';
  const DENSITIES = [
    { name: 'Mini', value: 'mini' },
    { name: 'Small', value: 'small' },
    { name: 'Medium', value: 'medium' },
    { name: 'Large', value: 'large' },
  ];
  const COLUMNS: TableColumnData[] = [];

  const isLoading = ref<boolean>(true);
  const size = ref<SizeProps>('medium');
  const cloneColumns = ref<Column[]>([]);
  const showColumns = ref<Column[]>([]);
  const listModulesParams = ref<ListModulesParams>({});
  const pagination = ref<Pagination>({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const modulesData = ref([]);
  const fetchTableData = () => {
    isLoading.value = true;

    listModules({
      ...listModulesParams.value,
      page: pagination.value.current,
      pageSize: pagination.value.pageSize,
    })
      .then((res) => {
        modulesData.value = res.data;
        pagination.value.total = res.total;
      })
      .finally(() => {
        isLoading.value = false;
      });
  };

  const handleSearch = () => {
    fetchTableData();
  };

  const handleReset = () => {
    listModulesParams.value = {
      enabledOnly: false,
    };
    fetchTableData();
  };

  const handlePageChange = (current: number) => {
    pagination.value.page = 1;
    pagination.value.current = current;
    fetchTableData();
  };

  const handleSelectDensity = (
    val: string | number | Record<string, any> | undefined,
    e: Event
  ) => {
    size.value = val as SizeProps;
  };

  const handleChangeColumns = (
    checked: boolean | (string | boolean | number)[],
    column: Column,
    index: number
  ) => {
    if (!checked) {
      cloneColumns.value = showColumns.value.filter(
        (item) => item.dataIndex !== column.dataIndex
      );
    } else {
      cloneColumns.value.splice(index, 0, column);
    }
  };

  const exchangeArray = <T extends Array<any>>(
    array: T,
    beforeIdx: number,
    newIdx: number,
    isDeep = false
  ): T => {
    const newArray = isDeep ? cloneDeep(array) : array;
    if (beforeIdx > -1 && newIdx > -1) {
      newArray.splice(
        beforeIdx,
        1,
        newArray.splice(newIdx, 1, newArray[beforeIdx]).pop()
      );
    }
    return newArray;
  };

  const popupVisibleChange = (val: boolean) => {
    if (val) {
      nextTick(() => {
        const el = document.getElementById('tableSetting') as HTMLElement;
        const sortable = new Sortable(el, {
          onEnd(e: any) {
            const { oldIndex, newIndex } = e;
            exchangeArray(cloneColumns.value, oldIndex, newIndex);
            exchangeArray(showColumns.value, oldIndex, newIndex);
          },
        });
      });
    }
  };

  watch(
    () => COLUMNS,
    (val) => {
      cloneColumns.value = cloneDeep(val);
      cloneColumns.value.forEach((item, index) => {
        item.checked = true;
      });
      showColumns.value = cloneDeep(cloneColumns.value);
    },
    { deep: true, immediate: true }
  );

  const handleNewModule = () => {
    router.push({ name: 'new-module' });
  };

  onMounted(() => {
    fetchTableData();
  });
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }

  :deep(.arco-table-th) {
    &:last-child {
      .arco-table-th-item-title {
        margin-left: 16px;
      }
    }
  }

  .action-icon {
    margin-left: 12px;
    cursor: pointer;
  }

  .active {
    color: #fff;
    background-color: rgb(var(--primary-5));
  }

  .setting {
    display: flex;
    align-items: center;
    width: 200px;

    .title {
      margin-left: 12px;
      cursor: pointer;
    }
  }
</style>
