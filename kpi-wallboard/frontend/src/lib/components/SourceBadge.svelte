<script>
  /**
   * @type {{
   *  name?: string, ok?: boolean, stale?: boolean,
   *  lastSuccess?: string, lastError?: string, mode?: string
   * } | null}
   */
  export let source = null;
  /** Fallback label if source.name missing. */
  export let label = '';

  $: name = (source && source.name) || label || 'Source';
  $: ok = !!(source && source.ok);
  $: stale = !!(source && source.stale);
  // status: red (error / !ok), amber (stale), green (ok & fresh)
  $: status = !ok ? 'error' : stale ? 'stale' : 'ok';
  $: mode = (source && source.mode ? String(source.mode) : '').toUpperCase();
  $: statusText =
    status === 'error' ? 'ERROR' : status === 'stale' ? 'STALE' : 'LIVE';
  $: title =
    source && source.lastError
      ? `${name}: ${source.lastError}`
      : `${name}: ${statusText}`;
</script>

<div class="badge {status}" {title}>
  <span class="dot"></span>
  <span class="name">{name}</span>
  <span class="status">{statusText}</span>
  {#if mode}
    <span class="mode" class:demo={mode === 'DEMO'}>{mode}</span>
  {/if}
</div>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.35rem 0.7rem;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--bg-1);
    font-size: 0.85rem;
    font-weight: 700;
    line-height: 1;
    white-space: nowrap;
  }
  .name {
    color: var(--text);
  }
  .status {
    font-size: 0.72rem;
    letter-spacing: 0.06em;
  }
  .dot {
    width: 0.7rem;
    height: 0.7rem;
    border-radius: 50%;
    flex: 0 0 auto;
    box-shadow: 0 0 0.5rem currentColor;
  }
  .ok .dot {
    background: var(--green);
    color: var(--green);
  }
  .ok .status {
    color: var(--green);
  }
  .stale .dot {
    background: var(--amber);
    color: var(--amber);
    animation: pulse 1.4s ease-in-out infinite;
  }
  .stale .status {
    color: var(--amber);
  }
  .error .dot {
    background: var(--red);
    color: var(--red);
    animation: pulse 1.1s ease-in-out infinite;
  }
  .error .status {
    color: var(--red);
  }
  .mode {
    font-size: 0.66rem;
    letter-spacing: 0.06em;
    padding: 0.12rem 0.4rem;
    border-radius: 5px;
    background: rgba(255, 255, 255, 0.08);
    color: var(--text-dim);
  }
  .mode.demo {
    background: rgba(255, 177, 61, 0.18);
    color: var(--amber);
  }
  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.35;
    }
  }
</style>
