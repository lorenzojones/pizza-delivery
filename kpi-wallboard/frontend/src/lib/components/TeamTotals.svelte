<script>
  import CallBreakdown from '../charts/CallBreakdown.svelte';
  import { money, int } from '../format.js';

  /** `teamTotals` from payload. */
  export let totals = null;

  $: t = totals || {};
  $: callsIn = Number(t.callsIn) || 0;
  $: callsOut = Number(t.callsOut) || 0;
</script>

<section class="card totals">
  <div class="t-head">
    <h2>Team Totals</h2>
    <span class="headcount tabular">{int(t.headcount)} <span>people</span></span>
  </div>

  <div class="t-grid">
    <div class="kpi big">
      <div class="kpi-label">Total Sales</div>
      <div class="kpi-val tabular">{money(t.salesValue)}</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Deals</div>
      <div class="kpi-val tabular">{int(t.salesCount)}</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Total Calls</div>
      <div class="kpi-val tabular">{int(t.totalCalls)}</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Minutes</div>
      <div class="kpi-val tabular">{int(t.minutesOnPhone)}</div>
    </div>
  </div>

  <div class="t-chart">
    <CallBreakdown {callsIn} {callsOut} />
  </div>
</section>

<style>
  .totals {
    display: flex;
    flex-direction: column;
    min-height: 0;
    padding: clamp(0.7rem, 1.1vw, 1.3rem);
    gap: 0.6rem;
  }
  .t-head {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    flex: 0 0 auto;
  }
  h2 {
    font-size: clamp(1.1rem, 1.6vw, 1.9rem);
    font-weight: 900;
    color: #fff;
  }
  .headcount {
    color: var(--accent-2);
    font-weight: 900;
    font-size: 1.1rem;
  }
  .headcount span {
    color: var(--text-faint);
    font-weight: 600;
    font-size: 0.8rem;
  }
  .t-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5rem;
    flex: 0 0 auto;
  }
  .kpi {
    background: rgba(255, 255, 255, 0.04);
    border-radius: 0.6rem;
    padding: 0.5rem 0.7rem;
  }
  .kpi.big {
    grid-column: 1 / -1;
    background: linear-gradient(
      120deg,
      rgba(56, 224, 192, 0.16),
      rgba(255, 255, 255, 0.03)
    );
  }
  .kpi-label {
    font-size: 0.72rem;
    color: var(--text-faint);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-weight: 700;
  }
  .kpi-val {
    font-weight: 900;
    color: #fff;
    font-size: clamp(1.1rem, 1.6vw, 1.8rem);
    line-height: 1.1;
  }
  .kpi.big .kpi-val {
    font-size: clamp(1.7rem, 2.6vw, 3rem);
    color: var(--accent-2);
  }
  .t-chart {
    flex: 1 1 auto;
    min-height: 0;
  }
</style>
