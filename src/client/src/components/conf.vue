<template>
    <div style="width: 100%; height: auto; background: #fff;padding-top: 15px;">
        <h3 style="height: 40px;line-height:30px;text-align: left;padding-left: 15px;">修改配置</h3>
        <t-cell-group style="text-align: left;" title="订单选项">
            <t-cell arrow title="运动项目"
                :note="sportOptions.find(v => v.value == conf.sport_events_id)?.label || '选择运动项目'"
                @click="openSelect('sport_events_id')" />
            <t-cell arrow title="场地"
                :note="venueOptions[conf.sport_events_id].find(v => v.value == conf.venue_id)?.label || '选择运动场地'"
                @click="openSelect('venue_id')" />
            <t-cell arrow title="时间段" :note="periodShow || '选择时间段'" @click="ClickPeriod()" />
        </t-cell-group>
        <t-input label="价格" align="right" v-model="conf.money" placeholder="0" suffix="软妹" />
        <div>
            <t-button variant="text" block style="background-color: #FDF6EC;" @click="GetPrice()">{{ moneyLoading ?
                    "计算中..." : "计 算 价 格"
            }}</t-button>
            <t-button theme="primary" block @click="checkOk()" style="background-color:#85C88A ;">空 场 查 询</t-button>
        </div>
        <t-cell-group title="脚本参数" style="text-align: left;">
            <t-cell title="频率[ms]">
                <t-stepper v-model.number="conf.freq" :step="50" :min="50" :max="10000" />
            </t-cell>
        </t-cell-group>
        <div style="height: 50px;line-height: 50px;font-size: 12px;"></div>
        <div>
            <div style="display: inline-block;width: 50%;border-top: 1px solid #efefef;">
                <t-button variant="text" block @click="$emit('nosave')">取 消</t-button>
            </div>
            <div style="display: inline-block;width: 50%;border-top: 1px solid #efefef;">
                <t-button theme="primary" block :disabled="!props.sta.status" @click="up()">保 存</t-button>
            </div>
        </div>
        <div>
        </div>
    </div>
    <t-popup v-if="showSelect" v-model="showSelect" placement="bottom">
        <t-picker v-model="currentValue" title="请选择选项" @confirm="onConfirm" @cancel="onCancel"
            style="height: 30vh;padding-right:15px">
            <t-picker-item :options="currentOptions" />
        </t-picker>
    </t-popup>
    <t-popup v-if="showSelectPeriod" v-model="showSelectPeriod" placement="bottom">
        <div class="nav-sel">
            <div class="nav-btn" @click="CancelPeriod">取消</div>
            <div class="nav-btn">选择时间段</div>
            <div class="nav-btn" style="color: #0052d9;" @click="SavePeriod">确认</div>
        </div>
        <div style="min-height: 50vh;width:100%;background-color: #fff;overflow: scroll;">

            <t-checkbox-group v-model:value="periodCheckBox" :max="2" @change="changeFn">
                <t-checkbox style="height: 25px;line-height:25px;margin: 10px 100px;"
                    v-for="(index, key) in periodOptions[conf.sport_events_id]" name="periodCheckBox"
                    :value="index.value" :label="index.label"></t-checkbox>
            </t-checkbox-group>
            <div style="height: 30px;"></div>
        </div>

    </t-popup>
    <t-popup v-if="periodOkVisible" v-model="periodOkVisible" placement="bottom">

        <div class="periodok" v-if="okData&&periodOkVisible">
            <div style="height: 38px;line-height:38px;text-align: center;background-color: #efefef;"
                @click="periodOkVisible = false">
                <ChevronDownIcon size="large" />
            </div>
            <div style="height: 5px;"></div>
            <div style="height: 20px;text-align: center;font-size: 12px;color: #aaa;">今天是: {{ new
                    Date().toLocaleDateString('en-CA')
            }}</div>
            <div style="height: 46px;text-align: center;line-height: 46px;">
                <t-button theme="primary" shape="round" size="small" @click="RmDt" :disabled="n == 0 || addLoading">-
                </t-button>
                <div v-if="addLoading" style="display: inline-block;margin: 0 20px;width: 100px;">
                    <t-loading theme="spinner" text="解析中..." />
                </div>
                <div v-else style="display: inline-block;margin: 0 20px;width: 100px;">{{ dt }}</div>
                <t-button theme="primary" shape="round" :disabled="addLoading" size="small" @click="AddDt">+</t-button>
                <!-- <ChevronRightIcon size="large" /> -->

            </div>
            <div style="height: 5px;"></div>
            <div style="position: relative;">
                <div class="sqline" style="border-top: 1px solid #efefef;line-height:30px;color: #aaa;">
                    <div class="sqf">时\场</div>
                    <div class="sq" v-for="(i, k) in venueOptions[conf.sport_events_id]">{{ k + 1 }}</div>
                </div>
                <div class="sqline" v-for="(index, key) in periodOptions[conf.sport_events_id]">
                    <div class="sqf">{{ index.label.replace("-", " ") }}</div>
                    <div class="sq" v-for="(i, k) in venueOptions[conf.sport_events_id]">
                        <CheckCircleFilledIcon v-if="okData[i.value + '-' + index.value] == 0"
                            style="margin: 0;padding:0;position: absolute;left: 5px;top: 8px;color: green;" />
                        <CloseCircleFilledIcon v-else
                            style="margin: 0;padding:0;position: absolute;left: 5px;top: 8px;color: rgb(201, 58, 58);" />
                    </div>
                </div>
            </div>
            <div style="height: 10px;"></div>
            <t-divider align="center">没了没了</t-divider>
            <div style="height: 20px;"></div>
        </div>
        <div v-else style="height: 50vh;background-color: #fff;position: relative;">
            <t-loading theme="spinner" text="解析中" layout="vertical"
                style="position: absolute;left: 50%;top: 50%;transform: translate(-50%,-50%);" />
        </div>
    </t-popup>
