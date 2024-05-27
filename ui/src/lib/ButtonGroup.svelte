<script>
    import {apiCall} from "./api.js";

    export let prefix;
    export let options = [];
</script>

<main>
    {#each options as option}
        <input type="button" id={option} name={option} on:click={() => apiCall(prefix+option)}/>
        <label for={option}>{option}</label>
    {/each}
</main>

<style>
    main {
        width: fit-content;
        overflow: hidden;
        display: flex;
        flex-direction: row;
        flex-wrap: nowrap;
        margin-right: 6px;
        margin-bottom: 6px;
        border-radius: 6px;
        outline: none;
        background: var(--background);
    }
    input {
        display: none !important;
    }

    label {
        padding: 15px 24px !important;
        margin: 0 !important;
        cursor: pointer;
        transition: all 0.3s;
    }

    label:hover {
        background: var(--button-hover);
    }
</style>
