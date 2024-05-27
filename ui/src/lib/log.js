import {writable} from "svelte/store";

export const logs = writable([]);

export function addLog(time, text) {
    logs.update((currentLogs) => {
        return [...currentLogs, {
            time: time,
            text: text
        }];
    });
}
