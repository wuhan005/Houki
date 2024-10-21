<template>
  <a-layout class="layout">
    <div class="layout-navbar">
      <NavBar />
    </div>
    <a-layout>
      <a-layout>
        <a-layout class="layout-content">
          <a-layout-content>
            <a-card :bordered="false">
              <router-view />
            </a-card>
          </a-layout-content>
          <Footer />
        </a-layout>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" setup>
  import NavBar from '@/components/navbar.vue';
  import Footer from '@/components/footer.vue';
</script>

<style scoped lang="less">
  @nav-size-height: 60px;
  @layout-max-width: 1100px;

  .layout {
    width: 100%;
    height: 100%;
  }

  .layout-navbar {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 100;
    width: 100%;
    height: @nav-size-height;
  }

  .layout-content {
    min-height: 100vh;
    padding: 76px 20px 0;
    overflow-y: hidden;
    background-color: var(--color-fill-2);
    transition: padding 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
  }
</style>
