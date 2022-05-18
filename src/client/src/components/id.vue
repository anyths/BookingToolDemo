<template>
    <div v-if="ready" style="height: 100%;position: relative;">
        <t-message v-model="ctlVisible" :content="ctlmsg" :duration="0">
            <template #closeBtn>
                <t-button theme="primary" variant="outline" size="small" shape="round">重试</t-button>
            </template>
        </t-message>
        <t-message v-model="okVisible" :content="okmsg" theme="success" />
        <t-message v-model="errVisible" :content="errmsg" theme="error" />
        <div v-if="dropdownVisible"
            style="height: 40px;width: 100px;position: absolute;top: 48px;right: 0;background-color: #fff;border:1px solid #ddd;border-radius: 5px;z-index: 99999999;line-height: 40px;"
            @click="clearOrder()">清理订单</div>
        <t-navbar @click-right="dropdownVisible = !dropdownVisible"> [ {{ route.params.key }} ] #{{ route.params.id }}
        </t-navbar>
        <div style="width: 100vw;height: calc(100%-44px); background: #fff;text-align: left;padding-top: 1px;">
            <t-divider align="left" style="background-color: #fff;">配置</t-divider>
            <div class="line">
                <div class="w30">项目：{{ current.sport_events_id }}</div>
                <div class="w30">场地：{{ current.venue_id }}</div>
                <div class="w40">价格：{{ current.money }}</div>
            </div>
            <div class="line">
                <div class="w35">频率：{{ current.freq }} ms</div>
                <div class="w65">时段： {{ current.period }}</div>
            </div>
            <t-divider align="left" style="background-color: #fff;">状态</t-divider>

            <div style="height: 10px;"></div>
            <div class="square">
                <div v-for="(value, index) in msgList" :key="index">
                    <div>
                        <div class="msgline">
                            <div :class="{ 'successcmd': value.ok, 'errorcmd': !value.ok }"
                                style="display: inline-block;word-break:normal;">
                                <CheckCircleFilledIcon v-if="value.ok" />
                                <CloseCircleFilledIcon v-else />
                            </div> {{ " " + value.time }} > {{ value.content }}
                        </div>
                    </div>
                </div>
            </div>
            <div style="height: 10px;"></div>
            <div class="line" style="text-align:center;">
                <div class="w5">
                    <TimeFilledIcon v-if="current.status" class="success" size="large" />
                    <loading-icon v-else class="running" size="large" />
                </div>
                <div class="w20">
                    <span style="margin-left: ;20px" :class="current.status ? 'success' : 'running'"> {{ current.status
                            ?
                            '空闲' :
                            '执行中...'
                    }}</span>
                </div>
                <div class="w20">
                    <NextIcon size="large" />
                </div>
                <div class="w20">{{ current.cmd ? current.cmd : '无' }}</div>
            </div>

            <div style="height:15px;"></div>

            <div class="line" style="text-align: center;">
                <t-button theme="primary" size="small" style="margin-right: 10px;" @click="Cmd('auto')"
                    :disabled="!current.status">定时</t-button>
                <t-button theme="primary" size="small" style="margin-right: 10px;" @click="Cmd('now')"
                    :disabled="!current.status">下单</t-button>
                <t-button theme="primary" size="small" @click="Cmd('lock')" style="margin-right: 10px;"
                    :disabled="!current.status">锁单</t-button>
                <t-button theme="danger" size="small" @click="Cmd('break')" :disabled="current.lock"
                    style="margin-right: 10px;">终止</t-button>
                <t-button class="success" size="small" @click="configBtn()">配置</t-button>
            </div>
            <div style="height: 15px;"></div>


            <div style="text-align: center;">
                <t-divider align="left" style="background-color: #fff;">支付</t-divider>
                <div style="height: 10px;"></div>
                <div class="line" style="margin: 0 0px;">

                    <div class="w25" style="font-size: 14px;">定时/下单/锁单 </div>
                    <div class="w5" style="font-size: 14px;">
                        <ChevronRightDoubleIcon size="large" />
                    </div>
                    <div class="w25" style="font-size: 14px;">订单：{{ current.form ? "成功!" : "未生成" }} </div>
                    <div class="w5" style="font-size: 14px;">
                        <ChevronRightDoubleIcon size="large" />
                    </div>
                    <div class="w25" style="font-size: 14px;">
                        <t-button theme="primary" @click="Cmd('pay')" :disabled="!current.form" size="small">生成支付
                        </t-button>
                    </div>

                </div>

                <div style="height: 30px;"></div>
                <div class="line" style="margin: 0;">
                    <div class="w50">
                        <t-button variant="outline" block @click="Check()" :disabled="!wxReady">验 证 状 态</t-button>
                    </div>
                    <div class="w50">
                        <t-button theme="primary" block @click="Go()" :disabled="!wxReady">微 信 支 付</t-button>
                    </div>
                </div>
                <div style="height: 30px;"></div>
                <t-divider :dashed="true">这是我的底线</t-divider>
                <div style="height: 30px;"></div>
            </div>


            <!-- <t-button theme="danger" variant="outline" block @click="Delete()">取 消 订 单</t-button> -->


        </div>
        <t-popup v-model="configVisible" placement="top">
            <ConfVue :sta="current" v-if="configVisible" @nosave="configVisible = false" @save="confSave"></ConfVue>
        </t-popup>
    </div>
    <div v-else class="else">
        <div class="loading" v-if="ready == null">
            <t-loading theme="dots" layout="vertical" size="large" text="连接中..." />
        </div>
        <div class="info" v-else>
            <h3 style="margin: 4vh 2vw;">{{ errRst }}</h3>
            <t-button theme="primary" style="width: 100%;" @click="reload()">重连</t-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, onBeforeUnmount, ref } from 'vue';
