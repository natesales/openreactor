<script>
    import {logs} from "$lib/log";

    let textarea;
    $: {
        if (textarea) {
            textarea.innerHTML = $logs.map(log => `${log.time}: ${log.text}`).join("\n");
            textarea.scrollTop = textarea.scrollHeight + 10;
        }
    }
</script>

<main>
    <textarea bind:this={textarea} readonly></textarea>
</main>

<style>
    textarea {
        margin-top: 30px;
        font-family: monospace;
        font-size: 1.25em;
    }
</style>
