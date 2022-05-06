

<template>

<h2>H5拉起微信支付</h2>
  <div class="kuang">
    <t-textarea v-model="url" placeholder="请粘贴脚本生成的支付链接" autosize style="min-height: 200px;"></t-textarea>
  </div>
  <div style="height:30px"></div>
  <t-button theme="primary" style="width: 100%;" @click="Go()">调起微信支付</t-button>
  <div style="height:15px"></div>
  <t-button variant="outline" style="width: 100%;" @click="Check()">验证订单状态</t-button>

  
</template>
<script setup lang="ts">
import { ref } from 'vue';

const url = ref("")

function Go() {
  let newUrl = url.value.split("@")[0]
  window.location.href = newUrl
}
function Check() {
  let checkUrl = url.value.split("@")[1]
  window.location.href = checkUrl
}
</script>
<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
.kuang {
  background-color: #eee;
  border: 1px solid #2c3e50;
  margin-top: 10px;
  min-height: 200px;
}

</style>
