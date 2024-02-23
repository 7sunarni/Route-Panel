<script lang="ts">
  import type { route } from "wailsjs/go/models.js";
  import { ListRoutes } from "../wailsjs/go/main/App.js";
  import { onMount } from "svelte";

  let searchIP: string = "";

  let currentPage: number = 0;

  let totalPagesRoutes: route.Route[][] = [];
  let currentPageRoutes: route.Route[] = [];

  let itemsPerPage: number = 10;
  let totalRoutes: route.Route[] = [];

  $: currentPageRoutes =
    totalPagesRoutes.length > 0 ? totalPagesRoutes[currentPage] : [];

  function deleteRow(r: route.Route) {
    totalRoutes = totalRoutes.filter(
      (row: route.Route) => row.destination != r.destination,
    );
  }

  let header = [
    "Destination",
    "Mask",
    "Gateway",
    "Interface Name",
    "Interface IP",
    "Metric",
    "Type",
    "Protocol",
    "Options",
  ];

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

  function search() {
    // todo
  }

  function refresh() {
    ListRoutes().then((result: route.Route[]) => {
      totalRoutes = result;
      paginate(totalRoutes);
      setPage(currentPage);
    });
  }
</script>

<main>
  <div class="input-box" id="input">
    <button class="btn" on:click={refresh}>Refresh</button>
    <input
      autocomplete="off"
      placeholder="Search IP"
      bind:value={searchIP}
      class="input"
      id="name"
      type="text"
    />
    <button class="btn" on:click={search}>Search</button>
  </div>
  <table>
    <tr>
      {#each header as column}
        <th style="background-color: green">{column}</th>
      {/each}
    </tr>

    {#each currentPageRoutes as route}
      <tr>
        <td contenteditable="false" bind:innerHTML={route.destination} />
        <td contenteditable="false" bind:innerHTML={route.mask} />
        <td contenteditable="false" bind:innerHTML={route.gateway} />
        <td contenteditable="false" bind:innerHTML={route.interfaceName} />
        <td contenteditable="false" bind:innerHTML={route.interfaceIp} />
        <td contenteditable="false" bind:innerHTML={route.metric} />
        <td contenteditable="false" bind:innerHTML={route.type} />
        <td contenteditable="false" bind:innerHTML={route.protocol} />
        <button on:click={() => deleteRow(route)}>Delete</button>
        <button>Edit</button>
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
          PREV
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
          NEXT
        </button>
      </li>
    </ul>
  </nav>
</main>

<style>
  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 40px 0 0 40px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin: 30px;
  }
  th,
  td {
    border: 1px solid black;
    padding: 8px;
    text-align: left;
  }
  th {
    background-color: #f2f2f2;
  }
  nav {
    justify-content: center; /* 水平居中 */
  }
  nav > ul {
    list-style-type: none;
    display: flex;
  }
</style>