</template>
<script lang="ts" setup>
import { computed } from '@vue/reactivity';
import { ref } from 'vue';
import { CheckCircleFilledIcon, CloseCircleFilledIcon, ChevronDownIcon } from "tdesign-icons-vue-next"
import { http } from "../interface"
import { add } from 'lodash';

const props = defineProps(["sta"])
const emits = defineEmits(['save', "nosave"])
const dt = ref(new Date().toLocaleDateString('en-CA'))
const okData = ref<any>()
const n = ref(0)
const addLoading = ref(false)
const checklistloading = ref(false)
interface Conf {
    sport_events_id: string,
    venue_id: string,
    period: string,
    money: string,
    freq: number,
}
const conf = ref<Conf>({
    sport_events_id: props.sta.sport_events_id,
    venue_id: props.sta.venue_id,
    period: props.sta.period,
    money: props.sta.money,
    freq: props.sta.freq,
})
const periodOkVisible = ref(false)
const moneyLoading = ref(false)

const periodCheckBox = ref<any>([])

const showSelect = ref(false)
const showSelectPeriod = ref(false)
const sportOptions = [
    { label: "乒乓球", value: "33" },
    { label: "羽毛球", value: "34" },
    { label: "网球", value: "35" },
]
interface VenueOptions {
    [x: string]: LabelValue[]
}
interface LabelValue {
    label: string,
    value: string
}
const venueOptions: VenueOptions = {
    "33": [
        { label: '乒乓球1', value: '145' },
        { label: '乒乓球2', value: '146' },
        { label: '乒乓球3', value: '147' },
        { label: '乒乓球4', value: '148' },
        { label: '乒乓球5', value: '149' },
        { label: '乒乓球6', value: '150' },
        { label: '乒乓球7', value: '151' },
        { label: '乒乓球8', value: '152' },
        { label: '乒乓球9', value: '153' },
        { label: '乒乓球10', value: '154' },
        { label: '乒乓球11', value: '171' },
        { label: '乒乓球12', value: '172' },
        { label: '乒乓球13', value: '173' },
    ],
    "34": [
        { label: '羽毛球1', value: '155' },
        { label: '羽毛球2', value: '156' },
        { label: '羽毛球3', value: '157' },
        { label: '羽毛球4', value: '158' },
        { label: '羽毛球5', value: '159' },
        { label: '羽毛球6', value: '160' },
        { label: '羽毛球7', value: '161' },
        { label: '羽毛球8', value: '162' },
        { label: '羽毛球9', value: '163' },
        { label: '羽毛球10', value: '164' },
        { label: '羽毛球11', value: '165' },
        { label: '羽毛球12', value: '166' },
    ],
    "35": [
        { label: '网球1', value: '167' },
        { label: '网球2', value: '168' },
        { label: '网球3', value: '169' },
        { label: '网球4', value: '170' },
    ]
}
interface PeriodOptions {
    [x: string]: PeriodLabelValue[]
}
interface PeriodLabelValue extends LabelValue {
    daytype: string
}
const periodOptions: PeriodOptions = {
    "33": [
        { daytype: 'morning', label: '07:00-08:00', value: '328094' },
        { daytype: 'morning', label: '08:00-09:00', value: '328095' },
        { daytype: 'morning', label: '09:00-10:00', value: '328096' },
        { daytype: 'day', label: '10:00-11:00', value: '328097' },
        { daytype: 'day', label: '11:00-12:00', value: '328098' },
        { daytype: 'night', label: '12:00-13:00', value: '328099' },
        { daytype: 'night', label: '13:00-14:00', value: '328100' },
    ],
    "34": [
        { daytype: 'morning', label: '07:00-08:00', value: '328101' },
        { daytype: 'morning', label: '08:00-09:00', value: '328102' },
        { daytype: 'morning', label: '09:00-10:00', value: '328103' },
        { daytype: 'morning', label: '10:00-11:00', value: '328104' },
        { daytype: 'morning', label: '11:00-12:00', value: '328105' },
        { daytype: 'morning', label: '12:00-13:00', value: '328106' },
        { daytype: 'morning', label: '13:00-14:00', value: '328107' },
        { daytype: 'day', label: '14:00-15:00', value: '328125' },
        { daytype: 'day', label: '15:00-16:00', value: '328126' },
        { daytype: 'day', label: '16:00-17:00', value: '328127' },
        { daytype: 'day', label: '17:00-18:00', value: '328128' },
        { daytype: 'night', label: '18:00-19:00', value: '328129' },
        { daytype: 'night', label: '19:00-20:00', value: '328130' },
        { daytype: 'night', label: '20:00-21:00', value: '328131' },
        { daytype: 'night', label: '21:00-22:00', value: '328132' },
    ],
    "35": [
        { daytype: 'morning', label: '07:00-08:00', value: '328108' },
        { daytype: 'morning', label: '08:00-09:00', value: '328109' },
        { daytype: 'morning', label: '09:00-10:00', value: '328110' },
        { daytype: 'day', label: '10:00-11:00', value: '328111' },
        { daytype: 'day', label: '11:00-12:00', value: '328112' },
        { daytype: 'night', label: '12:00-13:00', value: '328113' },
        { daytype: 'night', label: '13:00-14:00', value: '328114' },
    ]
}