import { Status, Msg, wsconf } from '../interface';
import { useRoute } from 'vue-router';
import { LoadingIcon, TimeFilledIcon, NextIcon, CloseCircleFilledIcon, CheckCircleFilledIcon, ChevronRightDoubleIcon } from 'tdesign-icons-vue-next';
import ConfVue from './conf.vue';

const route = useRoute()
const ready = ref<boolean | null>(null)

const dropdownVisible = ref(false)

const errRst = computed(() => {
    if (ready.value == false) {
        return "自动断开连接"
    }
    if (ready.value == null) {
        return "连接中..."
    }
})

const ctlVisible = ref(false)
const ctlmsg = ref("")
const configVisible = ref(false)

const okmsg = ref("")
const okVisible = ref(false)
const errmsg = ref("")
const errVisible = ref(false)
const current = ref<Status>({
    id: "",
    status: false,
    cmd: "",
    form: false,
    msg: "",
    pay: "",
    pct: "",
    oct: "",
    sport_events_id: "",
    money: "",
    venue_id: "",
    period: "",
    freq: 0,
    lock: false,
})
const conf = ref({
    sport_events_id: "",
    venue_id: "",
    period: "",
    money: "",
    freq: 0,
})

function confSave(payload: any) {
    conf.value = payload
    Cmd("conf")
    configVisible.value = false
}
const msgList = ref<Msg[]>([])

const reload = () => {
    window.location.reload()
}

const wxReady = computed(() => {
    if (current.value.pay.length > 0) {
        return true
    }
    return false
})

const configBtn = () => {
    configVisible.value = true
}



function Go() {
    let url = current.value.pay
    let newUrl = url.split("@")[0]
    window.location.href = newUrl
}
function Check() {
    let url = current.value.pay
    let newUrl = url.split("@")[1]
    window.location.href = newUrl
}


var ws = new WebSocket(`${wsconf}${window.location.host}/api/uws/${route.params.key}/${route.params.id}`)


