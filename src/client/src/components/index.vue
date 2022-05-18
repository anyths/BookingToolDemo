<template>
    <div class="kuang">
        <h1>为了羽毛球</h1>
        <t-input placeholder="请输入Key" align="center" v-model="key" clearable />
        <t-button theme="primary" style="width: 100%;border-radius: 0;margin-top: 50px;" @click="Go()">起 飞</t-button>
        <t-message v-model="msgVisible" :content="msg" theme="error" />
    </div>

</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter()

const key = ref("")
const msgVisible = ref(false)
const msg = ref("")

const Go = () => {
    key.value = key.value.trim()
    if (key.value == "") {
        msg.value = "Key不能为空! "
        msgVisible.value = true
        return
    }
    router.push(`/${key.value}`)
}

</script>

<style lang="scss">
.kuang {
    padding-top: 100px;

    >h1 {
        margin-bottom: 130px;
    }
}
</style>