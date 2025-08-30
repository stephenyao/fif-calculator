function handleHoldingClicked(event, symbol) {
    console.log(event.target.dataset?.symbol);
    window.viewedHoldings[symbol] = true;
    // After marking visited, see if we can show the Calculate button
    if (window.allHoldingsVisited && window.allHoldingsVisited()) {
        const btn = document.getElementById('calculate-btn');
        if (btn) btn.classList.remove('hidden');
    }
}

function initRequiredSymbols(data) {
    console.log('test');
    console.log(data);
    window.requiredSymbols = new Set((data || []).map(h => h.Symbol));
    // on init, keep button hidden
    const btn = document.getElementById('calculate-btn');
    if (btn) btn.classList.add('hidden');
}

(function () {
    if (window.__fifInit) return;
    window.__fifInit = true;

    // global cache of values keyed by input id, e.g. price_start_AAPL
    window.fifState = window.fifState || {};
    window.viewedHoldings = {};
    window.requiredSymbols = window.requiredSymbols || new Set();

    // returns true when every required symbol has been clicked
    window.allHoldingsVisited = function () {
        for (const sym of window.requiredSymbols) {
            if (!window.viewedHoldings[sym]) return false;
        }
        return window.requiredSymbols.size > 0;
    };

    // Debug so you can see the wiring happened
    console.log("[FIF] input capture initialized");

    // Only fire when the event comes from within the current form
    function inForm(target) {
        const form = document.getElementById('fif-form');
        return form && form.contains(target);
    }

    // 1) Capture on committed change (blur or enter)
    document.addEventListener('change', (e) => {
        const el = e.target;
        if (!inForm(el)) return;
        if (!(el instanceof HTMLInputElement)) return;
        if (el.type !== 'number') return;
        if (!el.id) return;

        window.fifState[el.id] = el.value;
        // console.log('[FIF] change:', el.id, '=', el.value);
    });

    // 3) Hydrate inputs inside #holding-panel after HTMX swaps
    document.body.addEventListener('htmx:afterSwap', (e) => {
        if (e.target && e.target.id === 'holding-panel') {
            e.target.querySelectorAll('input[type="number"][id]').forEach((el) => {
                if (el.id in window.fifState) el.value = window.fifState[el.id];
            });
        }
    });

    // 4) Snapshot current panel before a request in case user hasn't blurred
    document.body.addEventListener('htmx:beforeRequest', () => {
        const panel = document.getElementById('holding-panel');
        if (!panel) return;
        panel.querySelectorAll('input[type="number"][id]').forEach((el) => {
            window.fifState[el.id] = el.value;
        });
        // console.log('[FIF] snapshot before request');
    });

    // 5) Your HTMX POST will include these values:
    window.buildFIFSubmission = function () {
        return window.fifState; // keys like price_start_AAPL, etc.
    };
})();