<script>
    import {
        ArrowPath,
        ArrowUturnLeft, ArrowUturnRight,
        BoltSlash, ChevronRight,
        ExclamationTriangle,
        Power, Signal, SignalSlash,
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
    import StepperGroup from "$lib/StepperGroup.svelte";
    import {wsURL} from "$lib/api";

    let muted = false;
    let refreshStepper;

    function eStop() {
        apiCall("/hv/set?v=0");
        apiCall("/mfc/set?slpm=0");
        apiCall("/turbo/off");
    }

    let wsConnected;
    const ws = new WebSocket(wsURL);

    function emit(name) {
        ws.send(JSON.stringify({
            name: name
        }))
    }

    function wsConnect() {
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
            switch (data["name"]) {
                case "logMessage":
                    addLog(data["message"]);
                    break;
                case "audioAlert":
                    addLog(data["text"]);
                    if ('speechSynthesis' in window) {
                        if (!muted) {
                            window.speechSynthesis.speak(new SpeechSynthesisUtterance(data["text"]));
                        }
                    } else {
                        addLog("Speech synthesis not supported");
                    }
                    break;
                case "fsmStateChange":
                    if (data["state"]) {
                        addLog("FSM state changed to " + data["state"]);
                    }
                    refreshStepper();
                    break;
                default:
                    addLog("Unknown message type: " + data["name"]);
            }
        };
    }

    onMount(wsConnect);
</script>

<svelte:head>
    <title>OpenReactor | Mobile Control</title>
</svelte:head>

<div class="row">
    <h1>OpenReactor Mobile Control</h1>
    <ConnectionChip connected={wsConnected}/>
</div>

<div class="section">
    <div class="row">
        <IconToggle onIcon={SpeakerXMark} offIcon={SpeakerWave} bind:value={muted}/>
        <ActionButton small noMargin icon={ArrowPath} action={() => emit("fsmReset")} label=""/>
        <ActionButton small noMargin icon={ChevronRight} action={() => emit("fsmNext")} label=""/>
        <ActionButton danger wide noMargin icon={ExclamationTriangle} label="Emergency Stop" action={eStop}/>
    </div>

    <StepperGroup bind:refresh={refreshStepper}/>
</div>

<div class="row section">
    <div class="group">
        <h2>HV</h2>
        <div>
            <ActionButton icon={BoltSlash} label="Off" action="/hv/set?v=0"/>
            <SettableField label="Voltage Setpoint" prefix="/hv/set?v="/>
        </div>
        <h3>Quick Actions (x10 kV)</h3>
        <ButtonGroup prefix="/hv/set?v=" options={[0.5, 1, 1.75, 2.25, 3.75]}/>
    </div>

    <div class="group">
        <h2>MFC</h2>
        <div>
            <ActionButton icon={XMark} label="Close" action="/mfc/set?slpm=0"/>
            <SettableField label="Flow Rate Setpoint" prefix="/mfc/set?slpm="/>
        </div>
        <h3>Quick Actions (SLPM)</h3>
        <ButtonGroup prefix="/mfc/set?slpm=" options={[0.02, 0.05, 0.1, 0.2]}/>
    </div>
</div>

<div class="row section">
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
    }

    .section {
        border-bottom: 1px solid var(--border);
    }

    .group {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }

    h1 {
        margin: 0;
        padding: 0;
    }
</style>
