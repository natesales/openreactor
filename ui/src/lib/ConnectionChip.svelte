<script>
    export let connected = false;
</script>

<main>
    {#if connected}
        <span class="dot connected"></span>
        <span>Connected</span>
    {:else}
        <span class="dot disconnected"></span>
        <span>Disconnected</span>
    {/if}
</main>

<style>
    main {
        display: flex;
        align-items: center;
    }

    .dot {
        height: 10px;
        width: 10px;
        border-radius: 50%;
        margin-right: 5px;
    }

    .connected {
        background-color: #0dbd0d;
    }

    .disconnected {
        background-color: red;
    }
</style>
