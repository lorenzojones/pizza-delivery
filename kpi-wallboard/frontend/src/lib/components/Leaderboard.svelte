<script>
  import LeaderboardBars from '../charts/LeaderboardBars.svelte';

  /** Raw `people` array from the payload. */
  export let people = [];

  // Map the API shape into the flat rows the chart wants.
  $: rows = (people || [])
    .filter((p) => p && p.staffMember)
    .map((p) => ({
      id: p.staffMember.id || p.staffMember.name || Math.random().toString(36),
      rank: Number(p.rank) || 999,
      name: p.staffMember.name || p.staffMember.spreadsheetName || '—',
      score: Number(p.score) || 0,
      salesValue: p.salesMetrics ? Number(p.salesMetrics.salesValue) || 0 : 0,
      hasSalesData: !!p.hasSalesData,
    }));
</script>

<section class="card leaderboard">
  <div class="lb-head">
    <h2>Leaderboard</h2>
    <span class="lb-sub">ranked by performance score</span>
  </div>
  <div class="lb-body">
    {#if rows.length === 0}
      <div class="empty">No staff data yet</div>
    {:else}
      <LeaderboardBars people={rows} />
    {/if}
  </div>
</section>

<style>
  .leaderboard {
    display: flex;
    flex-direction: column;
    min-height: 0;
    height: 100%;
    padding: clamp(0.7rem, 1.2vw, 1.4rem);
  }
  .lb-head {
    display: flex;
    align-items: baseline;
    gap: 0.8rem;
    flex: 0 0 auto;
    margin-bottom: 0.4rem;
  }
  h2 {
    font-size: clamp(1.3rem, 1.9vw, 2.2rem);
    font-weight: 900;
    color: #fff;
  }
  .lb-sub {
    color: var(--text-faint);
    font-size: 0.95rem;
    font-weight: 600;
  }
  .lb-body {
    flex: 1 1 auto;
    min-height: 0;
    position: relative;
  }
  .empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--text-faint);
    font-size: 1.4rem;
    font-weight: 700;
  }
</style>
