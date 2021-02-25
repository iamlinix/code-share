<template>
  <div id="app">
    <router-view name="blank"></router-view>
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'app',
}
</script>

<style lang="scss">
  #app {
    background-color: #222d32;
    width: 100%;
    height: 100%;
  }

  .el-upload{
    width: 100%;
  }

  html, body{
      width: 100%;
      height: 100%;
      overflow: hidden;
      margin: 0;
      padding: 0;
  }
</style>
