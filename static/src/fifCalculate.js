// fifCalculate.js
(function () {
    if (window.__fifInit) return;
    window.__fifInit = true;

    // ---------- Global state ----------
    window.fifState = window.fifState || {};          // input snapshot by id
    window.viewedHoldings = window.viewedHoldings || {}; // VERIFIED symbols only
    window.requiredSymbols = window.requiredSymbols || new Set();

    // ---------- Helpers ----------
    function inForm(target) {
        const form = document.getElementById('fif-form');
        return form && form.contains(target);
    }

    window.allHoldingsVisited = function () {
        for (const sym of window.requiredSymbols) {
            if (!window.viewedHoldings[sym]) return false;
        }
        return window.requiredSymbols.size > 0;
    };

    // Update segment buttons' look (verified => green + ✔)
    window.refreshHoldingButtons = function () {
        const buttons = document.querySelectorAll('#holdings-segment > button[data-symbol]');
        buttons.forEach(btn => {
            const sym = btn.dataset.symbol;
            const verified = !!window.viewedHoldings[sym];
            const tick = btn.querySelector('.tick');

            if (verified) {
                btn.classList.remove(
                    'bg-white','hover:bg-gray-50','text-gray-800',
                    'dark:bg-gray-800','dark:hover:bg-gray-700','dark:text-gray-100'
                );
                btn.classList.add(
                    'bg-green-100','hover:bg-green-200','text-green-800',
                    'dark:bg-green-900','dark:hover:bg-green-800','dark:text-green-200'
                );
                if (tick) tick.classList.remove('hidden');
            } else {
                btn.classList.remove(
                    'bg-green-100','hover:bg-green-200','text-green-800',
                    'dark:bg-green-900','dark:hover:bg-green-800','dark:text-green-200'
                );
                btn.classList.add(
                    'bg-white','hover:bg-gray-50','text-gray-800',
                    'dark:bg-gray-800','dark:hover:bg-gray-700','dark:text-gray-100'
                );
                if (tick) tick.classList.add('hidden');
            }
        });
    };

    // Toggle bottom-panel Verify visibility and global Calculate visibility
    window.updateActionButtons = function () {
        const allVisited = window.allHoldingsVisited && window.allHoldingsVisited();
        const calc = document.getElementById('calculate-btn');
        const verify = document.getElementById('verify-btn'); // in the panel

        if (calc) {
            if (allVisited) calc.classList.remove('hidden');
            else calc.classList.add('hidden');
        }
        if (verify) {
            // Verify button can remain visible; hide it only when all are visited
            if (allVisited) verify.classList.add('hidden');
            else verify.classList.remove('hidden');
        }

        window.refreshHoldingButtons();
    };

    // ---------- Public API (called from templates) ----------

    // Called once (Templ onload) with holdings slice
    window.initRequiredSymbols = function (data) {
        window.requiredSymbols = new Set((data || []).map(h => h.Symbol));
        // DO NOT auto-mark anything as visited here
        window.updateActionButtons();
    };

    // Called by each segment button onclick
    // NOTE: This now ONLY navigates; it DOES NOT mark as visited.
    window.handleHoldingClicked = function (_event, _symbol) {
        // intentionally do nothing to viewedHoldings
        // HTMX will load the panel; buttons/panel will sync after swap
    };

    // Called by the panel's full-width Verify button
    // This marks the CURRENT symbol as verified and updates UI.
    window.verifyCurrent = function (symbol) {
        window.viewedHoldings[symbol] = true;
        window.updateActionButtons();
    };

    // HTMX will post these values
    window.buildFIFSubmission = function () {
        return window.fifState;
    };

    // ---------- Event wiring ----------
    // Capture committed changes on number inputs (inside the form)
    document.addEventListener('change', (e) => {
        const el = e.target;
        if (!inForm(el)) return;
        if (!(el instanceof HTMLInputElement)) return;
        if (el.type !== 'number') return;
        if (!el.id) return;
        window.fifState[el.id] = el.value;
    });

    // Snapshot the current panel before HTMX sends a request
    document.body.addEventListener('htmx:beforeRequest', () => {
        const panel = document.getElementById('holding-panel');
        if (!panel) return;
        panel.querySelectorAll('input[type="number"][id]').forEach((el) => {
            window.fifState[el.id] = el.value;
        });
    });

    // After HTMX swaps in a new holding panel, rehydrate values and sync buttons
    document.body.addEventListener('htmx:afterSwap', (e) => {
        if (e.target && e.target.id === 'holding-panel') {
            e.target.querySelectorAll('input[type="number"][id]').forEach((el) => {
                if (el.id in window.fifState) el.value = window.fifState[el.id];
            });
            window.updateActionButtons(); // hide Verify if all done; refresh ticks/colors
        }
    });

    console.log('[FIF] verify-only flow initialized');
})();