<script>
  import { money, int } from '../format.js';

  /** Full dashboard payload. */
  export let data = null;
  /** Top-level stale flag. */
  export let stale = false;
  /** Is the connection currently online? */
  export let online = true;

  // Build a list of lively competitive messages from the payload.
  $: items = buildItems(data, stale, online);

  function personName(p) {
    if (!p) return null;
    if (p.staffMember && p.staffMember.name) return p.staffMember.name;
    if (p.name) return p.name;
    return null;
  }

  function buildItems(d, isStale, isOnline) {
    const out = [];
    if (!isOnline) {
      out.push('🔌 Connection lost — showing last known data');
    }
    if (isStale) {
      out.push('⚠️ Data is stale — awaiting fresh refresh');
    }
    if (!d) {
      out.push('⏳ Connecting to backend…');
      return out;
    }

    const people = Array.isArray(d.people) ? d.people : [];

    // Leader
    const leader = d.topPerformer
      ? (d.topPerformer.staffMember &&
        (d.topPerformer.staffMember.callMetrics ||
          d.topPerformer.staffMember.salesMetrics)
          ? d.topPerformer.staffMember
          : d.topPerformer)
      : people.find((p) => p.rank === 1);
    const leaderName = personName(leader);
    if (leaderName) {
      const sv =
        leader.salesMetrics && leader.salesMetrics.salesValue != null
          ? money(leader.salesMetrics.salesValue)
          : null;
      out.push(sv ? `🏆 ${leaderName} leads with ${sv}` : `🏆 ${leaderName} leads the board`);
    }

    // Biggest mover
    if (d.biggestMover && d.biggestMover.name) {
      const delta = Number(d.biggestMover.delta) || 0;
      if (delta > 0) {
        out.push(`📈 Biggest mover: ${d.biggestMover.name} (+${delta})`);
      } else if (delta < 0) {
        out.push(`📉 ${d.biggestMover.name} slipped ${delta} spot${delta === -1 ? '' : 's'}`);
      } else {
        out.push(`🔁 Biggest mover: ${d.biggestMover.name}`);
      }
    }

    // Team totals
    const tt = d.teamTotals;
    if (tt) {
      if (tt.salesValue != null)
        out.push(`💰 Team sales today: ${money(tt.salesValue)}`);
      if (tt.totalCalls != null)
        out.push(
          `📞 ${int(tt.totalCalls)} calls — ${int(tt.callsIn)} in / ${int(
            tt.callsOut
          )} out`
        );
      if (tt.salesCount != null)
        out.push(`✅ ${int(tt.salesCount)} deals closed by ${int(tt.headcount)} reps`);
    }

    // Runner-up callout
    const second = people.find((p) => p.rank === 2);
    if (second && leader) {
      const ls = Number(leader.score) || 0;
      const ss = Number(second.score) || 0;
      const gap = Math.max(0, ls - ss);
      if (gap > 0 && personName(second)) {
        out.push(`⚡ ${personName(second)} is ${gap.toFixed(1)} pts off the lead`);
      }
    }

    if (out.length === 0) out.push('📊 Sales wallboard live');
    return out;
  }
</script>

<div class="ticker" aria-label="Live stats ticker">
  <div class="ticker-tag">LIVE</div>
  <div class="ticker-viewport">
    <!-- Duplicate the sequence so the marquee loops seamlessly. -->
    <div class="ticker-track" class:paused={items.length === 0}>
      {#each [...items, ...items] as msg, i}
        <span class="ticker-item">{msg}</span>
        <span class="sep" aria-hidden="true">•</span>
      {/each}
    </div>
  </div>
</div>

<style>
  .ticker {
    display: flex;
    align-items: stretch;
    height: 100%;
    gap: 0;
    background: var(--bg-1);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    overflow: hidden;
  }
  .ticker-tag {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    padding: 0 1rem;
    font-weight: 900;
    letter-spacing: 0.14em;
    font-size: 0.85rem;
    color: #06210f;
    background: linear-gradient(120deg, var(--green), var(--accent-2));
  }
  .ticker-viewport {
    flex: 1 1 auto;
    overflow: hidden;
    position: relative;
    display: flex;
    align-items: center;
  }
  .ticker-track {
    display: inline-flex;
    align-items: center;
    white-space: nowrap;
    will-change: transform;
    animation: marquee 45s linear infinite;
  }
  .ticker-track.paused {
    animation: none;
  }
  .ticker-item {
    font-size: 1.25rem;
    font-weight: 800;
    color: var(--text);
    padding: 0 0.4rem;
  }
  .sep {
    color: var(--accent);
    padding: 0 0.9rem;
    font-size: 1.2rem;
  }
  @keyframes marquee {
    from {
      transform: translateX(0);
    }
    to {
      /* track contains the items twice; shift by half to loop seamlessly */
      transform: translateX(-50%);
    }
  }
</style>
