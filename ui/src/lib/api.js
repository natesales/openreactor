export function apiCall(s) {
    let route = `http://${window.location.host}${s}`
    console.log(route)
    fetch(route)
}

export let apiURL = "http://localhost:8084";

// export let wsURL = `ws://${window.location.host}/maestro/ws`;
export let wsURL = "ws://localhost:8084/ws";