const periodShow = computed(() => {
    let arr = conf.value.period.split(",")
    let str1 = periodOptions[conf.value.sport_events_id].find(i => i.value == arr[0])?.label
    let str2 = periodOptions[conf.value.sport_events_id].find(i => i.value == arr[1])?.label
    if (str1 == undefined || str2 == undefined) {
        return "请选择时段"
    }
    return str1 + "," + str2
})
const currentSelect = ref<string>("")
const currentOptions = ref<any>([])
const currentValue = ref<any>([])
function openSelect(index: string) {
    if (index == "sport_events_id") {
        currentSelect.value = index
        currentValue.value = [conf.value[index]]
        currentOptions.value = sportOptions
    }
    if (index == "venue_id") {
        currentSelect.value = index
        currentValue.value = [conf.value[index]]
        currentOptions.value = venueOptions[conf.value.sport_events_id]
    }
    showSelect.value = true

}
function changeFn() {
    // 实时选择多选时间段的触发函数
}
function GetPrice() {
    moneyLoading.value = true
    // 价格表
    let s: any = {}
    let m: number = 0
    fetch(`${http}${window.location.host}/api/price/${conf.value.sport_events_id}`, {
        method: "GET",
        headers: {
            "token": "87d8a2bf-27b5-46b1-bad6-1796e488e1ac"
        }
    })
        .then((res) => {
            return res.json()
        })
        .then((res) => {
            let data = res.data
            for (let v in data) {
                s[data[v].daytype] = data[v].price
            }
            let rst = []
            let arr = conf.value.period.split(",")
            let str1 = periodOptions[conf.value.sport_events_id].find(i => i.value == arr[0])?.value
            if (str1 != undefined) {
                rst.push(str1)
            }
            let str2 = periodOptions[conf.value.sport_events_id].find(i => i.value == arr[1])?.value
            if (str2 != undefined) {
                rst.push(str2)
            }
            for (let i of rst) {
                let daytype = periodOptions[conf.value.sport_events_id].find(item => item.value == i)?.daytype
                if (daytype != undefined) {
                    m += s[daytype]
                }
            }
            conf.value.money = m.toString()
            moneyLoading.value = false
            return
        })
}

