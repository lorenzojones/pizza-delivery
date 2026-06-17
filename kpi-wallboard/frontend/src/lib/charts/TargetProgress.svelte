<script>
  import { onMount, onDestroy } from 'svelte';
  import * as d3 from 'd3';
  import { pct } from '../format.js';

  /** Progress 0..1+ toward target. May exceed 1. */
  export let progress = 0;
  /** Whether a target exists at all (target > 0). */
  export let hasTarget = true;

  let container;
  let svgEl;
  let ro;
  let width = 0;
  let height = 0;

  const START = -Math.PI * 0.75;
  const END = Math.PI * 0.75; // 270deg open-bottom gauge

  function color(p) {
    if (p >= 1) return 'var(--gold)';
    if (p >= 0.66) return 'var(--green)';
    if (p >= 0.33) return 'var(--amber)';
    return 'var(--red)';
  }

  function draw() {
    if (!svgEl || width === 0 || height === 0) return;

    const p = Math.max(0, Number(progress) || 0);
    const shown = Math.min(1, p); // arc fill caps at full sweep

    const svg = d3.select(svgEl);
    svg.attr('width', width).attr('height', height);

    const size = Math.min(width, height);
    const radius = size / 2 - 4;
    const thickness = Math.max(10, radius * 0.22);
    const cx = width / 2;
    const cy = height / 2;

    let g = svg.select('g.tp-root');
    if (g.empty()) g = svg.append('g').attr('class', 'tp-root');
    g.attr('transform', `translate(${cx}, ${cy})`);

    const arc = d3
      .arc()
      .innerRadius(radius - thickness)
      .outerRadius(radius)
      .cornerRadius(thickness / 2)
      .startAngle(START);

    // Track (full sweep, dim).
    let track = g.select('path.tp-track');
    if (track.empty()) track = g.append('path').attr('class', 'tp-track');
    track.attr('d', arc.endAngle(END)());

    // Value arc.
    let val = g.select('path.tp-val');
    if (val.empty()) {
      val = g.append('path').attr('class', 'tp-val');
      val._cur = START;
    }
    val.attr('fill', hasTarget ? color(p) : 'var(--text-faint)');

    const targetAngle = START + (END - START) * (hasTarget ? shown : 0);
    const fromAngle = val._cur != null ? val._cur : START;

    val
      .transition()
      .duration(900)
      .ease(d3.easeCubicOut)
      .attrTween('d', function () {
        const i = d3.interpolate(fromAngle, targetAngle);
        return (t) => {
          const a = i(t);
          val._cur = a;
          return arc.endAngle(a)();
        };
      });

    // Center text.
    let label = g.select('text.tp-label');
    if (label.empty()) {
      label = g
        .append('text')
        .attr('class', 'tp-label')
        .attr('text-anchor', 'middle')
        .attr('dy', '0.35em');
    }
    label
      .style('font-size', `${Math.max(16, radius * 0.5)}px`)
      .attr('fill', hasTarget ? color(p) : 'var(--text-faint)')
      .text(hasTarget ? pct(p) : '—');
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

  $: progress, hasTarget, draw();
</script>

<div class="tp-wrap" bind:this={container}>
  <svg bind:this={svgEl} role="img" aria-label="Target progress"></svg>
</div>

<style>
  .tp-wrap {
    width: 100%;
    height: 100%;
    min-height: 0;
  }
  svg {
    display: block;
  }
  :global(.tp-track) {
    fill: rgba(255, 255, 255, 0.07);
  }
  :global(.tp-label) {
    font-weight: 800;
    font-variant-numeric: tabular-nums;
  }
</style>
