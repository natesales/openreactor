<script>
    import { onMount } from "svelte";

    function apiCall(s) {
        let route = "https://ctl-reactor.westland.as34553.net"+s
        console.log(route)
        fetch(route)
    }

    let v = 0.0;
</script>

<h1>OpenReactor HV Supply Controller</h1>
<label for="vset">Voltage Setpoint</label>
<input name="vset" type="number" bind:value={v}>

<button on:click={() => {
    apiCall(`/hv/set?v=${v}`)
}}>Set</button>

<button on:click={() => {
    apiCall(`/hv/set?v=0`)
}}>Disable</button>



<button on:click={() => {
    apiCall(`/turbo/turbo/off`)
}}>Off</button>
