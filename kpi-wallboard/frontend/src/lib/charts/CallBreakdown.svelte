<script>
  import { onMount, onDestroy } from 'svelte';
  import * as d3 from 'd3';
  import { int } from '../format.js';

  /** Team-level inbound vs outbound call counts. */
  export let callsIn = 0;
  export let callsOut = 0;

  let container;
  let svgEl;
  let ro;
  let width = 0;
  let height = 0;

  const COLORS = { in: 'var(--callin)', out: 'var(--callout)' };

  let arcGen, prevArcs;

  function draw() {
    if (!svgEl || width === 0 || height === 0) return;

    const cin = Math.max(0, Number(callsIn) || 0);
    const cout = Math.max(0, Number(callsOut) || 0);
    const total = cin + cout;

    const svg = d3.select(svgEl);
    svg.attr('width', width).attr('height', height);

    const size = Math.min(width, height);
    const radius = size / 2;
    const thickness = Math.max(14, radius * 0.34);
    const cx = width / 2;
    const cy = height / 2;

    const pieData = [
      { key: 'out', label: 'Out', value: cout, color: 'var(--callout)' },
      { key: 'in', label: 'In', value: cin, color: 'var(--callin)' },
    ];

    // Avoid an empty pie collapsing.
    const pie = d3
      .pie()
      .sort(null)
      .value((d) => d.value)
      .padAngle(0.02);

    const arcs = pie(total === 0 ? pieData.map((d) => ({ ...d, value: 1 })) : pieData);

    arcGen = d3
      .arc()
      .innerRadius(radius - thickness)
      .outerRadius(radius)
      .cornerRadius(Math.min(8, thickness / 2));

    let g = svg.select('g.cb-root');
    if (g.empty()) {
      g = svg.append('g').attr('class', 'cb-root');
    }
    g.attr('transform', `translate(${cx}, ${cy})`);

    const T = svg.transition().duration(800).ease(d3.easeCubicOut);

    const paths = g.selectAll('path.cb-arc').data(arcs, (d) => d.data.key);

    paths
      .enter()
      .append('path')
      .attr('class', 'cb-arc')
      .attr('fill', (d) => d.data.color)
      .each(function (d) {
        this._current = { ...d, startAngle: d.startAngle, endAngle: d.startAngle };
      })
      .merge(paths)
      .transition(T)
      .attrTween('d', function (d) {
        const i = d3.interpolate(this._current || d, d);
        this._current = i(1);
        return (t) => arcGen(i(t));
      });

    paths.exit().remove();

    // Center label: total calls.
    let label = g.select('g.cb-center');
    if (label.empty()) {
      label = g.append('g').attr('class', 'cb-center');
      label.append('text').attr('class', 'cb-total').attr('text-anchor', 'middle');
      label
        .append('text')
        .attr('class', 'cb-total-label')
        .attr('text-anchor', 'middle')
        .text('TOTAL CALLS');
    }
    label
      .select('.cb-total')
      .attr('dy', '-0.05em')
      .style('font-size', `${Math.max(18, radius * 0.42)}px`)
      .text(int(total));
    label
      .select('.cb-total-label')
      .attr('dy', `${Math.max(14, radius * 0.34)}px`)
      .style('font-size', `${Math.max(9, radius * 0.13)}px`);
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

  $: callsIn, callsOut, draw();

  $: total = (Number(callsIn) || 0) + (Number(callsOut) || 0);
</script>

<div class="cb-wrap">
  <div class="cb-chart" bind:this={container}>
    <svg bind:this={svgEl} role="img" aria-label="Call volume breakdown"></svg>
  </div>
  <div class="cb-legend">
    <div class="cb-leg-item">
      <span class="cb-dot" style="background: var(--callout)"></span>
      <span class="cb-leg-label">Out</span>
      <span class="cb-leg-val tabular">{int(callsOut)}</span>
    </div>
    <div class="cb-leg-item">
      <span class="cb-dot" style="background: var(--callin)"></span>
      <span class="cb-leg-label">In</span>
      <span class="cb-leg-val tabular">{int(callsIn)}</span>
    </div>
  </div>
</div>

<style>
  .cb-wrap {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    min-height: 0;
    gap: 0.3rem;
  }
  .cb-chart {
    flex: 1 1 auto;
    min-height: 0;
    width: 100%;
  }
  svg {
    display: block;
  }
  .cb-legend {
    flex: 0 0 auto;
    display: flex;
    justify-content: center;
    gap: 1.4rem;
  }
  .cb-leg-item {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.95rem;
  }
  .cb-dot {
    width: 0.85rem;
    height: 0.85rem;
    border-radius: 3px;
    display: inline-block;
  }
  .cb-leg-label {
    color: var(--text-dim);
  }
  .cb-leg-val {
    color: var(--text);
    font-weight: 800;
  }

  :global(.cb-total) {
    fill: var(--text);
    font-weight: 800;
    font-variant-numeric: tabular-nums;
  }
  :global(.cb-total-label) {
    fill: var(--text-faint);
    font-weight: 700;
    letter-spacing: 0.08em;
  }
</style>
