<script>
  import TargetProgress from '../charts/TargetProgress.svelte';
  import { money, int, initials, DASH } from '../format.js';

  /**
   * `topPerformer` from payload. The contract shows it nested under
   * `staffMember`, but it may also arrive as a plain person object.
   * We normalise both shapes.
   */
  export let topPerformer = null;

  // Normalise: a "person" object has staffMember/callMetrics/salesMetrics.
  // topPerformer may be { staffMember: <person> } or the person directly.
  $: person = (() => {
    if (!topPerformer) return null;
    const tp = topPerformer;
    // If it already looks like a person (has callMetrics/salesMetrics), use it.
    if (tp.callMetrics || tp.salesMetrics || tp.score != null) {
      // But staffMember could still be the nested person per the contract.
      if (
        tp.staffMember &&
        (tp.staffMember.callMetrics || tp.staffMember.salesMetrics)
      ) {
        return tp.staffMember;
      }
      return tp;
    }
    if (tp.staffMember) return tp.staffMember;
    return tp;
  })();

  $: name =
    (person && person.staffMember && person.staffMember.name) ||
    (person && person.name) ||
    DASH;
  $: score = person && person.score != null ? Number(person.score) : null;
  $: sales = person && person.salesMetrics ? person.salesMetrics : null;
  $: calls = person && person.callMetrics ? person.callMetrics : null;
  $: hasSales = !!(person && person.hasSalesData);
  $: hasCalls = !!(person && person.hasAircallData);
  $: target = sales ? Number(sales.target) || 0 : 0;
  $: progress = sales ? Number(sales.targetProgress) || 0 : 0;
</script>

<section class="card hero">
  <div class="hero-head">
    <span class="crown">👑</span>
    <span class="hero-title">Top Performer</span>
  </div>

  {#if !person}
    <div class="empty">No leader yet</div>
  {:else}
    <div class="hero-main">
      <div class="avatar">{initials(name)}</div>
      <div class="hero-id">
        <div class="hero-name">{name}</div>
        {#if score != null}
          <div class="hero-score tabular">
            {score.toFixed(1)}<span class="score-unit">pts</span>
          </div>
        {/if}
      </div>
      <div class="hero-gauge">
        <TargetProgress {progress} hasTarget={target > 0} />
        <div class="gauge-cap">of target</div>
      </div>
    </div>

    <div class="hero-stats">
      <div class="stat">
        <div class="stat-label">Sales</div>
        <div class="stat-val tabular">
          {hasSales ? money(sales.salesValue) : DASH}
        </div>
      </div>
      <div class="stat">
        <div class="stat-label">Deals</div>
        <div class="stat-val tabular">
          {hasSales ? int(sales.salesCount) : DASH}
        </div>
      </div>
      <div class="stat">
        <div class="stat-label">Calls</div>
        <div class="stat-val tabular">
          {hasCalls ? int(calls.totalCalls) : DASH}
        </div>
      </div>
    </div>
  {/if}
</section>

<style>
  .hero {
    display: flex;
    flex-direction: column;
    padding: clamp(0.7rem, 1.1vw, 1.3rem);
    gap: 0.6rem;
    min-height: 0;
    background: linear-gradient(
      160deg,
      rgba(245, 166, 35, 0.18),
      var(--card) 55%
    );
    border-color: rgba(255, 207, 63, 0.35);
  }
  .hero-head {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .crown {
    font-size: 1.3rem;
  }
  .hero-title {
    font-size: 0.85rem;
    font-weight: 800;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--gold);
  }
  .hero-main {
    display: flex;
    align-items: center;
    gap: 0.8rem;
  }
  .avatar {
    flex: 0 0 auto;
    width: clamp(2.6rem, 3.6vw, 4rem);
    height: clamp(2.6rem, 3.6vw, 4rem);
    border-radius: 50%;
    display: grid;
    place-items: center;
    font-weight: 900;
    font-size: 1.2rem;
    color: #1a1205;
    background: linear-gradient(150deg, #ffd95e, #f5a623);
    box-shadow: 0 0 1.4rem rgba(245, 166, 35, 0.5);
  }
  .hero-id {
    flex: 1 1 auto;
    min-width: 0;
  }
  .hero-name {
    font-size: clamp(1.2rem, 1.9vw, 2.1rem);
    font-weight: 900;
    color: #fff;
    line-height: 1.05;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .hero-score {
    font-size: clamp(1.6rem, 2.4vw, 2.8rem);
    font-weight: 900;
    color: var(--gold);
    line-height: 1;
  }
  .score-unit {
    font-size: 0.45em;
    color: var(--text-dim);
    margin-left: 0.2em;
    font-weight: 700;
  }
  .hero-gauge {
    flex: 0 0 auto;
    width: clamp(4rem, 6vw, 6.5rem);
    height: clamp(4rem, 6vw, 6.5rem);
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
  }
  .gauge-cap {
    margin-top: -0.4rem;
    font-size: 0.65rem;
    color: var(--text-faint);
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .hero-stats {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.6rem;
    margin-top: auto;
  }
  .stat {
    background: rgba(255, 255, 255, 0.04);
    border-radius: 0.6rem;
    padding: 0.5rem 0.6rem;
    text-align: center;
  }
  .stat-label {
    font-size: 0.7rem;
    color: var(--text-faint);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-weight: 700;
  }
  .stat-val {
    font-size: clamp(1.1rem, 1.6vw, 1.7rem);
    font-weight: 900;
    color: #fff;
  }
  .empty {
    color: var(--text-faint);
    font-size: 1.2rem;
    font-weight: 700;
    padding: 1rem 0;
  }
</style>
