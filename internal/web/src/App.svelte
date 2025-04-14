<script>
  import { metricsStore } from "./stores/metrics";
  import WebsocketHandler from "./components/Websocket.svelte";
  import Plot from "./components/Plot.svelte";
</script>

<main>
  <WebsocketHandler />
  <!-- the series keys should come from the config store, not from the metrics, otherwise it will
   be refreshed each time the metrics store is updated. -->

  {#if $metricsStore}
    {#each Array.from($metricsStore.series.keys()) as key}
      <Plot name={key} />
    {/each}
  {/if} <button on:click={onclick}> Set key-a </button>
</main>

<style>
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>
