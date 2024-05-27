import {writable} from "svelte/store";

export const logs = writable([]);

export function addLog(text) {
    logs.update((currentLogs) => {
        return [...currentLogs, {
            time: new Date().toLocaleString(),
            text: text
        }];
    });
}
