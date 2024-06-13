<script lang="ts">
    import type { route } from "wailsjs/go/models";
    import { i18n } from "../i18n/i18n";
    import { validIP, validNetMask } from "../utils/utils";

    export let show: boolean = false;
    export let rt: route.Route;
    export let onCancel = () => {};
    export let onOkay = (old: route.Route, want: route.Route) => {};
    export let interfaceNames: string[];

    let ip: string = "";
    let destination: string = "";
    let mask: string = "";
    let interfaceName: string = "";

    let onChange = () => {};

    function _onCancel() {
        onCancel();
        close();
    }

    function selectInterfaceName(event: any) {
        interfaceName = event.target.value;
    }

    let onDestinationChange = (event: any) => {
        let ipAndMask = cidrToIpAndMask(event.target.value);
        if (ipAndMask.mask != "") {
            mask = ipAndMask.mask;
        }
        mask = ipAndMask.mask;
        ip = ipAndMask.ip;
    };

    let onMaskChange = (event: any) => {
        mask = event.target.value;
    };

    function cidrToIpAndMask(cidr: string) {
        const parts = cidr.split("/");
        const ip = parts[0];
        let maskLength = parseInt(parts[1], 10);
        let maskString = "";
        if (maskLength) {
            const mask = [];
            for (let i = 0; i < 4; i++) {
                let octet = 0;
                if (maskLength >= 8) {
                    octet = 255;
                    maskLength -= 8;
                } else {
                    octet = (255 << (8 - maskLength)) & 255;
                    maskLength = 0;
                }
                mask.push(octet);
                maskString = mask.join(".");
            }
        }

        return {
            ip: ip,
            mask: maskString,
        };
    }

    function _onOkay() {
        if (interfaceName == "") {
            alert("select interface");
            return;
        }

        if (!validIP(ip)) {
            alert("invalid ip address");
            return;
        }
        if (!validNetMask(mask)) {
            alert("invalid net mask");
            return;
        }
        let want: route.Route = {
            destination: ip,
            mask: mask,
            gateway: "",
            interfaceIp: "",
            interfaceName: interfaceName,
        };
        onOkay(rt, want);
        close();
    }

    $: {
        if (rt) {
            destination = rt.destination;
            mask = rt.mask;
            interfaceName = rt.interfaceName;
        }
        onChange();
    }
</script>

{#if show}
    <div class="modal">
        <div class="content">
            <div class="edit-title">{$i18n("addpage.title")}</div>
            <div class="input-box">
                <label for={destination}>{$i18n("table.destination")}</label>
                <input
                    class="input"
                    value={destination}
                    on:input={onDestinationChange}
                />
                <label for={mask}>{$i18n("table.mask")}</label>
                <input value={mask} class="input" on:input={onMaskChange} />
                <label for={interfaceName}
                    >{$i18n("table.interfacaName")}
                </label>
                <select
                    class="input"
                    value={interfaceName}
                    on:change={selectInterfaceName}
                >
                    {#each interfaceNames as interfaceName}
                        <option value={interfaceName}>
                            {interfaceName}
                        </option>
                    {/each}
                </select>
            </div>
            <button on:click={_onOkay} class="btn">
                {$i18n("addpage.configbutton")}
            </button>
            <button on:click={_onCancel} class="btn">
                {$i18n("addpage.cancelbuton")}
            </button>
            <slot />
        </div>
    </div>
{/if}
