<script>
    import Dot from "$lib/Dot.svelte";

    export let steps = [];
    export let connected = true;
</script>

<main class:connected>
    {#each steps as step, i}
        <div>
            <Dot --color={step.color}/>
            <p class:current={step.current}>{step.name}</p>
        </div>

        {#if i < steps.length - 1}  <!-- If not last dash segment -->
            <span
                    class="spacer"
                    class:connected
                    class:active={steps[i + 1].active}
            ></span>
        {/if}
    {/each}
</main>

<style>
    main {
        padding: 0;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

    .connected {
        width: 90%;
    }

    div {
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        width: 64px;
    }

    .spacer {
        margin-top: 3px;
        width: 100px;
        height: 2px;
        flex-grow: 1;
        background-color: transparent;
    }

    .spacer.active {
        background-color: #bdbdbd;
    }

    .spacer.connected {
        background-color: #444444;
    }

    p {
        margin-top: 10px;
        color: #9f9f9f;
    }

    p.current {
        color: white;
    }
</style>