ws.onopen = () => {
    console.log('WebSocket Client Connected');
    let data = {
        "cmd": "sta"
    }
    ws.send(JSON.stringify(data))
    msgList.value.push({
        ok: true,
        time: new Date().toLocaleDateString('en-CA'),
        content: "连接成功..."
    })
}

const heart = ref(true)
var heartCheck = setInterval(() => {
    if (heart.value == false) {
        // 说明已经断开连接了
        ready.value = false
    }
}, 12000)
var timer = setInterval(() => {
    heart.value = false
    ws.send("")
}, 8000)

ws.onerror = (e) => {
    ctlVisible.value = true
    ready.value = false
    errmsg.value = "脚本离线或服务器关闭"
    errVisible.value = true
    console.log(e.type)
    clearInterval(timer)
};

ws.onmessage = (e) => {
    if (typeof e.data === 'string') {
        if (e.data.length == 0) {
            heart.value = true
        } else {
            let data = JSON.parse(e.data)
            if (data.type == "sta") {
                current.value = data.sta
                ready.value = true
            }
            if (data.type == "msg") {
                if (msgList.value.length >= 20) {
                    msgList.value.shift()
                    msgList.value.push(data.msg)
                } else {
                    msgList.value.push(data.msg)
                }
            }
        }

    }
};

function clearOrder() {
    Cmd('clear')
    dropdownVisible.value = false
}

function Cmd(str: string) {
    let data: any = {
        cmd: str
    }
    if (str == "conf") {
        data["data"] = conf.value
    }
    ws.send(JSON.stringify(data))
    if (configVisible.value) {
        configVisible.value = false
    }
}


onBeforeUnmount(() => {
    clearInterval(timer)
    clearInterval(heartCheck)
    ws.close()
})


</script>

<style lang="scss">
.line {
    display: block;
    height: 30px;
    line-height: 30px;
    margin: 0px 30px;
    font-size: 14px;

    .w10 {
        display: inline-block;
        width: 10%;
    }

    .w5 {
        display: inline-block;
        width: 5%;
    }

    .w3 {
        display: inline-block;
        width: 3%;
    }

    .w20 {
        display: inline-block;
        width: 20%;
    }

    .w25 {
        display: inline-block;
        width: 25%;
    }

    .w30 {
        display: inline-block;
        width: 30%;
    }

    .w35 {
        display: inline-block;
        width: 35%;
    }

    .w40 {
        display: inline-block;
        width: 40%;
    }

    .w50 {
        display: inline-block;
        width: 50%;
    }

    .w70 {
        display: inline-block;
        width: 70%;
    }

    .w65 {
        display: inline-block;
        width: 65%;
    }

    .w90 {
        display: inline-block;
        width: 90%;
    }
}

.square {
    box-sizing: border-box;
    background-color: rgb(51, 51, 51);
    color: #fff;
    display: block;
    height: 343px;
    width: calc(100% - 30px);
    padding: 20px 0px;
    margin: 0 15px;
    border: 1px solid #efefef;
    border-radius: 20px;

    .msgline {
        box-sizing: border-box;
        height: 15px;
        line-height: 15px;
        display: block;
        width: 100%;
        padding: 0px 30px;
        font-size: 10px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        -o-text-overflow: ellipsis;
    }
}



.disable {
    color: #BBD3FB;
}

.success {
    color: #008000;
}

.successcmd {
    color: #60be60;
}

.running {
    color: orange;
}

.error {
    color: red;
}

.errorcmd {
    color: rgb(243, 135, 135);
}


.else {
    box-sizing: border-box;
    height: 100vh;
    width: 100%;
    background-color: #ddd;
    position: relative;

    .info {
        position: absolute;
        top: 40%;
        left: 50%;
        width: 60vw;
        background-color: #fff;
        transform: translate(-50%, -50%);
        border-radius: 10px;
    }

    .loading {
        position: absolute;
        top: 50%;
        left: 50%;
        height: 30vh;
        width: 70vw;
        transform: translate(-50%, -50%);
    }
}
</style>