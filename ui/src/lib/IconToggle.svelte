<script>
    import ActionButton from "./ActionButton.svelte";

    export let value = false;
    export let onIcon;
    export let offIcon;
</script>

<ActionButton small icon={value ? onIcon : offIcon} action={() => {value = !value}} label=""/>
