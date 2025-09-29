(function () {
    if (window.__fifInit) return;
    window.__fifInit = true;

    window.fifState = window.fifState || {};
    window.viewedHoldings = window.viewedHoldings || {};
    window.requiredSymbols = window.requiredSymbols || new Set();

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

    window.updateActionButtons = function () {
        const allVisited = window.allHoldingsVisited && window.allHoldingsVisited();
        const calc = document.getElementById('calculate-btn');
        const markBtn = document.getElementById('markdone-btn');

        if (calc) {
            if (allVisited) calc.classList.remove('hidden');
            else calc.classList.add('hidden');
        }
        if (markBtn) {
            if (allVisited) markBtn.classList.add('hidden');
            else markBtn.classList.remove('hidden');
        }

        window.refreshHoldingButtons();
    };

    // Called once with holdings slice
    window.initRequiredSymbols = function (data) {
        window.requiredSymbols = new Set((data || []).map(h => h.Symbol));
        window.updateActionButtons();
    };

    // Only marks as done when the Mark button is clicked
    window.markHoldingDone = function (symbol) {
        window.viewedHoldings[symbol] = true;
        window.updateActionButtons();
    };

    // Reset "done" if any input changes in that holding’s panel
    function resetDoneOnInputChange(symbol) {
        window.viewedHoldings[symbol] = false;
        window.updateActionButtons();
    }

    // Track changes in numeric inputs
    document.addEventListener('input', (e) => {
        const el = e.target;
        if (!inForm(el)) return;
        if (!(el instanceof HTMLInputElement)) return;
        if (el.type !== 'number') return;
        if (!el.id) return;

        window.fifState[el.id] = el.value;

        // symbol is encoded in input id e.g. "price_start_AAPL"
        const parts = el.id.split('_');
        const sym = parts[parts.length - 1];
        if (sym && window.requiredSymbols.has(sym)) {
            resetDoneOnInputChange(sym);
        }
    });

    // Snapshot before request
    document.body.addEventListener('htmx:beforeRequest', () => {
        const panel = document.getElementById('holding-panel');
        if (!panel) return;
        panel.querySelectorAll('input[type="number"][id]').forEach((el) => {
            window.fifState[el.id] = el.value;
        });
    });

    // Rehydrate and update after swap
    document.body.addEventListener('htmx:afterSwap', (e) => {
        if (e.target && e.target.id === 'holding-panel') {
            e.target.querySelectorAll('input[type="number"][id]').forEach((el) => {
                if (el.id in window.fifState) el.value = window.fifState[el.id];
            });
            window.updateActionButtons();
        }
    });

    window.buildFIFSubmission = function () {
        return window.fifState;
    };

    window.clearFIFState = function() {
        window.fifState = {};
        window.viewedHoldings = {};
    }

    console.log('[FIF] mark-as-done flow with reset-on-input-change initialized');
})();