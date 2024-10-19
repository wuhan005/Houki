<template>
  <div class="navbar">
    <div class="left-side">
      <a-space>
        <a-typography-title
          :style="{ margin: 0, fontSize: '18px' }"
          :heading="5"
        >
          🧹 Houki
        </a-typography-title>
      </a-space>
    </div>
    <div class="center-side">
      <a-menu
        mode="horizontal"
        :selected-keys="[route.name]"
        @menu-item-click="handleMenuItemClick"
      >
        <a-menu-item key="modules">Modules</a-menu-item>
        <a-menu-item key="proxy">Proxy</a-menu-item>
        <a-menu-item key="certificate">Certificate</a-menu-item>
      </a-menu>
    </div>
    <ul class="right-side">
      <li>
        <a-button
          class="nav-btn"
          type="outline"
          shape="circle"
          @click="handleToggleTheme"
        >
          <template #icon>
            <icon-moon-fill v-if="theme === 'dark'" />
            <icon-sun-fill v-else />
          </template>
        </a-button>
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import { useDark, useToggle } from '@vueuse/core';
  import { useAppStore } from '@/store';
  import { useRoute, useRouter } from 'vue-router';

  const route = useRoute();
  const router = useRouter();
  const appStore = useAppStore();
  const theme = computed(() => {
    return appStore.theme;
  });

  const isDark = useDark({
    selector: 'body',
    attribute: 'arco-theme',
    valueDark: 'dark',
    valueLight: 'light',
    storageKey: 'arco-theme',
    onChanged(dark: boolean) {
      appStore.toggleTheme(dark);
    },
  });
  const toggleTheme = useToggle(isDark);
  const handleToggleTheme = () => {
    toggleTheme();
  };

  const handleMenuItemClick = (item: string) => {
    router.push({ name: item });
  };
</script>

<style scoped lang="less">
  .navbar {
    display: flex;
    justify-content: space-between;
    height: 100%;
    background-color: var(--color-bg-2);
    border-bottom: 1px solid var(--color-border);
  }

  .left-side {
    display: flex;
    align-items: center;
    padding-left: 20px;
  }

  .center-side {
    flex: 1;
  }

  .right-side {
    display: flex;
    padding-right: 20px;
    list-style: none;

    :deep(.locale-select) {
      border-radius: 20px;
    }

    li {
      display: flex;
      align-items: center;
      padding: 0 10px;
    }

    a {
      color: var(--color-text-1);
      text-decoration: none;
    }

    .nav-btn {
      border-color: rgb(var(--gray-2));
      color: rgb(var(--gray-8));
      font-size: 16px;
    }

    .trigger-btn,
    .ref-btn {
      position: absolute;
      bottom: 14px;
    }

    .trigger-btn {
      margin-left: 14px;
    }
  }
</style>

<style lang="less">
  .message-popover {
    .arco-popover-content {
      margin-top: 0;
    }
  }
</style>
