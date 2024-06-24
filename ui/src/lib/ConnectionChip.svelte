<script>
    import Dot from "./Dot.svelte";

    export let connected = false;
</script>

<main>
    <Dot --color="var(--{connected ? 'green' : 'red'})" />
    <span>{connected ? "Connected" : "Disconnected"}</span>
</main>

<style>
    main {
        display: flex;
        align-items: center;
    }
</style>
