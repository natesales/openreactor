<script>
    import {Icon} from "svelte-hero-icons";
    import {apiCall} from "./api.js";

    export let icon = null;
    export let label;
    export let action;
    export let noMargin = false;

    export let danger = false;
    export let wide = false;
    export let small = false;
    export let right = false;
</script>

<button
        class:danger={danger}
        class:wide={wide}
        class:small={small}
        class:right={right}
        class:noMargin={noMargin}
        on:click={() => {
                if (typeof action === "string") {
                    apiCall(action);
                } else {
                    action();
                }
            }}>
        <span>
            {#if icon}
                <span class:icon={label}><Icon src={icon} size="24"/></span>
            {/if}{label}
        </span>
</button>

<style>
    button {
        padding: 1em 3em;
        margin-bottom: 1em;
    }

    button span {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    button.wide {
        width: 75%;
    }

    .noMargin {
        margin: 0;
    }

    button.small {
        padding: 1em 1em;
    }

    button.right {
        border-top-left-radius: 0;
        border-bottom-left-radius: 0;
        width: 0;
        margin-bottom: 0;
    }

    button.danger {
        background-color: #ff0000;
    }

    .icon {
        margin-right: 0.5em;
    }
</style>
