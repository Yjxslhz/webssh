<template>
  <div id="app">
    <el-container>
      <el-header style="height: auto">
        <vheader @ssh-select="() => {
          this.$refs.tabs.openTerm()
        }"/>
      </el-header>
      <el-container>
        <el-main style="padding: 0">
            <tabs ref="tabs"></tabs>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import Header from '@/components/Header'
import Tabs from '@/components/Tabs'

export default {
    name: 'App',
    components: {
        vheader: Header,
        tabs: Tabs
    },
    mounted() {
        console.log('App mounted, loading configurations...');
        
        // 加载默认配置
        this.$store.dispatch('loadDefaultConfig').then((config) => {
            console.log('Default config loaded successfully:', config);
        }).catch(error => {
            console.error('Failed to load default config:', error);
        });
        
        // 加载服务器列表
        this.$store.dispatch('loadServerList').then((serverList) => {
            console.log('Server list loaded successfully, count:', serverList.length);
        }).catch(error => {
            console.error('Failed to load server list:', error);
        });
    }
}
</script>

<style lang="scss">
#app {
    height: 100%;
    > div {
        height: 100%;
    }
}
</style>