function checkOk() {
    dt.value = new Date().toLocaleDateString('en-CA')
    n.value = 0
    periodOkVisible.value = true
    searchOk()
}
function searchOk() {
    addLoading.value = true
    fetch(`${http}${window.location.host}/api/search/${conf.value.sport_events_id}?day=${dt.value}`, {
        method: "GET",
        headers: {
            "token": "87d8a2bf-27b5-46b1-bad6-1796e488e1ac"
        }
    }).then(res => {
        return res.json()
    }).then(res => {
        okData.value = res.data
        checklistloading.value = false
        addLoading.value = false

    }).catch(err => {
        console.log(err);
    })
}
function AddDt() {
    addLoading.value = true
    checklistloading.value = true
    n.value++
    dt.value = new Date(new Date().setDate(new Date().getDate() + n.value)).toLocaleDateString('en-CA')
    searchOk()
}
function RmDt() {
    if (n.value == 0) {
        return
    }
    addLoading.value = true
    n.value--
    dt.value = new Date(new Date().setDate(new Date().getDate() + n.value)).toLocaleDateString('en-CA')
    searchOk()
}
// fetch("http://gym.dazuiwl.cn/api/sport_events/hour/id/35")
//     .then((res) => {
//         return res.json()
//     })
//     .then((res) => {
//         let data = res.data
//         let rst = data.map((res: { begintime_text: string; endtime_text: string; id: number; daytype: string }) => {
//             return {
//                 daytype: res.daytype,
//                 label: res.begintime_text + "-" + res.endtime_text,
//                 value: res.id.toString()
//             }
//         })
//         console.log(rst)
//     })
// fetch("http://gym.dazuiwl.cn/api/sport_events/field/id/34")
//     .then((res) => {
//         return res.json()
//     })
//     .then((res) => {
//         let data: d = res.data
//         let nd = []
//         for (let v in data) {
//             let n = { label: "", value: "" }
//             n.label = data[v].name
//             n.value = data[v].id.toString()
//             nd.push(n)
//         }
//         console.log(nd)
//     })


function onConfirm(val: string[]) {
    if (currentSelect.value == "sport_events_id") {
        conf.value.sport_events_id = val[0]
    }
    if (currentSelect.value == "venue_id") {
        conf.value.venue_id = val[0]
    }
    showSelect.value = false
}
function onCancel() {
    showSelect.value = false
}
function CancelPeriod() {
    showSelectPeriod.value = false
}
function SavePeriod() {
    console.log(periodCheckBox.value);
    conf.value.period = periodCheckBox.value.join(',')
    showSelectPeriod.value = false
}
function ClickPeriod() {
    let rst = []
    let rs = conf.value.period.includes(",", 0)
    if (rs) {
        let arr = conf.value.period.split(",")
        let str1 = periodOptions[conf.value.sport_events_id].find(i => i.value == arr[0])?.value
        if (str1 != undefined) {
            rst.push(str1)
        }
        let str2 = periodOptions[conf.value.sport_events_id].find(i => i.value == arr[1])?.value
        if (str2 != undefined) {
            rst.push(str2)
        }
    }
    periodCheckBox.value = rst
    showSelectPeriod.value = true
}

const up = () => {
    console.log(conf.value);
    emits('save', conf.value)
}
</script>

<style lang="scss">
.nav-sel {
    height: 48px;
    background-color: #fff;
    border-bottom: 1px solid #efefef;
    display: flex;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    -webkit-box-pack: justify;
    -ms-flex-pack: justify;
    justify-content: space-between;
    padding-right: 15px;

    .nav-btn {
        font-size: 16px;
        height: 100%;
        padding: 12px 16px;
    }
}

.periodok {
    min-height: 50vh;
    max-height: 100%;
    background-color: #fff;
    width: 100%;
    text-align: left;
    vertical-align: unset;
    box-sizing: border-box;

    .sqline {

        border-bottom: 1px solid #efefef;

        .sqf {
            padding-left: 8px;
            position: relative;
            display: inline-block;
            height: 26px;
            width: 44px;
            font-size: 10px;
            font-weight: bold;
            border-right: 1px solid #efefef;
        }

        .sq {
            text-align: center;
            position: relative;
            display: inline-block;
            height: 26px;
            width: 26px;
            border-right: 1px solid #efefef;
        }

        .red {
            color: rgb(201, 58, 58);
        }
    }

}
</style>