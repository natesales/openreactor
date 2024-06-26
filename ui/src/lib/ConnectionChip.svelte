<script>
    import Dot from "./Dot.svelte";

    export let connected = false;
</script>

<main>
    <div>
        <Dot --color="var(--{connected ? 'green' : 'red'})" />
        <span>{connected ? "Connected" : "Disconnected"}</span>
    </div>
</main>

<style>
    div {
        display: flex;
        align-items: center;
    }
</style>
