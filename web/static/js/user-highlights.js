/**
 * user-highlights.js  v2
 * User highlight + notes system with right-click context menu.
 *
 * Features:
 *  • Right-click text → colour highlight + optional note bubble (visible, dismissible)
 *  • Right-click diagram container → outline highlight + optional note badge
 *  • Left-click image inside diagram → pin annotation at X/Y % position
 *  • Slide-in panel with colour-filter chips, grouped by page → phase
 *  • "Focus mode" – dim phases without highlights; reset restores all
 *  • Dismissible note bubbles (X button); dismissed state persisted in localStorage
 *
 * Storage: localStorage key "sd-user-highlights"
 */
(function () {
    'use strict';

    var STORE_KEY = 'sd-user-highlights';
    var _colorFilter = 'all';   // current panel colour filter
    var _focusMode   = false;
    var _focusDimmed = [];      // elements that got .hl-focus-dim

    // ── Utilities ─────────────────────────────────────────────────────────────

    function esc(s) {
        return String(s || '')
            .replace(/&/g, '&amp;').replace(/</g, '&lt;')
            .replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    }

    function genId() {
        return 'hl_' + Date.now() + '_' + Math.random().toString(36).slice(2, 6);
    }

    // ── Storage ───────────────────────────────────────────────────────────────

    function loadStore() {
        try { return JSON.parse(localStorage.getItem(STORE_KEY) || '{"highlights":[]}'); }
        catch (e) { return { highlights: [] }; }
    }

    function saveStore(store) {
        try { localStorage.setItem(STORE_KEY, JSON.stringify(store)); }
        catch (e) {}
    }

    function storeFind(store, id) {
        for (var i = 0; i < store.highlights.length; i++) {
            if (store.highlights[i].id === id) return store.highlights[i];
        }
        return null;
    }

    // ── Page context ──────────────────────────────────────────────────────────

    function pageCtx() {
        var path = window.location.pathname;
        var h1   = document.querySelector('.detail-header h1, h1.page-title');
        var title = h1
            ? h1.textContent.trim()
            : document.title.replace(' — System Design Prep', '').replace(' | System Design Prep', '').trim();
        var type = path.startsWith('/problem/') ? 'problem'
                 : path.startsWith('/fund/')    ? 'fundamental'
                 : path.startsWith('/algo/')    ? 'algorithm'
                 : path.startsWith('/pattern/') ? 'pattern' : 'page';
        return { path: path, title: title, type: type };
    }

    // ── Nearest phase header ──────────────────────────────────────────────────

    function nearestPhase(node) {
        var detail = document.getElementById('detail');
        var el = (node && node.nodeType === Node.TEXT_NODE) ? node.parentElement : node;
        while (el && el !== detail) {
            var sib = el.previousElementSibling;
            while (sib) {
                if (sib.classList && sib.classList.contains('phase-header')) {
                    var titleEl = sib.querySelector('.phase-title');
                    return { id: sib.id, title: titleEl ? titleEl.textContent.trim() : sib.id };
                }
                sib = sib.previousElementSibling;
            }
            el = el.parentElement;
        }
        return null;
    }

    // ── Range serialisation ───────────────────────────────────────────────────

    function serializeRange(range) {
        var text = range.toString();
        if (!text.trim()) return null;
        var startEl = (range.startContainer.nodeType === Node.TEXT_NODE)
            ? range.startContainer.parentElement : range.startContainer;
        var anchor = startEl;
        while (anchor && !anchor.id && anchor !== document.body) anchor = anchor.parentElement;
        if (!anchor || anchor === document.body || anchor === document.documentElement)
            anchor = document.getElementById('detail');
        var fullText = anchor ? anchor.textContent : '';
        var idx   = fullText.indexOf(text);
        var before = idx >= 0 ? fullText.slice(Math.max(0, idx - 40), idx) : '';
        var after  = idx >= 0 ? fullText.slice(idx + text.length, idx + text.length + 40) : '';
        return { anchorId: anchor ? anchor.id : null, text: text, before: before, after: after };
    }

    function deserializeRange(serial) {
        if (!serial || !serial.text) return null;
        var container = (serial.anchorId && document.getElementById(serial.anchorId))
            || document.getElementById('detail');
        if (!container) return null;
        var fullText = container.textContent;
        var startIdx = -1;
        if (serial.before || serial.after) {
            var searchStr = serial.before + serial.text + serial.after;
            var idx = fullText.indexOf(searchStr);
            if (idx >= 0) startIdx = idx + serial.before.length;
        }
        if (startIdx < 0) startIdx = fullText.indexOf(serial.text);
        if (startIdx < 0) return null;
        return charRangeToDOM(container, startIdx, startIdx + serial.text.length);
    }

    function charRangeToDOM(container, start, end) {
        var walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT);
        var pos = 0, startNode = null, startOff = 0, endNode = null, endOff = 0;
        while (walker.nextNode()) {
            var node = walker.currentNode;
            var len  = node.textContent.length;
            if (!startNode && pos + len > start) { startNode = node; startOff = start - pos; }
            if (startNode && !endNode && pos + len >= end) { endNode = node; endOff = end - pos; break; }
            pos += len;
        }
        if (!startNode || !endNode) return null;
        try {
            var r = document.createRange();
            r.setStart(startNode, startOff);
            r.setEnd(endNode, endOff);
            return r;
        } catch (e) { return null; }
    }

    // ── Note bubble (inline after <mark>, comic speech bubble style) ─────────
    // Expanded: shows note text + × to minimize.
    // Minimized: shrinks to small 💬 icon that re-expands on click.

    function createNoteBubble(hlId, note, minimized) {
        var bubble = document.createElement('span');
        bubble.className    = 'hl-note-bubble';
        bubble.dataset.hlId = hlId;
        if (minimized) bubble.classList.add('hl-nb-minimized');

        var preview = esc(note.slice(0, 100)) + (note.length > 100 ? '\u2026' : '');

        bubble.innerHTML =
            // ── full state ──
            '<span class="hl-nb-full">' +
                '<span class="hl-nb-text">' + preview + '</span>' +
                '<button class="hl-nb-close" title="Minimize note">\u00D7</button>' +
            '</span>' +
            // ── minimized state ──
            '<button class="hl-nb-mini" title="Show note">\uD83D\uDCAC</button>';

        // × → minimize
        bubble.querySelector('.hl-nb-close').addEventListener('click', function (e) {
            e.stopPropagation();
            bubble.classList.add('hl-nb-minimized');
            var store = loadStore();
            var hl = storeFind(store, hlId);
            if (hl) { hl.noteDismissed = true; saveStore(store); }
        });

        // 💬 → re-expand
        bubble.querySelector('.hl-nb-mini').addEventListener('click', function (e) {
            e.stopPropagation();
            bubble.classList.remove('hl-nb-minimized');
            var store = loadStore();
            var hl = storeFind(store, hlId);
            if (hl) { hl.noteDismissed = false; saveStore(store); }
        });

        return bubble;
    }

    // ── Apply / remove text highlight ─────────────────────────────────────────

    function applyTextMark(range, color, hlId, note, dismissed) {
        if (!range || range.collapsed) return false;
        var startEl = (range.startContainer.nodeType === Node.TEXT_NODE)
            ? range.startContainer.parentElement : range.startContainer;
        if (startEl && startEl.closest('mark.user-hl')) return false;

        // Clamp cross-cell table selections to the start cell only
        var startTd = startEl && startEl.closest('td,th');
        if (startTd) {
            var endEl = (range.endContainer.nodeType === Node.TEXT_NODE)
                ? range.endContainer.parentElement : range.endContainer;
            var endTd = endEl && endEl.closest('td,th');
            if (endTd && endTd !== startTd) {
                // Rebuild range within start cell only
                try {
                    var newRange = document.createRange();
                    newRange.setStart(range.startContainer, range.startOffset);
                    newRange.setEnd(startTd, startTd.childNodes.length);
                    range = newRange;
                } catch (e) { /* use original range */ }
            }
        }

        var mark = null;
        try {
            mark = document.createElement('mark');
            mark.className    = 'user-hl user-hl--' + color;
            mark.dataset.hlId = hlId;
            range.surroundContents(mark);
        } catch (e) {
            try {
                var frag = range.extractContents();
                mark = document.createElement('mark');
                mark.className    = 'user-hl user-hl--' + color;
                mark.dataset.hlId = hlId;
                mark.appendChild(frag);
                range.insertNode(mark);
            } catch (e2) { return false; }
        }

        // Bubble is a CHILD of mark so position:relative on mark is the containing block
        if (note && mark) {
            var bubble = createNoteBubble(hlId, note, dismissed);
            mark.appendChild(bubble);
        }
        return true;
    }

    function removeTextMark(hlId) {
        var mark = document.querySelector('mark.user-hl[data-hl-id="' + hlId + '"]');
        if (mark) {
            // Remove bubble child first (so it doesn't get re-inserted during unwrap)
            var bubbleChild = mark.querySelector('.hl-note-bubble');
            if (bubbleChild) mark.removeChild(bubbleChild);
            var p = mark.parentNode;
            while (mark.firstChild) p.insertBefore(mark.firstChild, mark);
            p.removeChild(mark);
            p.normalize();
        }
        // Backwards compat: bubble may be a sibling on older highlights
        var bubble = document.querySelector('.hl-note-bubble[data-hl-id="' + hlId + '"]');
        if (bubble && bubble.parentNode) bubble.parentNode.removeChild(bubble);
    }

    // ── Apply / remove diagram highlight ──────────────────────────────────────

    function applyDiagramMark(diagEl, color, hlId, note, dismissed) {
        diagEl.classList.add('user-hl-diagram', 'user-hl-diagram--' + color);
        diagEl.dataset.hlId = hlId;
        var old = diagEl.querySelector('.user-hl-diagram-note[data-hl-id="' + hlId + '"]');
        if (old && old.parentNode) old.parentNode.removeChild(old);
        if (note) {
            var badge = document.createElement('div');
            badge.className    = 'user-hl-diagram-note';
            badge.dataset.hlId = hlId;
            if (dismissed) badge.classList.add('hl-nb-minimized');
            var preview = esc(note.slice(0, 100)) + (note.length > 100 ? '\u2026' : '');
            badge.innerHTML =
                '<span class="hl-nb-full">' +
                    '<span class="hl-nb-text">\uD83D\uDCAC ' + preview + '</span>' +
                    '<button class="hl-nb-close" title="Minimize">\u00D7</button>' +
                '</span>' +
                '<button class="hl-nb-mini" title="Show note">\uD83D\uDCAC</button>';
            badge.querySelector('.hl-nb-close').addEventListener('click', function (e) {
                e.stopPropagation();
                badge.classList.add('hl-nb-minimized');
                var store = loadStore();
                var hl = storeFind(store, hlId);
                if (hl) { hl.noteDismissed = true; saveStore(store); }
            });
            badge.querySelector('.hl-nb-mini').addEventListener('click', function (e) {
                e.stopPropagation();
                badge.classList.remove('hl-nb-minimized');
                var store = loadStore();
                var hl = storeFind(store, hlId);
                if (hl) { hl.noteDismissed = false; saveStore(store); }
            });
            diagEl.appendChild(badge);
        }
    }

    function removeDiagramMark(hlId) {
        var el = document.querySelector('.user-hl-diagram[data-hl-id="' + hlId + '"]');
        if (el) {
            el.classList.remove(
                'user-hl-diagram',
                'user-hl-diagram--yellow', 'user-hl-diagram--green',
                'user-hl-diagram--blue',   'user-hl-diagram--pink'
            );
            delete el.dataset.hlId;
        }
        var note = document.querySelector('.user-hl-diagram-note[data-hl-id="' + hlId + '"]');
        if (note && note.parentNode) note.parentNode.removeChild(note);
    }

    // ── Apply / remove pin annotation (image X/Y %) ───────────────────────────

    function applyPinMark(diagEl, pctX, pctY, color, hlId, note, dismissed) {
        // Make container position:relative so pins can be placed
        var cs = window.getComputedStyle(diagEl);
        if (cs.position === 'static') diagEl.style.position = 'relative';

        var pin = document.createElement('div');
        pin.className    = 'hl-pin hl-pin--' + color;
        pin.dataset.hlId = hlId;
        pin.style.left   = pctX + '%';
        pin.style.top    = pctY + '%';

        var bubbleHtml = note
            ? '<div class="hl-pin-bubble' + (dismissed ? ' hl-nb-hidden' : '') + '">' +
                '<span class="hl-nb-text">' + esc(note) + '</span>' +
                '<button class="hl-nb-close" data-hl-id="' + esc(hlId) + '" title="Dismiss">\u00D7</button>' +
              '</div>'
            : '';

        pin.innerHTML = '<div class="hl-pin-dot"></div>' + bubbleHtml;

        if (note) {
            pin.querySelector('.hl-nb-close').addEventListener('click', function (e) {
                e.stopPropagation();
                pin.querySelector('.hl-pin-bubble').classList.add('hl-nb-hidden');
                var store = loadStore();
                var hl = storeFind(store, hlId);
                if (hl) { hl.noteDismissed = true; saveStore(store); }
            });
        }

        // Right-click pin → delete context menu
        pin.addEventListener('contextmenu', function (e) {
            e.preventDefault();
            e.stopPropagation();
            showCtxMenu(e.clientX, e.clientY, { type: 'pin', hlId: hlId });
        });

        diagEl.appendChild(pin);
    }

    function removePinMark(hlId) {
        var pin = document.querySelector('.hl-pin[data-hl-id="' + hlId + '"]');
        if (pin && pin.parentNode) pin.parentNode.removeChild(pin);
    }

    // ── Restore page highlights ───────────────────────────────────────────────

    function restorePageHighlights() {
        var store = loadStore();
        var path  = window.location.pathname;
        store.highlights
            .filter(function (h) { return h.page === path; })
            .forEach(function (hl) {
                if (hl.isPin) {
                    var d = document.querySelector('.diagram-container[data-slug="' + hl.diagramSlug + '"]');
                    if (d) applyPinMark(d, hl.pctX, hl.pctY, hl.color, hl.id, hl.note, hl.noteDismissed);
                } else if (hl.isDiagram) {
                    var d2 = document.querySelector('.diagram-container[data-slug="' + hl.diagramSlug + '"]');
                    if (d2 && !d2.dataset.hlId) applyDiagramMark(d2, hl.color, hl.id, hl.note, hl.noteDismissed);
                } else {
                    var range = deserializeRange(hl.range);
                    if (range) applyTextMark(range, hl.color, hl.id, hl.note, hl.noteDismissed);
                }
            });
        updateClearBtn();
        updateFocusBtn();
    }

    // ── Clear page highlights ─────────────────────────────────────────────────

    function clearPageHighlights() {
        var path = window.location.pathname;

        document.querySelectorAll('mark.user-hl').forEach(function (m) {
            // Remove bubble child first so it's not re-inserted during unwrap
            var bc = m.querySelector('.hl-note-bubble');
            if (bc) m.removeChild(bc);
            var p = m.parentNode;
            while (m.firstChild) p.insertBefore(m.firstChild, m);
            p.removeChild(m);
        });
        // Remove any remaining sibling bubbles (backwards compat)
        document.querySelectorAll('.hl-note-bubble').forEach(function (b) {
            b.parentNode && b.parentNode.removeChild(b);
        });
        document.querySelectorAll('.user-hl-diagram').forEach(function (d) {
            d.classList.remove(
                'user-hl-diagram',
                'user-hl-diagram--yellow', 'user-hl-diagram--green',
                'user-hl-diagram--blue',   'user-hl-diagram--pink'
            );
            delete d.dataset.hlId;
        });
        document.querySelectorAll('.user-hl-diagram-note').forEach(function (n) {
            n.parentNode && n.parentNode.removeChild(n);
        });
        document.querySelectorAll('.hl-pin').forEach(function (p) {
            p.parentNode && p.parentNode.removeChild(p);
        });
        var detail = document.getElementById('detail');
        if (detail) detail.normalize();

        var store = loadStore();
        store.highlights = store.highlights.filter(function (h) { return h.page !== path; });
        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
        updateFocusBtn();

        if (_focusMode) exitFocusMode();
    }

    window.clearPageHighlights = clearPageHighlights;

    // ── Focus mode: "show highlighted content only" ───────────────────────────

    function enterFocusMode() {
        if (_focusMode) return;
        _focusMode = true;
        document.body.classList.add('hl-focus-mode');
        updateFocusBtn();
        _focusDimmed = [];

        var detail = document.getElementById('detail');
        if (!detail) return;

        // Build phase→content mapping by walking direct children
        var sections = [];
        var current  = { header: null, content: [] };
        Array.prototype.forEach.call(detail.children, function (el) {
            if (el.classList && el.classList.contains('phase-header')) {
                sections.push(current);
                current = { header: el, content: [] };
            } else {
                current.content.push(el);
            }
        });
        sections.push(current);

        sections.forEach(function (sec) {
            if (!sec.header) return; // pre-phase content, always show
            // Check if any highlight element exists in this section
            var allEls = [sec.header].concat(sec.content);
            var hasHighlight = allEls.some(function (el) {
                return el.querySelector && (
                    el.querySelector('mark.user-hl') ||
                    el.querySelector('.hl-pin') ||
                    el.querySelector('.user-hl-diagram')
                );
            });
            if (!hasHighlight) {
                sec.header.classList.add('hl-focus-dim');
                _focusDimmed.push(sec.header);
                sec.content.forEach(function (el) {
                    el.classList.add('hl-focus-dim');
                    _focusDimmed.push(el);
                });
            } else {
                // Phase has highlights — collapse non-highlighted diagrams inside it
                sec.content.forEach(function (el) {
                    if (!el.classList) return;
                    var diags = el.querySelectorAll
                        ? el.querySelectorAll('.diagram-container:not(.user-hl-diagram)')
                        : [];
                    Array.prototype.forEach.call(diags, function (d) {
                        if (!d.querySelector('.hl-pin')) {
                            d.classList.add('hl-diagram-collapsed');
                            _focusDimmed.push({ el: d, isCollapsed: true });
                            // Wire title click to expand
                            var titleEl = d.querySelector('.diagram-title');
                            if (titleEl && !titleEl._focusClickWired) {
                                titleEl._focusClickWired = true;
                                titleEl.addEventListener('click', function () {
                                    d.classList.remove('hl-diagram-collapsed');
                                });
                            }
                        }
                    });
                });
            }
        });

        showFocusBanner();
    }

    function exitFocusMode() {
        if (!_focusMode) return;
        _focusMode = false;
        document.body.classList.remove('hl-focus-mode');
        updateFocusBtn();
        _focusDimmed.forEach(function (item) {
            if (item && item.isCollapsed) {
                item.el.classList.remove('hl-diagram-collapsed');
            } else if (item && item.classList) {
                item.classList.remove('hl-focus-dim');
            }
        });
        _focusDimmed = [];
        hideFocusBanner();
    }

    window.toggleHighlightFocusMode = function () {
        if (_focusMode) exitFocusMode(); else enterFocusMode();
    };

    function updateFocusBtn() {
        var btn = document.getElementById('hl-focus-btn');
        if (!btn) return;
        // Only show if there are highlights on the current page
        var path = window.location.pathname;
        var has  = loadStore().highlights.some(function (h) { return h.page === path; });
        btn.style.display = has ? '' : 'none';
        btn.classList.toggle('hl-focus-btn--active', _focusMode);
        btn.title = _focusMode ? 'Exit focus mode (show all)' : 'Show only highlighted content';
    }

    function showFocusBanner() {
        var banner = document.getElementById('hl-focus-banner');
        if (!banner) {
            banner = document.createElement('div');
            banner.id        = 'hl-focus-banner';
            banner.className = 'hl-focus-banner';
            banner.innerHTML =
                '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">' +
                    '<circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>' +
                '</svg>' +
                '<span>Focus mode &mdash; showing highlighted sections only.</span>' +
                '<button onclick="window.toggleHighlightFocusMode && window.toggleHighlightFocusMode()">Reset view</button>';
            var detail = document.getElementById('detail');
            if (detail) detail.insertBefore(banner, detail.firstChild);
        }
        banner.style.display = 'flex';
    }

    function hideFocusBanner() {
        var banner = document.getElementById('hl-focus-banner');
        if (banner) banner.style.display = 'none';
    }

    // ── Context Menu ──────────────────────────────────────────────────────────

    var _ctxMenu = null;
    var _pending = null;

    function buildCtxMenu() {
        if (_ctxMenu) return _ctxMenu;
        var m = document.createElement('div');
        m.id        = 'hl-ctx-menu';
        m.className = 'hl-ctx-menu';
        m.innerHTML =
            '<div class="ctx-section ctx-colors-section">' +
                '<span class="ctx-label">Highlight</span>' +
                '<div class="ctx-swatches">' +
                    '<button class="ctx-swatch ctx-swatch--y" data-color="yellow" title="Yellow"></button>' +
                    '<button class="ctx-swatch ctx-swatch--g" data-color="green"  title="Green"></button>' +
                    '<button class="ctx-swatch ctx-swatch--b" data-color="blue"   title="Blue"></button>' +
                    '<button class="ctx-swatch ctx-swatch--p" data-color="pink"   title="Pink"></button>' +
                '</div>' +
            '</div>' +
            '<div class="ctx-section ctx-note-section">' +
                '<button class="ctx-item" id="ctx-note-btn">' +
                    '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">' +
                        '<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>' +
                    '</svg>' +
                    'Highlight + Add Note' +
                '</button>' +
            '</div>' +
            '<div class="ctx-section ctx-delete-section" style="display:none">' +
                '<button class="ctx-item ctx-item--danger" id="ctx-del-btn">' +
                    '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">' +
                        '<polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/>' +
                    '</svg>' +
                    'Remove highlight' +
                '</button>' +
            '</div>';
        document.body.appendChild(m);

        m.querySelectorAll('.ctx-swatch').forEach(function (btn) {
            btn.addEventListener('mousedown', function (e) { e.preventDefault(); });
            btn.addEventListener('click', function () {
                commitHighlight(_pending, btn.dataset.color, '');
                hideCtxMenu();
            });
        });
        document.getElementById('ctx-note-btn').addEventListener('mousedown', function (e) { e.preventDefault(); });
        document.getElementById('ctx-note-btn').addEventListener('click', function () {
            var saved = _pending; hideCtxMenu(); showNoteDialog(saved);
        });
        document.getElementById('ctx-del-btn').addEventListener('click', function () {
            if (_pending && _pending.hlId) deleteHighlight(_pending.hlId);
            hideCtxMenu();
        });

        _ctxMenu = m;
        return m;
    }

    function showCtxMenu(x, y, action) {
        _pending = action;
        var m = buildCtxMenu();
        var hasNew = !!(action.range || (action.type === 'diagram' && !action.hlId));
        m.querySelector('.ctx-colors-section').style.display = hasNew     ? '' : 'none';
        m.querySelector('.ctx-note-section').style.display   = hasNew     ? '' : 'none';
        m.querySelector('.ctx-delete-section').style.display = action.hlId ? '' : 'none';
        m.style.display = 'block';
        requestAnimationFrame(function () {
            m.style.left = Math.min(x, window.innerWidth  - m.offsetWidth  - 8) + 'px';
            m.style.top  = Math.min(y, window.innerHeight - m.offsetHeight - 8) + 'px';
        });
    }

    function hideCtxMenu() {
        if (_ctxMenu) _ctxMenu.style.display = 'none';
        _pending = null;
    }

    // ── Note Dialog (text/diagram highlights) ─────────────────────────────────

    var _noteDlg    = null;
    var _noteAction = null;

    function buildNoteDialog() {
        if (_noteDlg) return _noteDlg;
        var d = document.createElement('div');
        d.id        = 'hl-note-dialog';
        d.className = 'hl-note-dialog';
        d.innerHTML =
            '<div class="hl-note-inner">' +
                '<div class="hl-note-hdr">' +
                    '<span class="hl-note-hdr-title">Highlight + Note</span>' +
                    '<button class="hl-note-close" id="hl-note-close">&#x2715;</button>' +
                '</div>' +
                '<div class="hl-note-colors">' +
                    '<button class="hl-nc-btn hl-nc-btn--yellow active" data-color="yellow"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--green"         data-color="green"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--blue"          data-color="blue"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--pink"          data-color="pink"></button>' +
                '</div>' +
                '<textarea class="hl-note-txt" id="hl-note-txt" placeholder="Add a note (optional)..." rows="3"></textarea>' +
                '<div class="hl-note-footer">' +
                    '<button class="btn btn-sm btn-secondary" id="hl-note-cancel">Cancel</button>' +
                    '<button class="btn btn-sm btn-primary"   id="hl-note-save">Save Highlight</button>' +
                '</div>' +
            '</div>';
        document.body.appendChild(d);

        var selColor = 'yellow';
        d.querySelectorAll('.hl-nc-btn').forEach(function (btn) {
            btn.addEventListener('click', function () {
                d.querySelectorAll('.hl-nc-btn').forEach(function (b) { b.classList.remove('active'); });
                btn.classList.add('active');
                selColor = btn.dataset.color;
            });
        });
        document.getElementById('hl-note-close').addEventListener('click',  closeNoteDialog);
        document.getElementById('hl-note-cancel').addEventListener('click', closeNoteDialog);
        document.getElementById('hl-note-save').addEventListener('click', function () {
            var note = (document.getElementById('hl-note-txt').value || '').trim();
            commitHighlight(_noteAction, selColor, note);
            closeNoteDialog();
        });
        d.addEventListener('click', function (e) { if (e.target === d) closeNoteDialog(); });
        _noteDlg = d;
        return d;
    }

    function showNoteDialog(action) {
        _noteAction = action;
        var d = buildNoteDialog();
        d.querySelectorAll('.hl-nc-btn').forEach(function (b) {
            b.classList.toggle('active', b.dataset.color === 'yellow');
        });
        document.getElementById('hl-note-txt').value = '';
        d.style.display = 'flex';
        setTimeout(function () {
            var t = document.getElementById('hl-note-txt'); if (t) t.focus();
        }, 30);
    }

    function closeNoteDialog() {
        if (_noteDlg) _noteDlg.style.display = 'none';
        _noteAction = null;
    }

    // ── Pin Dialog (image click annotations) ──────────────────────────────────

    var _pinDlg    = null;
    var _pinAction = null;

    function buildPinDialog() {
        if (_pinDlg) return _pinDlg;
        var d = document.createElement('div');
        d.id        = 'hl-pin-dialog';
        d.className = 'hl-pin-dialog';
        d.innerHTML =
            '<div class="hl-pin-inner">' +
                '<div class="hl-pin-hdr">' +
                    '<span class="hl-pin-hdr-title">\uD83D\uDCCD Pin Annotation</span>' +
                    '<button class="hl-note-close" id="hl-pin-close">&#x2715;</button>' +
                '</div>' +
                '<div class="hl-note-colors">' +
                    '<button class="hl-nc-btn hl-nc-btn--yellow active" data-color="yellow"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--green"         data-color="green"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--blue"          data-color="blue"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--pink"          data-color="pink"></button>' +
                '</div>' +
                '<textarea class="hl-note-txt" id="hl-pin-txt" placeholder="Note for this spot (optional)..." rows="2"></textarea>' +
                '<div class="hl-note-footer">' +
                    '<button class="btn btn-sm btn-secondary" id="hl-pin-cancel">Cancel</button>' +
                    '<button class="btn btn-sm btn-primary"   id="hl-pin-save">Add Pin</button>' +
                '</div>' +
            '</div>';
        document.body.appendChild(d);

        var selColor = 'yellow';
        d.querySelectorAll('.hl-nc-btn').forEach(function (btn) {
            btn.addEventListener('click', function () {
                d.querySelectorAll('.hl-nc-btn').forEach(function (b) { b.classList.remove('active'); });
                btn.classList.add('active');
                selColor = btn.dataset.color;
            });
        });
        document.getElementById('hl-pin-close').addEventListener('click',  closePinDialog);
        document.getElementById('hl-pin-cancel').addEventListener('click', closePinDialog);
        document.getElementById('hl-pin-save').addEventListener('click', function () {
            var note = (document.getElementById('hl-pin-txt').value || '').trim();
            commitPinAnnotation(_pinAction, selColor, note);
            closePinDialog();
        });
        d.addEventListener('click', function (e) { if (e.target === d) closePinDialog(); });
        _pinDlg = d;
        return d;
    }

    function showPinDialog(action) {
        _pinAction = action;
        var d = buildPinDialog();
        d.querySelectorAll('.hl-nc-btn').forEach(function (b) {
            b.classList.toggle('active', b.dataset.color === 'yellow');
        });
        document.getElementById('hl-pin-txt').value = '';
        d.style.display = 'flex';
        setTimeout(function () {
            var t = document.getElementById('hl-pin-txt'); if (t) t.focus();
        }, 30);
    }

    function closePinDialog() {
        if (_pinDlg) _pinDlg.style.display = 'none';
        _pinAction = null;
    }

    function commitPinAnnotation(action, color, note) {
        if (!action) return;
        var ctx   = pageCtx();
        var store = loadStore();
        var id    = genId();
        var phase = nearestPhase(action.diagEl);
        var titleEl = action.diagEl.querySelector('.diagram-title');
        var hl = {
            id: id,
            page: ctx.path, pageTitle: ctx.title, pageType: ctx.type,
            color: color, note: note, createdAt: new Date().toISOString(),
            isPin: true, isDiagram: false,
            diagramSlug: action.diagEl.dataset.slug || '',
            pctX: action.pctX, pctY: action.pctY,
            text: '\uD83D\uDCCD ' + (titleEl ? titleEl.textContent.trim() : 'Image annotation'),
            phase: phase
        };
        applyPinMark(action.diagEl, action.pctX, action.pctY, color, id, note, false);
        store.highlights.push(hl);
        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
        updateFocusBtn();
    }

    // ── Commit a highlight ────────────────────────────────────────────────────

    function commitHighlight(action, color, note) {
        if (!action) return;
        var ctx   = pageCtx();
        var store = loadStore();
        var id    = genId();

        if (action.type === 'pin') {
            // Pin right-click = delete only, not commit
            return;
        } else if (action.type === 'diagram') {
            var diagEl = action.diagramEl;
            if (!diagEl || diagEl.dataset.hlId) return;
            var slug    = diagEl.dataset.slug || '';
            var titleEl = diagEl.querySelector('.diagram-title');
            var phase   = nearestPhase(diagEl);
            var hl = {
                id: id,
                page: ctx.path, pageTitle: ctx.title, pageType: ctx.type,
                color: color, note: note, createdAt: new Date().toISOString(),
                isDiagram: true, isPin: false,
                diagramSlug: slug,
                text: (titleEl ? titleEl.textContent.trim() : slug) || 'Diagram',
                phase: phase
            };
            applyDiagramMark(diagEl, color, id, note, false);
            store.highlights.push(hl);
        } else {
            var range = action.range;
            if (!range || range.collapsed) return;
            var text  = range.toString().trim();
            if (!text) return;
            var startEl = (range.startContainer.nodeType === Node.TEXT_NODE)
                ? range.startContainer.parentElement : range.startContainer;
            var phase2  = nearestPhase(startEl);
            var serial  = serializeRange(range);
            if (!serial) return;
            window.getSelection().removeAllRanges();
            var hl2 = {
                id: id,
                page: ctx.path, pageTitle: ctx.title, pageType: ctx.type,
                color: color, note: note, createdAt: new Date().toISOString(),
                isDiagram: false, isPin: false,
                text:  text.slice(0, 200),
                range: serial,
                phase: phase2
            };
            if (applyTextMark(range, color, id, note, false)) store.highlights.push(hl2);
        }

        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
        updateFocusBtn();
    }

    // ── Delete a highlight ────────────────────────────────────────────────────

    function deleteHighlight(hlId) {
        var store = loadStore();
        var hl    = storeFind(store, hlId);
        if (hl) {
            if (hl.isPin)     removePinMark(hlId);
            else if (hl.isDiagram) removeDiagramMark(hlId);
            else              removeTextMark(hlId);
        }
        store.highlights = store.highlights.filter(function (h) { return h.id !== hlId; });
        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
        updateFocusBtn();
        var panel = document.getElementById('hl-panel');
        if (panel && panel.classList.contains('open')) renderPanel();
    }

    // ── Header button helpers ─────────────────────────────────────────────────

    function updateClearBtn() {
        var btn  = document.getElementById('hl-clear-btn');
        if (!btn) return;
        var has  = loadStore().highlights.some(function (h) { return h.page === window.location.pathname; });
        btn.style.display = has ? '' : 'none';
    }

    function updatePanelBadge() {
        var badge = document.getElementById('hl-panel-badge');
        if (!badge) return;
        var count = loadStore().highlights.length;
        badge.textContent  = count > 0 ? String(count) : '';
        badge.style.display = count > 0 ? '' : 'none';
    }

    // ── Highlights Panel ──────────────────────────────────────────────────────

    var _dragSelecting = false;
    var _selected      = {};

    function openPanel() {
        var panel = buildPanel();
        renderPanel();
        panel.classList.add('open');
    }
    window.openHighlightsPanel = openPanel;

    function closePanel() {
        var panel = document.getElementById('hl-panel');
        if (panel) panel.classList.remove('open');
        _selected = {};
    }
    window.closeHighlightsPanel = closePanel;

    function buildPanel() {
        var existing = document.getElementById('hl-panel');
        if (existing) return existing;

        var panel = document.createElement('div');
        panel.id        = 'hl-panel';
        panel.className = 'hl-panel';
        panel.innerHTML =
            '<div class="hl-panel-hdr">' +
                '<span class="hl-panel-title">' +
                    '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">' +
                        '<path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>' +
                    '</svg>' +
                    ' My Highlights' +
                '</span>' +
                '<div class="hl-panel-hdr-right">' +
                    '<button class="hl-panel-btn hl-del-sel-btn" id="hl-del-sel" style="display:none">Delete selected</button>' +
                    '<button class="hl-panel-btn hl-clear-all-btn" id="hl-clear-all">Clear all</button>' +
                    '<button class="hl-panel-close" id="hl-panel-close">&#x2715;</button>' +
                '</div>' +
            '</div>' +
            '<div class="hl-panel-filters" id="hl-panel-filters">' +
                '<button class="hl-filter-chip hl-filter-chip--all active" data-filter="all">All</button>' +
                '<button class="hl-filter-chip" data-filter="yellow"><span class="hl-fc-dot hl-fc-dot--yellow"></span>Yellow</button>' +
                '<button class="hl-filter-chip" data-filter="green"><span class="hl-fc-dot hl-fc-dot--green"></span>Green</button>' +
                '<button class="hl-filter-chip" data-filter="blue"><span class="hl-fc-dot hl-fc-dot--blue"></span>Blue</button>' +
                '<button class="hl-filter-chip" data-filter="pink"><span class="hl-fc-dot hl-fc-dot--pink"></span>Pink</button>' +
            '</div>' +
            '<div class="hl-panel-body" id="hl-panel-body"></div>';
        document.body.appendChild(panel);

        document.getElementById('hl-panel-close').addEventListener('click', closePanel);

        document.getElementById('hl-clear-all').addEventListener('click', function () {
            if (!confirm('Delete ALL highlights across all pages?')) return;
            saveStore({ highlights: [] });
            clearPageHighlights();
            updatePanelBadge();
            renderPanel();
        });

        document.getElementById('hl-del-sel').addEventListener('click', function () {
            Object.keys(_selected).forEach(function (hlId) { deleteHighlight(hlId); });
            _selected = {};
            updateDelSelBtn();
        });

        // Colour filter chips
        panel.querySelectorAll('.hl-filter-chip').forEach(function (chip) {
            chip.addEventListener('click', function () {
                panel.querySelectorAll('.hl-filter-chip').forEach(function (c) { c.classList.remove('active'); });
                chip.classList.add('active');
                _colorFilter = chip.dataset.filter;
                renderPanel();
            });
        });

        // Close on outside mousedown
        document.addEventListener('mousedown', function (e) {
            var p = document.getElementById('hl-panel');
            if (p && p.classList.contains('open') &&
                !p.contains(e.target) &&
                !e.target.closest('#hl-open-btn')) {
                closePanel();
            }
        });

        return panel;
    }

    function renderPanel() {
        var body = document.getElementById('hl-panel-body');
        if (!body) return;
        _selected = {};
        updateDelSelBtn();

        var store = loadStore();
        var hls   = store.highlights;

        // Apply colour filter
        if (_colorFilter !== 'all') {
            hls = hls.filter(function (h) { return h.color === _colorFilter; });
        }

        if (!hls.length) {
            body.innerHTML =
                '<p class="hl-panel-empty">' +
                (_colorFilter !== 'all'
                    ? 'No ' + _colorFilter + ' highlights.'
                    : 'No highlights yet.<br>Right-click any text or diagram to highlight it.') +
                '</p>';
            return;
        }

        // Group by page
        var pageOrder = [];
        var pages     = {};
        hls.forEach(function (h) {
            if (!pages[h.page]) {
                pages[h.page] = { title: h.pageTitle || h.page, type: h.pageType || 'page', items: [] };
                pageOrder.push(h.page);
            }
            pages[h.page].items.push(h);
        });

        var html = '';
        pageOrder.forEach(function (pg) {
            var g = pages[pg];
            g.items.sort(function (a, b) { return a.createdAt < b.createdAt ? -1 : 1; });

            html += '<div class="hl-pg-group">';
            html +=
                '<div class="hl-pg-hdr">' +
                    '<span class="hl-pg-type hl-pg-type--' + esc(g.type) + '">' + esc(g.type) + '</span>' +
                    '<a class="hl-pg-title" href="' + esc(pg) + '"' +
                        ' hx-get="' + esc(pg) + '" hx-target="#detail" hx-swap="innerHTML" hx-push-url="true">' +
                        esc(g.title) +
                    '</a>' +
                    '<span class="hl-pg-count">' + g.items.length + '</span>' +
                '</div>';

            g.items.forEach(function (hl) {
                var phaseId  = hl.phase ? hl.phase.id : '';
                var gotoHref = pg + (phaseId ? '#' + phaseId : '');
                var snippet  = (hl.text || '').slice(0, 100);
                var ellipsis = (hl.text && hl.text.length > 100) ? '\u2026' : '';
                var dotClass = hl.isPin
                    ? 'hl-pin-dot-sm'
                    : 'hl-item-dot hl-item-dot--' + esc(hl.color);

                html +=
                    '<div class="hl-item" data-hl-id="' + esc(hl.id) + '">' +
                        '<input type="checkbox" class="hl-item-cb" data-hl-id="' + esc(hl.id) + '">' +
                        '<div class="' + dotClass + '"></div>' +
                        '<div class="hl-item-body">' +
                            (hl.phase ? '<div class="hl-item-phase">' + esc(hl.phase.title) + '</div>' : '') +
                            '<div class="hl-item-text">' + esc(snippet + ellipsis) + '</div>' +
                            (hl.note ? '<div class="hl-item-note">\uD83D\uDCAC ' + esc(hl.note) + '</div>' : '') +
                            '<div class="hl-item-actions">' +
                                '<a class="hl-goto" href="' + esc(gotoHref) + '"' +
                                    ' hx-get="' + esc(pg) + '" hx-target="#detail"' +
                                    ' hx-swap="innerHTML" hx-push-url="true"' +
                                    ' data-phase-id="' + esc(phaseId) + '">Go to \u2192</a>' +
                                '<button class="hl-del-btn" data-hl-id="' + esc(hl.id) + '">Delete</button>' +
                            '</div>' +
                        '</div>' +
                    '</div>';
            });
            html += '</div>';
        });

        body.innerHTML = html;

        // Wire checkboxes and drag-select
        body.querySelectorAll('.hl-item').forEach(function (item) {
            var cb = item.querySelector('.hl-item-cb');
            if (!cb) return;
            cb.addEventListener('change', function () {
                if (cb.checked) _selected[cb.dataset.hlId] = true;
                else            delete _selected[cb.dataset.hlId];
                item.classList.toggle('hl-item--selected', cb.checked);
                updateDelSelBtn();
            });
            item.addEventListener('click', function (e) {
                if (e.target === cb) return;
                if (e.target.tagName === 'A' || e.target.tagName === 'BUTTON') return;
                cb.checked = !cb.checked;
                cb.dispatchEvent(new Event('change'));
            });
        });

        body.addEventListener('mousedown', function (e) {
            if (e.target.tagName === 'A' || e.target.tagName === 'BUTTON' || e.target.tagName === 'INPUT') return;
            if (e.target.closest('.hl-item')) _dragSelecting = true;
        });

        // "Go to" links — scroll to phase after navigation
        body.querySelectorAll('.hl-goto').forEach(function (a) {
            a.addEventListener('click', function () {
                var phaseId = a.dataset.phaseId;
                closePanel();
                if (phaseId) {
                    setTimeout(function () {
                        var ph = document.getElementById(phaseId);
                        if (ph) ph.scrollIntoView({ behavior: 'smooth', block: 'start' });
                    }, 450);
                }
            });
        });

        // Delete buttons
        body.querySelectorAll('.hl-del-btn').forEach(function (btn) {
            btn.addEventListener('click', function (e) {
                e.stopPropagation();
                deleteHighlight(btn.dataset.hlId);
            });
        });

        // Process HTMX on injected links
        if (typeof htmx !== 'undefined') {
            body.querySelectorAll('[hx-get]').forEach(function (el) { htmx.process(el); });
        }
    }

    function updateDelSelBtn() {
        var btn = document.getElementById('hl-del-sel');
        if (btn) btn.style.display = Object.keys(_selected).length > 0 ? '' : 'none';
    }

    // Global drag-select mousemove / mouseup
    document.addEventListener('mousemove', function (e) {
        if (!_dragSelecting) return;
        var item = e.target && e.target.closest && e.target.closest('.hl-item');
        if (!item) return;
        var cb = item.querySelector('.hl-item-cb');
        if (!cb || cb.checked) return;
        cb.checked = true;
        _selected[cb.dataset.hlId] = true;
        item.classList.add('hl-item--selected');
        updateDelSelBtn();
    });
    document.addEventListener('mouseup', function () { _dragSelecting = false; });

    // ── Context menu + image click wiring ────────────────────────────────────

    var _ctxAttached  = false;
    var _imgMousedown = null; // track mousedown pos to distinguish click vs drag

    function attachContextMenu() {
        if (_ctxAttached) return;
        _ctxAttached = true;

        // Image click → pin annotation (only if mouse didn't move significantly)
        document.addEventListener('mousedown', function (e) {
            var img = e.target.closest('.diagram-container img');
            if (img) _imgMousedown = { x: e.clientX, y: e.clientY };
            else     _imgMousedown = null;
        });

        document.addEventListener('click', function (e) {
            if (!_imgMousedown) return;
            var img = e.target.closest('.diagram-container img');
            if (!img) { _imgMousedown = null; return; }
            // Ignore if mouse moved too much (panning)
            var dx = Math.abs(e.clientX - _imgMousedown.x);
            var dy = Math.abs(e.clientY - _imgMousedown.y);
            _imgMousedown = null;
            if (dx > 5 || dy > 5) return;
            // Ignore if text is selected
            var sel = window.getSelection();
            if (sel && !sel.isCollapsed && sel.toString().trim()) return;

            var diagEl = img.closest('.diagram-container');
            if (!diagEl) return;
            var rect = img.getBoundingClientRect();
            var pctX = ((e.clientX - rect.left)  / rect.width  * 100).toFixed(1);
            var pctY = ((e.clientY - rect.top)   / rect.height * 100).toFixed(1);
            showPinDialog({ diagEl: diagEl, pctX: parseFloat(pctX), pctY: parseFloat(pctY) });
        });

        // Right-click context menu (text selection, existing marks, diagrams)
        document.addEventListener('contextmenu', function (e) {
            var detail = document.getElementById('detail');
            if (!detail || !detail.contains(e.target)) return;
            if (['INPUT', 'TEXTAREA', 'SELECT'].indexOf(e.target.tagName) >= 0) return;

            var sel       = window.getSelection();
            var hasText   = sel && !sel.isCollapsed && sel.toString().trim().length > 0;
            var existMark = e.target.closest('mark.user-hl');
            var existPin  = e.target.closest('.hl-pin');
            var diagEl    = e.target.closest('.diagram-container');

            if (!hasText && !existMark && !existPin && !diagEl) return;
            e.preventDefault();

            var action;
            if (hasText) {
                action = {
                    type:  'text',
                    range: sel.getRangeAt(0).cloneRange(),
                    hlId:  existMark ? existMark.dataset.hlId : null
                };
            } else if (existPin) {
                action = { type: 'pin', hlId: existPin.dataset.hlId };
            } else if (existMark) {
                action = { type: 'text', hlId: existMark.dataset.hlId };
            } else if (diagEl) {
                action = {
                    type:      'diagram',
                    diagramEl: diagEl,
                    hlId:      diagEl.dataset.hlId || null
                };
            }
            if (action) showCtxMenu(e.clientX, e.clientY, action);
        });

        // Close context menu on outside click
        document.addEventListener('click', function (e) {
            if (_ctxMenu && !e.target.closest('#hl-ctx-menu')) hideCtxMenu();
        });

        // Escape closes everything
        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape') {
                hideCtxMenu();
                closeNoteDialog();
                closePinDialog();
            }
        });
    }

    // ── Initialise ────────────────────────────────────────────────────────────

    function init() {
        attachContextMenu();
        updateClearBtn();
        updatePanelBadge();
        updateFocusBtn();
        requestAnimationFrame(restorePageHighlights);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // Re-run after every HTMX swap into #detail
    document.addEventListener('htmx:afterSwap', function (e) {
        if (e.detail && e.detail.target && e.detail.target.id === 'detail') {
            if (_focusMode) exitFocusMode();
            updateClearBtn();
            updatePanelBadge();
            updateFocusBtn();
            setTimeout(restorePageHighlights, 150);
        }
    });

}());
