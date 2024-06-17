<script lang="ts">
  import type { route } from "../wailsjs/go/models";
  import {
    AddRoute,
    DeleteRoute,
    ListRoutes,
    EditRoute,
  } from "../wailsjs/go/main/App.js";
  import { onMount } from "svelte";
  import Routeview from "./route/Route.svelte";
  import { i18n, locale, locales } from "./i18n/i18n.js";
  import { toast } from "@zerodevx/svelte-toast";
  import { SvelteToast } from "@zerodevx/svelte-toast";

  // Optionally set default options here
  const options = {};

  const itemsPerPage: number = 10;

  let searched: string = "";

  let currentPage: number = 0;

  let totalPagesRoutes: route.Route[][] = [];
  let currentPageRoutes: route.Route[] = [];

  let totalRoutes: route.Route[] = [];

  let showEdit: boolean = false;
  let choosenRoute: route.Route;
  let interfaceNames: { [key: string]: string } = {};
  let filteredRoutes: route.Route[] = [];

  $: currentPageRoutes =
    totalPagesRoutes.length > 0 ? totalPagesRoutes[currentPage] : [];

  function deleteRoute(r: route.Route) {
    DeleteRoute(r).then((msg: string) => {
      if (msg !== "") {
        toast.push(msg);
        return;
      }
      toast.push("Delete route success", {
        theme: {
          "--toastColor": "mintcream",
          "--toastBackground": "rgba(72,187,120,0.9)",
          "--toastBarBackground": "#2F855A",
          "--toastBarHeight": 0,
        },
      });
      refresh();
    });
  }

  function addRoute(r: route.Route) {
    AddRoute(r).then((msg: string) => {
      if (msg !== "") {
        toast.push(msg);
        return;
      }
      toast.push("Add route success", {
        theme: {
          "--toastColor": "mintcream",
          "--toastBackground": "rgba(72,187,120,0.9)",
          "--toastBarBackground": "#2F855A",
          "--toastBarHeight": 0,
        },
      });
      refresh();
    });
  }

  function editRoute(old: route.Route, want: route.Route) {
    EditRoute(old, want).then((msg: string) => {
      if (msg !== "") {
        toast.push(msg);
        return;
      }
      toast.push("Edit route success", {
        theme: {
          "--toastColor": "mintcream",
          "--toastBackground": "rgba(72,187,120,0.9)",
          "--toastBarBackground": "#2F855A",
          "--toastBarHeight": 0,
        },
      });
      refresh();
    });
  }

  function onSearch(event: any) {
    let search = event.target.value;
    let filter: route.Route[] = [];

    if (search !== "") {
      totalRoutes.forEach((r) => {
        if (
          r.destination.includes(search) ||
          r.gateway.includes(search) ||
          r.interfaceName.includes(search)
        ) {
          filter.push(r);
        }
      });
    } else {
      filter = totalRoutes;
    }

    filteredRoutes = filter;
    paginate(filteredRoutes);
    setPage(currentPage);
  }

  function onEditRoute(choosen: route.Route) {
    choosenRoute = choosen;
    setShowEdit(true);
  }

  const paginate = (rs: route.Route[]) => {
    const pages = Math.ceil(rs.length / itemsPerPage);
    const paginatedItems = Array.from({ length: pages }, (_, index) => {
      const start = index * itemsPerPage;
      return rs.slice(start, start + itemsPerPage);
    });
    totalPagesRoutes = [...paginatedItems];
  };

  onMount(() => {
    paginate(totalRoutes);
  });

  const setPage = (p: number) => {
    if (p >= 0) {
      currentPage = p;
    }
  };

  const setShowEdit = (e: boolean) => {
    showEdit = e;
  };

  function onSaveEdit(old: route.Route, want: route.Route) {
    want.interfaceIp = interfaceNames[want.interfaceName];
    if (old && old.destination !== "") {
      editRoute(old, want);
    } else {
      addRoute(want);
    }
    setShowEdit(false);
  }

  function refresh() {
    ListRoutes().then((result: route.Route[]) => {
      let filteNonInterface: route.Route[] = [];
      result.forEach((r: route.Route) => {
        if (r.interfaceName === "") {
          return;
        }
        filteNonInterface.push(r);
        if (r.interfaceName in interfaceNames) {
          return;
        }
        interfaceNames[r.interfaceName] = r.interfaceIp;
      });

      totalRoutes = filteNonInterface;
      paginate(totalRoutes);
      setPage(currentPage);
    });
  }
</script>

<main>
  <SvelteToast {options} />
  <Routeview
    rt={choosenRoute}
    show={showEdit}
    interfaceNames={Object.keys(interfaceNames)}
    onOkay={onSaveEdit}
    onCancel={() => setShowEdit(false)}
  />

  <div class="input-box" id="input">
    <input
      autocomplete="off"
      placeholder={$i18n("homepage.search")}
      on:input={onSearch}
      bind:value={searched}
      class="input"
      id="name"
      type="text"
    />
    <button class="btn" on:click={refresh}>{$i18n("homepage.refresh")}</button>
    <button
      class="btn"
      on:click={() =>
        onEditRoute({
          destination: "",
          mask: "",
          gateway: "",
          interfaceIp: "",
          interfaceName: "",
        })}>{$i18n("homepage.addRoute")}</button
    >
    <select bind:value={$locale}>
      {#each locales as l}
        <option value={l}>{l}</option>
      {/each}
    </select>
  </div>
  <table>
    <tr>
      <th>{$i18n("table.destination")}</th>
      <th>{$i18n("table.mask")}</th>
      <th>{$i18n("table.gateway")}</th>
      <th>{$i18n("table.interfacaName")}</th>
      <th>{$i18n("table.options")}</th>
    </tr>

    {#each currentPageRoutes as route}
      <tr>
        <td contenteditable="false" bind:innerHTML={route.destination} />
        <td contenteditable="false" bind:innerHTML={route.mask} />
        <td contenteditable="false" bind:innerHTML={route.gateway} />
        <td contenteditable="false" bind:innerHTML={route.interfaceName} />
        <td>
          <button on:click={() => onEditRoute(route)}
            >{$i18n("table.edit")}</button
          >
          <button on:click={() => deleteRoute(route)}
            >{$i18n("table.delete")}</button
          >
        </td>
      </tr>
    {/each}
  </table>
  <nav class="pagination">
    <ul>
      <li>
        <button
          type="button"
          class="btn-next-prev"
          on:click={() => setPage(currentPage - 1)}
        >
          {$i18n("page.prev")}
        </button>
      </li>

      {#each totalPagesRoutes as page, i}
        <li>
          <button
            type="button"
            class="btn-page-number"
            on:click={() => setPage(i)}
          >
            {i + 1}
          </button>
        </li>
      {/each}

      <li>
        <button
          type="button"
          class="btn-next-prev"
          on:click={() => setPage(currentPage + 1)}
        >
          {$i18n("page.next")}
        </button>
      </li>
    </ul>
  </nav>
</main>
