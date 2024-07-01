<script>
    import {onMount} from "svelte";
    import Stepper from "$lib/Stepper.svelte";
    import {apiURL} from "$lib/api";

    let steps = [];
    let errors = [];

    export function refresh() {
        fetch(apiURL + "/fsm/states")
            .then(response => response.json())
            .then(data => {
                errors = data["errorStates"].map(e => {
                    let active = data["errors"].includes(e);
                    return {
                        name: e,
                        active: false,
                        current: active,
                        color: active ? "red" : "",
                    }
                });

                let found = true;
                steps = data["states"].map(step => {
                    let e = {
                        name: step,
                        active: found,
                        current: step === data["active"],
                        color: found ? "white" : "",
                    };
                    if (step === data["active"]) {
                        found = false;
                        e.color = "var(--green)";
                    }
                    return e;
                });
            });
    }

    onMount(refresh);
</script>

<main>
    <h3>Startup</h3>
    <Stepper steps={steps.slice(0, Math.ceil(steps.length / 2))}/>

    <h3>Runtime</h3>
    <Stepper steps={steps.slice(Math.ceil(steps.length / 2))}/>

    <h3>Errors</h3>
    <Stepper steps={errors} connected={false}/>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    h3 {
        align-self: flex-start;
        margin: 16px;
        font-weight: normal;
    }
</style>
