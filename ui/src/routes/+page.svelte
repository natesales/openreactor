<script>
    import ActionButton from "$lib/ActionButton.svelte";
    import {
        BoltSlash,
        ExclamationTriangle,
        Power,
        SpeakerWave,
        SpeakerXMark,
        StopCircle,
        XMark
    } from "svelte-hero-icons";
    import SettableField from "$lib/SettableField.svelte";
    import {apiCall} from "$lib/api";
    import ButtonGroup from "../lib/ButtonGroup.svelte";
    import IconToggle from "../lib/IconToggle.svelte";

    function eStop() {
        apiCall("/hv/set?v=0");
        apiCall("/mfc/set?slpm=0");
        apiCall("/turbo/turbo/off")
    }

    let muted = false;
</script>

<h1>OpenReactor Mobile Control</h1>

<IconToggle onIcon={SpeakerXMark} offIcon={SpeakerWave} bind:value={muted}/>
Muted? {muted ? "Yes" : "No"}

<div class="row">
    <div class="group">
        <h2>HV</h2>
        <div>
            <ActionButton icon={BoltSlash} label="Off" action="/hv/set?v=0"/>
            <SettableField label="Voltage Setpoint" prefix="/hv/set?v="/>
        </div>
        <h3>Quick Actions (x10 kV)</h3>
        <ButtonGroup prefix="/hv/set?v=" options={[0.5, 1, 1.5, 2, 2.25]}/>
    </div>

    <div class="group">
        <h2>MFC</h2>
        <div>
            <ActionButton icon={XMark} label="Close" action="/mfc/set?v=0"/>
            <SettableField label="Flow Rate Setpoint" prefix="/mfc/set?slpm="/>

        </div>
        <h3>Quick Actions (SLPM)</h3>
        <ButtonGroup prefix="/mfc/set?slpm=" options={[0.02, 0.05, 0.1, 0.2]}/>
    </div>
</div>

<div class="row">
    <div class="group">
        <h2>Turbo Pump</h2>
        <ActionButton icon={Power} label="On" action="/turbo/turbo/on"/>
        <ActionButton icon={StopCircle} label="Off" action="/turbo/turbo/off"/>
    </div>
</div>

<ActionButton danger wide icon={ExclamationTriangle} label="Emergency Stop" action={eStop}/>

<style>
    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        padding-bottom: 20px;
        border-bottom: 1px solid var(--border);
    }

    .group {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }
</style>
