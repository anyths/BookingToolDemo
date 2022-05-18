interface Status {
    id: string,
    status: boolean,
    cmd: string,
    form: boolean,
    msg: string,
    pay: string,
    pct: string,
    oct: string,
    sport_events_id: string,
    money: string,
    venue_id: string,
    period: string,
    freq: number,
    lock: boolean,
}
interface Msg {
    ok:boolean,
    time: string,
    content: string,
}
const http = "http://"
const wsconf = "ws://" 
const api = {
    base : http + window.location.host,
    sta : http + window.location.host + "/api/sta/",
}

export {Status, api, Msg, wsconf,http}