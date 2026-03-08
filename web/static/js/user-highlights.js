/**
 * user-highlights.js
 * User highlight + notes system with right-click context menu.
 * No mouseup toolbar — context menu only (uniform for text and diagram/SVG).
 *
 * Storage: localStorage key "sd-user-highlights"
 * Highlights panel: slides in from right
 * Clear button: visible in header when current page has highlights
 */
(function () {
    'use strict';

    var STORE_KEY = 'sd-user-highlights';

    // ── Utilities ─────────────────────────────────────────────────────────────

    function esc(s) {
        return String(s || '')
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;');
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

    // ── Page context ──────────────────────────────────────────────────────────

    function pageCtx() {
        var path = window.location.pathname;
        var h1 = document.querySelector('.detail-header h1, h1.page-title');
        var title = h1
            ? h1.textContent.trim()
            : document.title.replace(' — System Design Prep', '').replace(' | System Design Prep', '').trim();
        var type = path.startsWith('/problem/') ? 'problem'
                 : path.startsWith('/fund/')    ? 'fundamental'
                 : path.startsWith('/algo/')    ? 'algorithm'
                 : path.startsWith('/pattern/') ? 'pattern' : 'page';
        return { path: path, title: title, type: type };
    }

    // ── Find nearest phase header above a node ────────────────────────────────
    // Walks up the DOM, checking previous siblings at each level for .phase-header.

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
    // Stores: anchorId (nearest ancestor with ID), selected text, and 40-char
    // context before/after for reliable restoration via text search.

    function serializeRange(range) {
        var text = range.toString();
        if (!text.trim()) return null;

        // Find nearest ancestor element that has an ID
        var startEl = (range.startContainer.nodeType === Node.TEXT_NODE)
            ? range.startContainer.parentElement
            : range.startContainer;
        var anchor = startEl;
        while (anchor && !anchor.id && anchor !== document.body) {
            anchor = anchor.parentElement;
        }
        if (!anchor || anchor === document.body || anchor === document.documentElement) {
            anchor = document.getElementById('detail');
        }

        // Capture surrounding context from anchor's full text
        var fullText = anchor ? anchor.textContent : '';
        var idx = fullText.indexOf(text);
        var before = (idx >= 0) ? fullText.slice(Math.max(0, idx - 40), idx) : '';
        var after  = (idx >= 0) ? fullText.slice(idx + text.length, idx + text.length + 40) : '';

        return {
            anchorId: anchor ? anchor.id : null,
            text:     text,
            before:   before,
            after:    after
        };
    }

    // Deserialise: find the text in the DOM and return a live Range.
    function deserializeRange(serial) {
        if (!serial || !serial.text) return null;

        var container = (serial.anchorId && document.getElementById(serial.anchorId))
            || document.getElementById('detail');
        if (!container) return null;

        var fullText = container.textContent;
        var target   = serial.text;

        // Try with context first for uniqueness, then fallback to text alone
        var startIdx = -1;
        if (serial.before || serial.after) {
            var searchStr = serial.before + target + serial.after;
            var idx = fullText.indexOf(searchStr);
            if (idx >= 0) startIdx = idx + serial.before.length;
        }
        if (startIdx < 0) {
            startIdx = fullText.indexOf(target);
        }
        if (startIdx < 0) return null;

        return charRangeToDOM(container, startIdx, startIdx + target.length);
    }

    // Convert start/end char offsets (within container.textContent) to a DOM Range.
    function charRangeToDOM(container, start, end) {
        var walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT);
        var pos = 0;
        var startNode = null, startOff = 0, endNode = null, endOff = 0;

        while (walker.nextNode()) {
            var node = walker.currentNode;
            var len  = node.textContent.length;

            if (!startNode && pos + len > start) {
                startNode = node;
                startOff  = start - pos;
            }
            if (startNode && !endNode && pos + len >= end) {
                endNode = node;
                endOff  = end - pos;
                break;
            }
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

    // ── Apply / remove text highlight in DOM ──────────────────────────────────

    function applyTextMark(range, color, hlId, note) {
        if (!range || range.collapsed) return false;
        // Don't highlight inside an existing user mark
        var startEl = (range.startContainer.nodeType === Node.TEXT_NODE)
            ? range.startContainer.parentElement : range.startContainer;
        if (startEl && startEl.closest('mark.user-hl')) return false;

        var cls  = 'user-hl user-hl--' + color;
        var mark = null;
        try {
            mark = document.createElement('mark');
            mark.className    = cls;
            mark.dataset.hlId = hlId;
            range.surroundContents(mark);
        } catch (e) {
            // Cross-element range: extract → wrap → reinsert
            try {
                var frag = range.extractContents();
                mark = document.createElement('mark');
                mark.className    = cls;
                mark.dataset.hlId = hlId;
                mark.appendChild(frag);
                range.insertNode(mark);
            } catch (e2) { return false; }
        }

        if (note && mark) {
            var icon = document.createElement('span');
            icon.className    = 'user-hl-note-icon';
            icon.textContent  = '\uD83D\uDCAC'; // 💬
            icon.title        = note;
            icon.dataset.hlId = hlId;
            if (mark.nextSibling) {
                mark.parentNode.insertBefore(icon, mark.nextSibling);
            } else {
                mark.parentNode.appendChild(icon);
            }
        }
        return true;
    }

    function removeTextMark(hlId) {
        var mark = document.querySelector('mark.user-hl[data-hl-id="' + hlId + '"]');
        if (mark) {
            var p = mark.parentNode;
            while (mark.firstChild) p.insertBefore(mark.firstChild, mark);
            p.removeChild(mark);
            p.normalize();
        }
        var icon = document.querySelector('.user-hl-note-icon[data-hl-id="' + hlId + '"]');
        if (icon) icon.parentNode && icon.parentNode.removeChild(icon);
    }

    // ── Apply / remove diagram highlight in DOM ───────────────────────────────

    function applyDiagramMark(diagEl, color, hlId, note) {
        diagEl.classList.add('user-hl-diagram', 'user-hl-diagram--' + color);
        diagEl.dataset.hlId = hlId;
        if (note) {
            var badge = document.createElement('div');
            badge.className    = 'user-hl-diagram-note';
            badge.dataset.hlId = hlId;
            badge.innerHTML    = '\uD83D\uDCAC ' + esc(note.slice(0, 80)) + (note.length > 80 ? '\u2026' : '');
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

    // ── Restore highlights for current page ───────────────────────────────────

    function restorePageHighlights() {
        var store = loadStore();
        var path  = window.location.pathname;
        store.highlights
            .filter(function (h) { return h.page === path; })
            .forEach(function (hl) {
                if (hl.isDiagram) {
                    var diagEl = document.querySelector(
                        '.diagram-container[data-slug="' + hl.diagramSlug + '"]'
                    );
                    if (diagEl && !diagEl.dataset.hlId) {
                        applyDiagramMark(diagEl, hl.color, hl.id, hl.note);
                    }
                } else {
                    var range = deserializeRange(hl.range);
                    if (range) applyTextMark(range, hl.color, hl.id, hl.note);
                }
            });
        updateClearBtn();
    }

    // ── Clear page highlights ─────────────────────────────────────────────────

    function clearPageHighlights() {
        var path = window.location.pathname;
        // Remove mark elements (unwrap)
        document.querySelectorAll('mark.user-hl').forEach(function (m) {
            var p = m.parentNode;
            while (m.firstChild) p.insertBefore(m.firstChild, m);
            p.removeChild(m);
        });
        // Remove note icons
        document.querySelectorAll('.user-hl-note-icon').forEach(function (i) {
            i.parentNode && i.parentNode.removeChild(i);
        });
        // Remove diagram highlights
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
        // Merge split text nodes
        var detail = document.getElementById('detail');
        if (detail) detail.normalize();

        // Purge from store
        var store = loadStore();
        store.highlights = store.highlights.filter(function (h) { return h.page !== path; });
        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
    }

    window.clearPageHighlights = clearPageHighlights;

    // ── Context Menu ──────────────────────────────────────────────────────────

    var _ctxMenu = null;
    var _pending = null; // { type:'text'|'diagram', range?, diagramEl?, hlId? }

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
                        '<polyline points="3 6 5 6 21 6"/>' +
                        '<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/>' +
                    '</svg>' +
                    'Remove highlight' +
                '</button>' +
            '</div>';
        document.body.appendChild(m);

        // Colour swatches — use mousedown+preventDefault to preserve selection
        m.querySelectorAll('.ctx-swatch').forEach(function (btn) {
            btn.addEventListener('mousedown', function (e) { e.preventDefault(); });
            btn.addEventListener('click', function () {
                commitHighlight(_pending, btn.dataset.color, '');
                hideCtxMenu();
            });
        });

        document.getElementById('ctx-note-btn').addEventListener('mousedown', function (e) { e.preventDefault(); });
        document.getElementById('ctx-note-btn').addEventListener('click', function () {
            var saved = _pending;
            hideCtxMenu();
            showNoteDialog(saved);
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

        // Toggle sections based on context
        var hasNewSelection = !!(action.range || (action.type === 'diagram' && !action.hlId));
        m.querySelector('.ctx-colors-section').style.display = hasNewSelection ? '' : 'none';
        m.querySelector('.ctx-note-section').style.display   = hasNewSelection ? '' : 'none';
        m.querySelector('.ctx-delete-section').style.display = action.hlId    ? '' : 'none';

        m.style.display = 'block';
        // Clamp to viewport after layout
        requestAnimationFrame(function () {
            var mw = m.offsetWidth, mh = m.offsetHeight;
            m.style.left = Math.min(x, window.innerWidth  - mw - 8) + 'px';
            m.style.top  = Math.min(y, window.innerHeight - mh - 8) + 'px';
        });
    }

    function hideCtxMenu() {
        if (_ctxMenu) _ctxMenu.style.display = 'none';
        _pending = null;
    }

    // ── Note Dialog ───────────────────────────────────────────────────────────

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
                    '<button class="hl-nc-btn hl-nc-btn--yellow active" data-color="yellow" title="Yellow"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--green"         data-color="green"  title="Green"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--blue"          data-color="blue"   title="Blue"></button>' +
                    '<button class="hl-nc-btn hl-nc-btn--pink"          data-color="pink"   title="Pink"></button>' +
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
        // Backdrop click closes
        d.addEventListener('click', function (e) { if (e.target === d) closeNoteDialog(); });

        _noteDlg = d;
        return d;
    }

    function showNoteDialog(action) {
        _noteAction = action;
        var d = buildNoteDialog();
        // Reset colour selection
        d.querySelectorAll('.hl-nc-btn').forEach(function (b) {
            b.classList.toggle('active', b.dataset.color === 'yellow');
        });
        document.getElementById('hl-note-txt').value = '';
        d.style.display = 'flex';
        setTimeout(function () {
            var txt = document.getElementById('hl-note-txt');
            if (txt) txt.focus();
        }, 30);
    }

    function closeNoteDialog() {
        if (_noteDlg) _noteDlg.style.display = 'none';
        _noteAction = null;
    }

    // ── Commit a highlight ────────────────────────────────────────────────────

    function commitHighlight(action, color, note) {
        if (!action) return;
        var ctx   = pageCtx();
        var store = loadStore();
        var id    = genId();

        if (action.type === 'diagram') {
            var diagEl = action.diagramEl;
            if (!diagEl) return;
            // Skip if already highlighted
            if (diagEl.dataset.hlId) return;
            var slug    = diagEl.dataset.slug || '';
            var titleEl = diagEl.querySelector('.diagram-title');
            var phase   = nearestPhase(diagEl);
            var hl = {
                id: id,
                page: ctx.path, pageTitle: ctx.title, pageType: ctx.type,
                color: color, note: note, createdAt: new Date().toISOString(),
                isDiagram: true,
                diagramSlug: slug,
                text: (titleEl ? titleEl.textContent.trim() : slug) || 'Diagram',
                phase: phase
            };
            applyDiagramMark(diagEl, color, id, note);
            store.highlights.push(hl);

        } else {
            var range = action.range;
            if (!range || range.collapsed) return;
            var text = range.toString().trim();
            if (!text) return;

            var startEl = (range.startContainer.nodeType === Node.TEXT_NODE)
                ? range.startContainer.parentElement
                : range.startContainer;
            var phase2  = nearestPhase(startEl);
            var serial  = serializeRange(range);
            if (!serial) return;

            // Clear browser selection before DOM modification
            window.getSelection().removeAllRanges();

            var hl2 = {
                id: id,
                page: ctx.path, pageTitle: ctx.title, pageType: ctx.type,
                color: color, note: note, createdAt: new Date().toISOString(),
                isDiagram: false,
                text:  text.slice(0, 200),
                range: serial,
                phase: phase2
            };
            if (applyTextMark(range, color, id, note)) {
                store.highlights.push(hl2);
            }
        }

        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
    }

    // ── Delete a single highlight ─────────────────────────────────────────────

    function deleteHighlight(hlId) {
        var store = loadStore();
        var hl    = store.highlights.filter(function (h) { return h.id === hlId; })[0];
        if (hl && hl.isDiagram) removeDiagramMark(hlId);
        else                    removeTextMark(hlId);
        store.highlights = store.highlights.filter(function (h) { return h.id !== hlId; });
        saveStore(store);
        updateClearBtn();
        updatePanelBadge();
        // Re-render panel if it's open
        var panel = document.getElementById('hl-panel');
        if (panel && panel.classList.contains('open')) renderPanel();
    }

    // ── Header button helpers ─────────────────────────────────────────────────

    function updateClearBtn() {
        var btn  = document.getElementById('hl-clear-btn');
        if (!btn) return;
        var path = window.location.pathname;
        var has  = loadStore().highlights.some(function (h) { return h.page === path; });
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
    var _selected      = {};   // hlId → true

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
                        '<path d="M12 20h9"/>' +
                        '<path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>' +
                    '</svg>' +
                    ' My Highlights' +
                '</span>' +
                '<div class="hl-panel-hdr-right">' +
                    '<button class="hl-panel-btn hl-del-sel-btn" id="hl-del-sel" style="display:none">Delete selected</button>' +
                    '<button class="hl-panel-btn hl-clear-all-btn" id="hl-clear-all">Clear all</button>' +
                    '<button class="hl-panel-close" id="hl-panel-close">&#x2715;</button>' +
                '</div>' +
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

        if (!hls.length) {
            body.innerHTML =
                '<p class="hl-panel-empty">' +
                'No highlights yet.<br>' +
                'Right-click any text or diagram to highlight it.' +
                '</p>';
            return;
        }

        // Group by page, preserve insertion order
        var pageOrder = [];
        var pages = {};
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

                html +=
                    '<div class="hl-item" data-hl-id="' + esc(hl.id) + '">' +
                        '<input type="checkbox" class="hl-item-cb" data-hl-id="' + esc(hl.id) + '">' +
                        '<div class="hl-item-dot hl-item-dot--' + esc(hl.color) + '"></div>' +
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

        // ── Wire checkboxes for multi-select (+ drag-to-select) ──────────────
        var items = body.querySelectorAll('.hl-item');
        items.forEach(function (item) {
            var cb = item.querySelector('.hl-item-cb');
            if (!cb) return;

            cb.addEventListener('change', function () {
                if (cb.checked) _selected[cb.dataset.hlId] = true;
                else            delete _selected[cb.dataset.hlId];
                item.classList.toggle('hl-item--selected', cb.checked);
                updateDelSelBtn();
            });

            // Click anywhere on item (except links/buttons) toggles checkbox
            item.addEventListener('click', function (e) {
                if (e.target === cb) return;
                if (e.target.tagName === 'A' || e.target.tagName === 'BUTTON') return;
                cb.checked = !cb.checked;
                cb.dispatchEvent(new Event('change'));
            });
        });

        // Drag-to-select: hold mousedown and drag over items
        var dragStart = null;
        body.addEventListener('mousedown', function (e) {
            if (e.target.tagName === 'A' || e.target.tagName === 'BUTTON' || e.target.tagName === 'INPUT') return;
            dragStart = e.target.closest('.hl-item');
            if (dragStart) _dragSelecting = true;
        });
        // mousemove wired globally (see bottom of file)

        // ── Wire "Go to" links ───────────────────────────────────────────────
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

        // ── Wire delete buttons ──────────────────────────────────────────────
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

    // Global mousemove / mouseup for drag-select in panel
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

    // ── Context menu event wiring (attached once, globally) ──────────────────

    var _ctxAttached = false;

    function attachContextMenu() {
        if (_ctxAttached) return;
        _ctxAttached = true;

        document.addEventListener('contextmenu', function (e) {
            var detail = document.getElementById('detail');
            if (!detail || !detail.contains(e.target)) return;
            // Pass through on form fields
            if (['INPUT', 'TEXTAREA', 'SELECT'].indexOf(e.target.tagName) >= 0) return;

            var sel         = window.getSelection();
            var hasText     = sel && !sel.isCollapsed && sel.toString().trim().length > 0;
            var existMark   = e.target.closest('mark.user-hl');
            var diagEl      = e.target.closest('.diagram-container');

            if (!hasText && !existMark && !diagEl) return; // no highlight target → native menu

            e.preventDefault();

            var action;
            if (hasText) {
                // New text selection (may overlap an existing mark)
                action = {
                    type:  'text',
                    range: sel.getRangeAt(0).cloneRange(),
                    hlId:  existMark ? existMark.dataset.hlId : null
                };
            } else if (existMark) {
                // Right-click on existing text highlight → delete only
                action = { type: 'text', hlId: existMark.dataset.hlId };
            } else if (diagEl) {
                // Diagram: highlight or delete existing
                action = {
                    type:      'diagram',
                    diagramEl: diagEl,
                    hlId:      diagEl.dataset.hlId || null
                };
            }
            if (action) showCtxMenu(e.clientX, e.clientY, action);
        });

        // Close menu on any click outside
        document.addEventListener('click', function (e) {
            if (_ctxMenu && !e.target.closest('#hl-ctx-menu')) hideCtxMenu();
        });

        // Close on Escape
        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape') { hideCtxMenu(); closeNoteDialog(); }
        });
    }

    // ── Initialise ───────────────────────────────────────────────────────────

    function init() {
        attachContextMenu();
        updateClearBtn();
        updatePanelBadge();
        // Restore after next paint so content is fully laid out
        requestAnimationFrame(restorePageHighlights);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // Re-run after every HTMX page swap
    document.addEventListener('htmx:afterSwap', function (e) {
        if (e.detail && e.detail.target && e.detail.target.id === 'detail') {
            updateClearBtn();
            updatePanelBadge();
            // Small delay lets HTMX fully settle the new DOM before we scan it
            setTimeout(restorePageHighlights, 150);
        }
    });

}());
