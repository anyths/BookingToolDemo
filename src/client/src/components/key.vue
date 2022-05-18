<template>
    <t-message v-model="msgVisible" :content="msg" theme="error" />
    <t-message v-model="okVisible" :content="okmsg" theme="success" />
    <t-navbar> [ {{ key }} ] </t-navbar>
    <template v-if="statusList">
        <TransitionGroup tag="ul" name="fade" class="container">
            <div v-for="(index, key) in statusList" class="item" :key="key" @click="Open(index.id)">
                <div class="hangpre"></div>
                <div class="hang">

                    <div class="lie20"># {{ key }}</div>
                    <div class="lie70">
                        <div style="display: inline-block;">
                            <TimeFilledIcon v-if="index.status" class="success" size="large" />
                            <LoadingIcon v-else class="running" size="large" />
                        </div><span style="padding-left:10px;padding-right:40px;"
                            :class="index.status ? 'success' : 'running'"> {{
                                    index.status
                                        ?
                                        '空闲' :
                                        '执行中...'
                            }}</span>
                        <NextIcon size="large" style="margin-right:10px" />{{ index.cmd ? index.cmd : "无" }}
                    </div>
                </div>
                <div class="hang">
                    <div class="lie50">项目: {{ index.sport_events_id }}</div>
                    <div class="lie50">场地: {{ index.venue_id }}</div>
                </div>
                <div class="hang">
                    <div class="lie60">时段: {{ index.period }}</div>
                    <div class="lie40">价格: {{ index.money }}</div>
                </div>
                <div class="hang">
                    <div class="lie100">
                        订单状态: {{ index.form ? "已生成" : "未生成" }}
                    </div>
                </div>
                <!-- <div class="hang">
                    有效时间：
                    <t-countdown content="ok" :time="new Date(index.oct).getTime() + 900000 - new Date().getTime()"
                        :auto-start="true" :millisecond="false" format="mm:ss:sss" theme="square" @finish="onFinish" />
                </div> -->
                <div class="hangpre"></div>

            </div>
        </TransitionGroup>
    </template>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Status, api } from '../interface';
import { TimeFilledIcon, LoadingIcon, NextIcon } from 'tdesign-icons-vue-next'
const route = useRoute()
const router = useRouter()


const msgVisible = ref(false)
const okVisible = ref(false)
const msg = ref("")
const okmsg = ref("")
const key = ref<string | string[]>("")
key.value = route.params.key




const statusList = ref<Status[] | null>(null)

GetList()
var timer = setInterval(() => {
    GetList()
}, 1000)

onBeforeUnmount(() => {
    clearInterval(timer)
})

async function GetList() {
    try {
        const response = await fetch(`${api.sta}${key.value}`);
        const json = await response.json();
        statusList.value = json.data
    } catch (err) {
        // msg.value = "请求发生错误！"
        console.log(err)
        // msgVisible.value = true
    }

}

function Open(index: string) {
    router.push(`/${route.params.key}/${index}`)
}


</script>

<style lang="scss">
.container {
    position: relative;
    padding: 0;
}

.item {
    background-color: #ffffff;
    text-align: center;
    border: 1px solid rgb(225, 225, 225);
    margin: 10px 20px;
    box-sizing: border-box;

    .hangpre {
        height: 10px;
        line-height: 10px;
        width: 100%;
    }

    .hang {
        height: 40px;
        line-height: 40px;
        width: 100%;
        text-align: left;
        box-sizing: border-box;
        padding-left: 30px;

        .lie50 {
            display: inline-block;
            width: 50%;
        }

        .lie30 {
            display: inline-block;
            width: 30%;
        }

        .lie20 {
            display: inline-block;
            width: 20%;
        }

        .lie100 {
            display: inline-block;
            width: 100%;
        }

        .lie70 {
            display: inline-block;
            width: 70%;
        }

        .lie60 {
            display: inline-block;
            width: 60%;
        }

        .lie50 {
            display: inline-block;
            width: 50%;
        }

        .lie40 {
            display: inline-block;
            width: 40%;
        }
    }
}

/* 1. declare transition */
.fade-move,
.fade-enter-active,
.fade-leave-active {
    transition: all 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

/* 2. declare enter from and leave to state */
.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: scaleY(0.01) translate(30px, 0);
}

/* 3. ensure leaving items are taken out of layout flow so that moving
      animations can be calculated correctly. */
.fade-leave-active {
    position: absolute;
}
</style>