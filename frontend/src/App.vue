<template>
  <div id="app">
    <div class="container">
      <b-navbar>
        <template #brand>
          <b-navbar-item tag="router-link" :to="{ path: '/' }">
            <span>🧹 Houki</span>
          </b-navbar-item>
        </template>
        <template #start>
          <b-navbar-item tag="router-link" :to="{ path: '/' }">
            Dashboard
          </b-navbar-item>
          <b-navbar-item tag="router-link" :to="{ path: '/modules' }">
            Modules
          </b-navbar-item>
          <b-navbar-item tag="router-link" :to="{ path: '/store' }">
            Store
          </b-navbar-item>
          <b-navbar-item tag="router-link" :to="{ path: '/logs' }">
            Logs
          </b-navbar-item>
        </template>
      </b-navbar>
      <div class="container-fluid">
        <router-view/>
      </div>
    </div>
  </div>
</template>

<style>
#app {
  display: flex;
  height: 100%;
  justify-content: center;
}
</style>