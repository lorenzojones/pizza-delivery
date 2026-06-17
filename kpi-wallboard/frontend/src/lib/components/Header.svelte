<script>
  import SourceBadge from './SourceBadge.svelte';
  import { longDate, clock, relativeAgo } from '../format.js';

  export let title = 'Vent Sales';
  export let lastUpdated = null;
  export let sourceStatus = null;
  export let online = true;
  /** ms epoch, ticks every second to keep "x ago" fresh. */
  export let now = Date.now();
  /** Are any sources in demo mode? */
  export let demo = false;

  $: aircall = sourceStatus && sourceStatus.aircall;
  $: excel = sourceStatus && sourceStatus.excel;
  $: rel = lastUpdated ? relativeAgo(lastUpdated, now) : '—';
  $: abs = lastUpdated ? clock(lastUpdated) : '—';
</script>

<header class="hdr">
  <div class="left">
    <div class="title-row">
      <span class="logo" aria-hidden="true">▲</span>
      <h1>{title}</h1>
      {#if demo}
        <span class="demo-pill">DEMO DATA</span>
      {/if}
    </div>
    <div class="date">{longDate(now)}</div>
  </div>

  <div class="right">
    <div class="updated" class:offline={!online}>
      <span class="upd-label">Updated</span>
      <span class="upd-rel tabular">{rel}</span>
      <span class="upd-abs tabular">{abs}</span>
      {#if !online}
        <span class="recon">reconnecting…</span>
      {/if}
    </div>
    <div class="badges">
      <SourceBadge source={aircall} label="Aircall" />
      <SourceBadge source={excel} label="Excel" />
    </div>
  </div>
</header>

<style>
  .hdr {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--gap);
    height: 100%;
    padding: 0 0.4rem;
  }
  .left {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    min-width: 0;
  }
  .title-row {
    display: flex;
    align-items: center;
    gap: 0.7rem;
  }
  .logo {
    color: var(--accent);
    font-size: 1.6rem;
    line-height: 1;
  }
  h1 {
    font-size: clamp(1.6rem, 2.6vw, 3rem);
    font-weight: 900;
    letter-spacing: -0.01em;
    color: #fff;
    white-space: nowrap;
  }
  .demo-pill {
    font-size: 0.7rem;
    font-weight: 800;
    letter-spacing: 0.1em;
    color: var(--amber);
    background: rgba(255, 177, 61, 0.16);
    border: 1px solid rgba(255, 177, 61, 0.4);
    padding: 0.25rem 0.55rem;
    border-radius: 6px;
  }
  .date {
    color: var(--text-dim);
    font-size: 1.05rem;
    font-weight: 600;
  }
  .right {
    display: flex;
    align-items: center;
    gap: 1.2rem;
    flex: 0 0 auto;
  }
  .updated {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    line-height: 1.15;
  }
  .upd-label {
    font-size: 0.7rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--text-faint);
  }
  .upd-rel {
    font-size: 1.15rem;
    font-weight: 800;
    color: var(--text);
  }
  .upd-abs {
    font-size: 0.85rem;
    color: var(--text-dim);
  }
  .updated.offline .upd-rel {
    color: var(--amber);
  }
  .recon {
    font-size: 0.72rem;
    color: var(--amber);
    font-weight: 700;
  }
  .badges {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    align-items: flex-end;
  }
</style>
