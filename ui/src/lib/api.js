export function apiCall(s) {
    let route = "https://ctl-reactor.westland.as34553.net" + s
    console.log(route)
    fetch(route)
}
