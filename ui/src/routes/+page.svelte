<script>
    import ActionButton from "$lib/ActionButton.svelte";
    import {BoltSlash, ExclamationTriangle, Power, StopCircle, XMark} from "svelte-hero-icons";
    import SettableField from "$lib/SettableField.svelte";
    import {apiCall} from "$lib/api";
    import ButtonGroup from "../lib/ButtonGroup.svelte";

    function eStop() {
        apiCall("/hv/set?v=0");
        apiCall("/mfc/set?slpm=0");
    }
</script>

<h1>OpenReactor HV Supply Controller</h1>

<div class="row">
    <div class="group">
        <h2>HV</h2>
        <div>
            <ActionButton icon={BoltSlash} label="Off" action="/hv/set?v=0"/>
            <SettableField label="Voltage Setpoint" prefix="/hv/set?v="/>
        </div>
        <h3>Quick Actions</h3>
        <ButtonGroup prefix="/hv/set?v=" options={[0, 1, 5, 10]}/>
    </div>

    <div class="group">
        <h2>MFC</h2>
        <div>
            <ActionButton icon={XMark} label="Close" action="/mfc/set?v=0"/>
            <SettableField label="Flow Rate Setpoint" prefix="/mfc/set?slpm="/>

        </div>
        <h3>Quick Actions</h3>
        <ButtonGroup prefix="/mfc/set?slpm=" options={[0, 0.1, 0.9]}/>
    </div>
</div>

<div class="row">
    <div class="group">
        <h2>Turbo Pump</h2>
        <ActionButton icon={Power} label="On" action="/turbo/turbo/on"/>
        <ActionButton icon={StopCircle} label="Off" action="/turbo/turbo/off"/>
    </div>
</div>

<div class="bottom">
    <ActionButton danger wide icon={ExclamationTriangle} label="Emergency Stop" action={eStop}/>
</div>

<style global>
    :root {
        --background: #151515;
        --background-body: #0f0f0f;
        --background-alt: #272727;
        --selection: #bb07d6;
        --text-main: #e3e3e3;
        --text-bright: #ffffff;
        --text-muted: #a9b1ba;
        --links: #c706e3;
        --focus: rgba(230, 230, 230, 0.67);
        --border: #272727;
        --code: #ffffff;
        --animation-duration: 0.1s;
        --button-base: #373737;
        --button-hover: #121212;
        --scrollbar-thumb: var(--button-hover);
        --scrollbar-thumb-hover: #ffffff;
        --form-placeholder: #a9a9a9;
        --form-text: #fff;
        --variable: #d941e2;
        --highlight: #ffffff;
    }

    .bottom {
        position: fixed;
        bottom: 0;
        left: 0;
        width: calc(100% - 20px);
        margin: 0 10px;
    }

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
