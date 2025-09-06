import {Events} from "@wailsio/runtime";
import * as Journal from "./bindings/changeme/backend/bindings/journalservice.js";
import * as Models from "./bindings/changeme/backend/models/models.js";

const timeElement = document.getElementById('time');

// Wire Ping button
const pingBtn = document.getElementById('ping-btn');
const pingResult = document.getElementById('ping-result');
if (pingBtn && pingResult) {
    pingBtn.addEventListener('click', async () => {
        try {
            const v = await Journal.Ping();
            pingResult.innerText = `version: ${v}`;
        } catch (e) {
            pingResult.innerText = 'ping failed';
        }
    });
}

// Filters → Query → GetAnalytics
const applyBtn = document.getElementById('apply-filters');
if (applyBtn) {
    applyBtn.addEventListener('click', async () => {
        const symbol = document.getElementById('filter-symbol')?.value || '';
        const side = document.getElementById('filter-side')?.value || '';
        const start = document.getElementById('filter-start')?.value || '';
        const end = document.getElementById('filter-end')?.value || '';
        /** @type {ReturnType<typeof Models.Query.createFrom>} */
        const q = Models.Query.createFrom({
            symbol,
            side,
            startTime: start ? new Date(start).toISOString() : undefined,
            endTime: end ? new Date(end).toISOString() : undefined,
            limit: 1000,
            offset: 0,
        });
        try {
            const a = await Journal.GetAnalytics(q);
            setMetric('winRate', a.winRate?.toFixed(2));
            setMetric('profitFactor', a.profitFactor?.toFixed(2));
            setMetric('maxDD', a.maxDD?.toFixed(2));
            setMetric('sharpe', a.sharpe?.toFixed(2));
            setMetric('sortino', a.sortino?.toFixed(2));
            setMetric('expectancy', a.expectancy?.toFixed(2));
        } catch (e) {
            console.error(e);
        }
    });
}

function setMetric(name, value) {
    const el = document.getElementById(`metric-${name}`);
    if (el) el.innerText = value ?? '-';
}

// CSV Import
const importBtn = document.getElementById('import-btn');
const importRes = document.getElementById('import-result');
if (importBtn && importRes) {
    importBtn.addEventListener('click', async () => {
        const txt = document.getElementById('csv-text')?.value || '';
        try {
            const report = await Journal.ImportCSV(txt);
            importRes.innerText = `imported: ${report.imported}, skipped: ${report.skipped}, errors: ${report.errors?.length}`;
        } catch (e) {
            importRes.innerText = 'import failed';
        }
    });
}

Events.On('time', (time) => {
    timeElement.innerText = time.data;
});
