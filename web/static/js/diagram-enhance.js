/*!
 * diagram-enhance.js
 * Enhancement layer for CSS-based diagrams. Zero diagram file changes.
 * Features: smart tooltips · counter animation · step glow · entrance animation
 *           flow pulse · legend toggle · keyboard shortcuts
 */
(function () {
  'use strict';

  /* ═══════════════════════════════════════════════════════════════
     1. SMART TOOLTIP
     Replaces CSS ::after tooltips with a JS-positioned overlay that:
     - Never clips (fixed position, viewport-aware)
     - Animates in smoothly
     - Follows cursor
  ══════════════════════════════════════════════════════════════════ */
  const tip = document.createElement('div');
  tip.id = 'dg-tooltip';
  tip.style.cssText = [
    'position:fixed',
    'z-index:99999',
    'max-width:300px',
    'padding:9px 13px',
    'background:#0F172A',
    'color:#E2E8F0',
    'font-size:0.71rem',
    'font-family:var(--font-sans,sans-serif)',
    'font-weight:450',
    'line-height:1.5',
    'border-radius:7px',
    'box-shadow:0 10px 30px rgba(0,0,0,0.35)',
    'pointer-events:none',
    'opacity:0',
    'transition:opacity 0.15s ease,transform 0.15s ease',
    'transform:translateY(6px)',
    'white-space:normal',
    'word-break:break-word',
    'border:1px solid rgba(255,255,255,0.07)',
  ].join(';');
  document.body.appendChild(tip);

  let tipHideTimer, tipShowTimer;

  function showTip(el, e) {
    clearTimeout(tipHideTimer);
    clearTimeout(tipShowTimer);
    const text = el.getAttribute('data-tip');
    if (!text) return;
    tipShowTimer = setTimeout(() => {
      tip.textContent = text;
      placeTip(e);
      tip.style.opacity = '1';
      tip.style.transform = 'translateY(0)';
    }, 80);
  }

  function moveTip(e) {
    if (tip.style.opacity === '1') placeTip(e);
  }

  function placeTip(e) {
    const GAP = 14;
    const tw = tip.offsetWidth || 260;
    const th = tip.offsetHeight || 60;
    const vw = window.innerWidth;
    const vh = window.innerHeight;
    let left = e.clientX + GAP;
    let top  = e.clientY - Math.round(th / 2);
    if (left + tw > vw - GAP) left = e.clientX - tw - GAP;
    if (top < GAP)             top  = GAP;
    if (top + th > vh - GAP)   top  = vh - th - GAP;
    tip.style.left = left + 'px';
    tip.style.top  = top  + 'px';
  }

  function hideTip() {
    clearTimeout(tipShowTimer);
    tipHideTimer = setTimeout(() => {
      tip.style.opacity = '0';
      tip.style.transform = 'translateY(6px)';
    }, 60);
  }

  /* Disable the CSS ::after tooltip only on d-box elements we handle */
  function injectTipDisableCSS() {
    if (document.getElementById('dg-tip-css-off')) return;
    const s = document.createElement('style');
    s.id = 'dg-tip-css-off';
    s.textContent = '.dg-js-tip.d-box::after,.dg-js-tip.d-box::before{display:none!important}';
    document.head.appendChild(s);
  }

  function initTooltips(root) {
    injectTipDisableCSS();
    root.querySelectorAll('[data-tip]').forEach(el => {
      if (el.dataset.dgTip) return; // already initialized
      el.dataset.dgTip = '1';
      if (el.classList.contains('d-box')) el.classList.add('dg-js-tip');
      el.addEventListener('mouseenter', e => showTip(el, e));
      el.addEventListener('mousemove',  moveTip);
      el.addEventListener('mouseleave', hideTip);
    });
  }


  /* ═══════════════════════════════════════════════════════════════
     2. STEP GLOW
     Hover a d-step badge → all boxes sharing that step number in the
     same diagram get a colored ring. Helps readers trace multi-column flows.
  ══════════════════════════════════════════════════════════════════ */
  function initStepGlow(root) {
    root.querySelectorAll('.d-step').forEach(badge => {
      if (badge.dataset.dgStep) return;
      badge.dataset.dgStep = '1';
      const diagram = badge.closest('.diagram-container');
      if (!diagram) return;

      const stepId = badge.textContent.trim();
      badge.style.cursor = 'pointer';
      badge.title = `Step ${stepId} — hover to highlight`;

      badge.addEventListener('mouseenter', () => {
        diagram.querySelectorAll('.d-step').forEach(b => {
          const box = b.closest('.d-box');
          if (!box) return;
          const match = b.textContent.trim() === stepId;
          box.style.transition = 'box-shadow 0.2s ease, transform 0.2s ease';
          if (match) {
            box.style.boxShadow = '0 0 0 2px var(--indigo),0 6px 20px rgba(99,102,241,0.3)';
            box.style.transform = 'translateY(-2px)';
          } else {
            box.style.opacity = '0.45';
          }
        });
      });

      badge.addEventListener('mouseleave', () => {
        diagram.querySelectorAll('.d-box').forEach(box => {
          box.style.boxShadow = '';
          box.style.transform = '';
          box.style.opacity   = '';
        });
      });
    });
  }


  /* ═══════════════════════════════════════════════════════════════
     3. ENTRANCE ANIMATION
     IntersectionObserver: stagger-reveal d-box / d-group / d-entity
     when diagram scrolls into view. Runs once per diagram.
  ══════════════════════════════════════════════════════════════════ */
  const ENTER_TARGETS = '.d-box,.d-group,.d-entity,.d-number,.d-subproblem';

  const enterObs = new IntersectionObserver(entries => {
    entries.forEach(entry => {
      if (!entry.isIntersecting) return;
      const el = entry.target;
      el.querySelectorAll(ENTER_TARGETS).forEach((item, i) => {
        // Only animate items that haven't been seen yet
        if (item.dataset.dgEntered) return;
        item.dataset.dgEntered = '1';
        const delay = i * 30; // 30ms stagger
        item.style.opacity  = '0';
        item.style.transform = 'translateY(8px)';
        item.style.transition = 'none';
        requestAnimationFrame(() => {
          setTimeout(() => {
            item.style.transition = `opacity 0.35s ease ${delay}ms,transform 0.35s ease ${delay}ms`;
            item.style.opacity  = '';
            item.style.transform = '';
          }, 40);
        });
      });
      enterObs.unobserve(el);
    });
  }, { threshold: 0.08 });

  function initEntrance(root) {
    root.querySelectorAll('.diagram-container').forEach(c => {
      if (!c.dataset.dgObserved) {
        c.dataset.dgObserved = '1';
        enterObs.observe(c);
      }
    });
  }


  /* ═══════════════════════════════════════════════════════════════
     4. COUNTER ANIMATION
     d-number-value elements count up from 0 when they enter view.
     Handles: integers, decimals, K/M/B suffix numbers, < > ~ prefixes.
  ══════════════════════════════════════════════════════════════════ */
  function parseCounterVal(text) {
    // Matches: optional prefix (<, ≈, ~, $) + number + optional suffix (K, M, B, ms, %, +, x, s)
    const m = text.trim().match(/^([<≈~>$]?\s*)(\d[\d,.]*)([KMBTkms%+×x]*)\s*$/);
    if (!m) return null;
    const raw = parseFloat(m[2].replace(/,/g, ''));
    if (isNaN(raw)) return null;
    return { prefix: m[1], num: raw, suffix: m[3], original: text, decimals: (m[2].split('.')[1] || '').length };
  }

  function runCounter(el, data, duration) {
    const start = performance.now();
    function tick(now) {
      const t = Math.min((now - start) / duration, 1);
      const eased = 1 - Math.pow(1 - t, 3); // ease-out cubic
      const val   = data.num * eased;
      const fmt   = data.decimals > 0
        ? val.toFixed(data.decimals)
        : Math.round(val).toLocaleString();
      el.textContent = data.prefix + fmt + data.suffix;
      if (t < 1) requestAnimationFrame(tick);
      else       el.textContent = data.original; // restore exact original
    }
    requestAnimationFrame(tick);
  }

  const counterObs = new IntersectionObserver(entries => {
    entries.forEach(entry => {
      if (!entry.isIntersecting) return;
      entry.target.querySelectorAll('.d-number-value').forEach(el => {
        if (el.dataset.dgCounted) return;
        el.dataset.dgCounted = '1';
        const data = parseCounterVal(el.textContent);
        if (data && data.num > 0) runCounter(el, data, 700);
      });
      counterObs.unobserve(entry.target);
    });
  }, { threshold: 0.5 });

  function initCounters(root) {
    root.querySelectorAll('.diagram-container').forEach(c => {
      if (!c.dataset.dgCounter) {
        c.dataset.dgCounter = '1';
        counterObs.observe(c);
      }
    });
  }


  /* ═══════════════════════════════════════════════════════════════
     5. FLOW PULSE ACTIVATION
     The CSS already defines flowDotV / flowDotH animations on
     d-arrow / d-arrow-down but they are paused by default on
     .diagram-container. Activate them once the diagram is visible.
  ══════════════════════════════════════════════════════════════════ */
  const pulseObs = new IntersectionObserver(entries => {
    entries.forEach(entry => {
      const play = entry.isIntersecting ? 'running' : 'paused';
      entry.target.querySelectorAll('.d-arrow,.d-arrow-down').forEach(a => {
        a.style.animationPlayState = play;
        const pseudo = a.querySelector('::before'); // won't work in JS but CSS handles it
        a.dataset.dgPulse = play;
      });
    });
  }, { threshold: 0.1 });

  function initFlowPulse(root) {
    root.querySelectorAll('.diagram-container').forEach(c => {
      if (!c.dataset.dgPulse) {
        c.dataset.dgPulse = '1';
        pulseObs.observe(c);
      }
    });
  }


  /* ═══════════════════════════════════════════════════════════════
     6. LEGEND TOGGLE
     Click d-legend to collapse/expand it. Useful on mobile.
  ══════════════════════════════════════════════════════════════════ */
  function initLegendToggle(root) {
    root.querySelectorAll('.d-legend').forEach(legend => {
      if (legend.dataset.dgLegend) return;
      legend.dataset.dgLegend = '1';
      legend.title = 'Click to collapse/expand legend';
      legend.style.cursor = 'pointer';
      legend.addEventListener('click', () => {
        const items = legend.querySelectorAll('.d-legend-item');
        const hidden = legend.dataset.collapsed === '1';
        items.forEach(item => {
          item.style.transition = 'opacity 0.2s, max-height 0.2s';
          item.style.overflow = 'hidden';
          item.style.maxHeight = hidden ? '40px' : '0px';
          item.style.opacity   = hidden ? ''    : '0';
        });
        legend.dataset.collapsed = hidden ? '' : '1';
      });
    });
  }


  /* ═══════════════════════════════════════════════════════════════
     7. KEYBOARD SHORTCUTS
     Escape → close fullscreen diagram
  ══════════════════════════════════════════════════════════════════ */
  document.addEventListener('keydown', e => {
    if (e.key === 'Escape') {
      document.querySelector('.diagram-container.fullscreen')
        ?.querySelector('.diagram-fullscreen-btn')
        ?.click();
    }
  });


  /* ═══════════════════════════════════════════════════════════════
     INIT + HTMX RE-INIT
  ══════════════════════════════════════════════════════════════════ */
  function initAll(root) {
    initTooltips(root);
    initStepGlow(root);
    initEntrance(root);
    initCounters(root);
    initFlowPulse(root);
    initLegendToggle(root);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => initAll(document.body));
  } else {
    initAll(document.body);
  }

  // Re-run after every HTMX partial swap (new diagram content loaded)
  document.addEventListener('htmx:afterSwap', e => {
    initAll(e.detail.target || document.body);
  });

})();
