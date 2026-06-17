<script>
  import { onMount, onDestroy } from 'svelte';
  import * as d3 from 'd3';
  import { money, initials } from '../format.js';

  /**
   * @typedef {Object} Row
   * @property {string} id
   * @property {number} rank
   * @property {string} name
   * @property {number} score
   * @property {number} salesValue
   * @property {boolean} hasSalesData
   */
  /** @type {Row[]} */
  export let people = [];

  let container;
  let svgEl;
  let ro;
  let width = 0;
  let height = 0;

  // Keep transitions stable across redraws by keying on id.
  function draw() {
    if (!svgEl || width === 0 || height === 0) return;

    const data = (people || []).slice().sort((a, b) => a.rank - b.rank);
    const n = Math.max(1, data.length);

    const svg = d3.select(svgEl);
    svg.attr('width', width).attr('height', height);

    const rowH = height / n;
    const barH = Math.min(rowH * 0.62, 90);
    const padLeft = Math.min(width * 0.4, 320); // space for rank + name
    const padRight = Math.min(width * 0.26, 280); // space for the £ headline
    const trackX0 = padLeft;
    const trackX1 = width - padRight;
    const trackW = Math.max(10, trackX1 - trackX0);

    const maxScore = d3.max(data, (d) => d.score) || 1;
    const x = d3.scaleLinear().domain([0, maxScore]).range([0, trackW]).clamp(true);

    const T = svg.transition().duration(800).ease(d3.easeCubicOut);

    const rows = svg
      .selectAll('g.lb-row')
      .data(data, (d) => d.id);

    // EXIT
    rows.exit().transition(T).style('opacity', 0).remove();

    // ENTER
    const enter = rows
      .enter()
      .append('g')
      .attr('class', 'lb-row')
      .attr('transform', (d, i) => `translate(0, ${i * rowH})`)
      .style('opacity', 0);

    enter
      .append('rect')
      .attr('class', 'lb-track')
      .attr('x', trackX0)
      .attr('rx', 8)
      .attr('ry', 8);

    enter
      .append('rect')
      .attr('class', 'lb-bar')
      .attr('x', trackX0)
      .attr('rx', 8)
      .attr('ry', 8)
      .attr('width', 0);

    enter.append('text').attr('class', 'lb-rank');
    enter.append('text').attr('class', 'lb-name');
    enter.append('text').attr('class', 'lb-value');

    // MERGE
    const merged = enter.merge(rows);

    merged
      .transition(T)
      .attr('transform', (d, i) => `translate(0, ${i * rowH})`)
      .style('opacity', 1);

    const cy = rowH / 2;

    merged.classed('is-top', (d) => d.rank === 1);

    merged
      .select('rect.lb-track')
      .attr('y', cy - barH / 2)
      .attr('height', barH)
      .attr('width', trackW);

    merged
      .select('rect.lb-bar')
      .attr('y', cy - barH / 2)
      .attr('height', barH)
      .attr('fill', (d) =>
        d.rank === 1 ? 'url(#lbGold)' : 'url(#lbBlue)'
      )
      .transition(T)
      .attr('width', (d) => Math.max(2, x(d.score)));

    merged
      .select('text.lb-rank')
      .attr('x', padLeft * 0.13)
      .attr('y', cy)
      .attr('dy', '0.35em')
      .text((d) => `#${d.rank}`);

    merged
      .select('text.lb-name')
      .attr('x', padLeft * 0.34)
      .attr('y', cy)
      .attr('dy', '0.35em')
      .text((d) => d.name);

    merged
      .select('text.lb-value')
      .attr('x', width - padRight * 0.08)
      .attr('y', cy)
      .attr('dy', '0.35em')
      .attr('text-anchor', 'end')
      .text((d) => (d.hasSalesData ? money(d.salesValue) : '—'));
  }

  function measure() {
    if (!container) return;
    const r = container.getBoundingClientRect();
    width = r.width;
    height = r.height;
    draw();
  }

  onMount(() => {
    ro = new ResizeObserver(() => measure());
    ro.observe(container);
    measure();
  });

  onDestroy(() => {
    if (ro) ro.disconnect();
  });

  // Redraw on data change.
  $: people, draw();
</script>

<div class="lb-wrap" bind:this={container}>
  <svg bind:this={svgEl} role="img" aria-label="Leaderboard">
    <defs>
      <linearGradient id="lbGold" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0%" stop-color="#f5a623" />
        <stop offset="100%" stop-color="#ffd95e" />
      </linearGradient>
      <linearGradient id="lbBlue" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0%" stop-color="#3a6fd0" />
        <stop offset="100%" stop-color="#4f8cff" />
      </linearGradient>
    </defs>
  </svg>
</div>

<style>
  .lb-wrap {
    width: 100%;
    height: 100%;
    min-height: 0;
  }
  svg {
    display: block;
  }

  :global(.lb-row .lb-track) {
    fill: rgba(255, 255, 255, 0.05);
  }
  :global(.lb-row .lb-rank) {
    fill: var(--text-faint);
    font-weight: 800;
    font-size: 1.6rem;
    font-variant-numeric: tabular-nums;
  }
  :global(.lb-row.is-top .lb-rank) {
    fill: var(--gold);
  }
  :global(.lb-row .lb-name) {
    fill: var(--text);
    font-weight: 700;
    font-size: 1.5rem;
  }
  :global(.lb-row.is-top .lb-name) {
    fill: #fff;
    font-weight: 800;
  }
  :global(.lb-row .lb-value) {
    fill: var(--text);
    font-weight: 800;
    font-size: 1.7rem;
    font-variant-numeric: tabular-nums;
  }
  :global(.lb-row.is-top .lb-value) {
    fill: var(--gold);
  }
</style>
