<script>
  import { onMount, onDestroy } from 'svelte';
  import { dashboard, startPolling, stopPolling } from './lib/api.js';
  import Header from './lib/components/Header.svelte';
  import StaleBanner from './lib/components/StaleBanner.svelte';
  import Leaderboard from './lib/components/Leaderboard.svelte';
  import TopPerformerCard from './lib/components/TopPerformerCard.svelte';
  import TeamTotals from './lib/components/TeamTotals.svelte';
  import Ticker from './lib/components/Ticker.svelte';

  // Live "now" clock so relative times tick.
  let now = Date.now();
  let clockTimer;

  onMount(() => {
    startPolling();
    clockTimer = setInterval(() => (now = Date.now()), 1000);
  });

  onDestroy(() => {
    stopPolling();
    if (clockTimer) clearInterval(clockTimer);
  });

  $: state = $dashboard;
  $: data = state.data;
  $: loading = state.loading;
  $: online = state.online;
  $: error = state.error;

  // Derived flags
  $: title = (data && data.title) || 'Sales Wallboard';
  $: topLevelStale = !!(data && data.stale) || !online;
  $: sourceStatus = data && data.sourceStatus;
  $: demo = isDemo(sourceStatus);
  $: lastUpdated = (data && (data.lastUpdated || data.generatedAt)) || null;

  // First-load with no data yet -> connecting screen.
  $: connecting = loading && !data;
  $: firstFailure = !loading && !data;

  function isDemo(ss) {
    if (!ss) return false;
    const modes = [];
    if (ss.aircall && ss.aircall.mode) modes.push(String(ss.aircall.mode).toLowerCase());
    if (ss.excel && ss.excel.mode) modes.push(String(ss.excel.mode).toLowerCase());
    return modes.length > 0 && modes.some((m) => m === 'demo');
  }
</script>

<div class="wallboard">
  {#if connecting}
    <div class="boot">
      <div class="boot-spinner"></div>
      <div class="boot-title">Connecting to backend…</div>
      <div class="boot-sub">Starting the sales wallboard</div>
    </div>
  {:else if firstFailure && !data}
    <div class="boot">
      <div class="boot-spinner error"></div>
      <div class="boot-title">Waiting for backend…</div>
      <div class="boot-sub">
        Can't reach the API{error ? ` (${error})` : ''} — retrying automatically
      </div>
    </div>
  {:else}
    <header class="row-header">
      <Header
        {title}
        {lastUpdated}
        {sourceStatus}
        {online}
        {now}
        {demo}
      />
    </header>

    {#if topLevelStale}
      <div class="row-stale">
        <StaleBanner stale={data && data.stale} {online} {error} />
      </div>
    {/if}

    <main class="row-main">
      <div class="col-left">
        <Leaderboard people={data && data.people} />
      </div>
      <div class="col-right">
        <TopPerformerCard topPerformer={data && data.topPerformer} />
        <TeamTotals totals={data && data.teamTotals} />
      </div>
    </main>

    <footer class="row-ticker">
      <Ticker {data} stale={!!(data && data.stale)} {online} />
    </footer>
  {/if}
</div>

<style>
  .wallboard {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    gap: var(--gap);
    padding: var(--gap);
    overflow: hidden;
  }

  .row-header {
    flex: 0 0 auto;
    height: clamp(3.2rem, 7vh, 6rem);
  }
  .row-stale {
    flex: 0 0 auto;
  }
  .row-main {
    flex: 1 1 auto;
    min-height: 0;
    display: grid;
    grid-template-columns: minmax(0, 1.85fr) minmax(0, 1fr);
    gap: var(--gap);
  }
  .col-left {
    min-width: 0;
    min-height: 0;
    display: flex;
  }
  .col-left :global(.leaderboard) {
    flex: 1 1 auto;
  }
  .col-right {
    min-width: 0;
    min-height: 0;
    display: grid;
    grid-template-rows: minmax(0, auto) minmax(0, 1fr);
    gap: var(--gap);
  }
  .row-ticker {
    flex: 0 0 auto;
    height: clamp(2.6rem, 6vh, 4.5rem);
  }

  /* Boot / connecting state */
  .boot {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    text-align: center;
  }
  .boot-spinner {
    width: 4rem;
    height: 4rem;
    border-radius: 50%;
    border: 0.4rem solid rgba(255, 255, 255, 0.1);
    border-top-color: var(--accent);
    animation: spin 1s linear infinite;
  }
  .boot-spinner.error {
    border-top-color: var(--amber);
  }
  .boot-title {
    font-size: 2rem;
    font-weight: 900;
    color: #fff;
  }
  .boot-sub {
    font-size: 1.1rem;
    color: var(--text-dim);
  }
  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
