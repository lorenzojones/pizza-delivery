<script>
  /** Top-level stale flag from payload. */
  export let stale = false;
  /** True when the last fetch attempt failed. */
  export let online = true;
  /** Last fetch error message, if any. */
  export let error = null;

  // Show whenever data is stale OR we're currently offline.
  $: show = stale || !online;
  $: message = !online
    ? 'Connection lost — showing last known data'
    : 'Data is stale — last refresh did not update';
</script>

{#if show}
  <div class="stale-banner" role="status">
    <span class="icon">⚠</span>
    <span class="msg">{message}</span>
    {#if error && !online}
      <span class="detail">({error})</span>
    {/if}
  </div>
{/if}

<style>
  .stale-banner {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.45rem 1rem;
    background: linear-gradient(
      90deg,
      rgba(255, 84, 104, 0.22),
      rgba(255, 177, 61, 0.18)
    );
    border: 1px solid rgba(255, 177, 61, 0.45);
    border-radius: var(--radius);
    color: #fff;
    font-weight: 700;
    font-size: 1rem;
    animation: glow 2s ease-in-out infinite;
  }
  .icon {
    color: var(--amber);
    font-size: 1.2rem;
  }
  .detail {
    color: var(--text-dim);
    font-weight: 600;
    font-size: 0.85rem;
  }
  @keyframes glow {
    0%,
    100% {
      box-shadow: 0 0 0 rgba(255, 177, 61, 0);
    }
    50% {
      box-shadow: 0 0 1.2rem rgba(255, 177, 61, 0.4);
    }
  }
</style>
