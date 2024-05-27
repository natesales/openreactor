<script>
    import {
        BoltSlash,
        ExclamationTriangle,
        Power,
        SpeakerWave,
        SpeakerXMark,
        StopCircle,
        XMark
    } from "svelte-hero-icons";

    import {apiCall} from "$lib/api";
    import {addLog} from "$lib/log";
    import {onMount} from "svelte";

    import ActionButton from "$lib/ActionButton.svelte";
    import SettableField from "$lib/SettableField.svelte";
    import ButtonGroup from "$lib/ButtonGroup.svelte";
    import IconToggle from "$lib/IconToggle.svelte";
    import Logs from "$lib/Logs.svelte";
    import ConnectionChip from "$lib/ConnectionChip.svelte";

    let muted = false;

    function eStop() {
        apiCall("/hv/set?v=0");
        apiCall("/mfc/set?slpm=0");
        apiCall("/turbo/off");
    }

    let wsConnected;

    function wsConnect() {
        const ws = new WebSocket(`ws://${window.location.host}/alert/ws`);
        addLog("Connecting to WebSocket server...");

        ws.onopen = () => {
            wsConnected = true;
            addLog("WebSocket connected");
        };
        ws.onclose = () => {
            wsConnected = false;
            addLog("WebSocket closed");
        };
        ws.onerror = (event) => {
            addLog("WebSocket error");
        };
        ws.onmessage = (event) => {
            let data = JSON.parse(event.data);
            switch (data["type"]) {
                case "logMessage":
                    addLog(data["message"]);
                    break;
                case "audioAlert":
                    addLog(data["text"]);
                    if (!muted) {
                        new Audio("/tts/audio?text=" + encodeURIComponent(data["text"])).play();
                    }
                    break;
                default:
                    addLog("Unknown message type: " + data["type"]);
            }
        };
    }

    onMount(wsConnect);
</script>

<svelte:head>
    <title>OpenReactor | Mobile Control</title>
</svelte:head>

<h1>OpenReactor Mobile Control</h1>

<ActionButton danger wide icon={ExclamationTriangle} label="Emergency Stop" action={eStop}/>

<div class="row">
    <IconToggle onIcon={SpeakerXMark} offIcon={SpeakerWave} bind:value={muted}/>
    <ConnectionChip connected={wsConnected}/>
</div>

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
        <ActionButton icon={Power} label="On" action="/turbo/on"/>
        <ActionButton icon={StopCircle} label="Off" action="/turbo/off"/>
    </div>
</div>

<Logs/>

<style>
    .row {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        align-items: center;
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
