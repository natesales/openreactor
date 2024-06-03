export function apiCall(s) {
    let route = `http://${window.location.host}${s}`
    console.log(route)
    fetch(route)
}
