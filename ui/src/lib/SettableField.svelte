<script>
    import {apiCall} from "./api.js";
    import ActionButton from "./ActionButton.svelte";

    export let label;
    export let prefix;
    export let v = 0.0;
</script>

<div>
    <input type="number" inputmode="decimal" bind:value={v}>
    <ActionButton right label="Set" action={() => apiCall(`${prefix}${v}`)}/>
</div>
<label>{label}</label>

<style>
    div {
        display: flex;
        flex-direction: row;
        align-items: center;
        height: 52px;
    }

    input {
        margin: 0;
        padding: 0 16px;
        width: 75px;
        height: 100%;
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
    }
</style>
