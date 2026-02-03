package api

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>股票监控管理</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }
        .container { max-width: 1000px; margin: 0 auto; padding: 20px; }
        h1 { text-align: center; margin-bottom: 30px; color: #333; }
        .card { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .card h2 { margin-bottom: 15px; color: #444; font-size: 18px; border-bottom: 1px solid #eee; padding-bottom: 10px; }
        .form-row { display: flex; gap: 10px; margin-bottom: 10px; flex-wrap: wrap; }
        input, select { padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
        input:focus, select:focus { outline: none; border-color: #4a90d9; }
        button { padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; }
        .btn-primary { background: #4a90d9; color: white; }
        .btn-danger { background: #e74c3c; color: white; }
        .btn-success { background: #27ae60; color: white; }
        button:hover { opacity: 0.9; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f9f9f9; font-weight: 500; }
        .tag { display: inline-block; padding: 2px 8px; border-radius: 3px; font-size: 12px; }
        .tag-info { background: #e3f2fd; color: #1976d2; }
        .tag-warning { background: #fff3e0; color: #f57c00; }
        .tag-critical { background: #ffebee; color: #c62828; }
        .switch { position: relative; display: inline-block; width: 44px; height: 24px; }
        .switch input { opacity: 0; width: 0; height: 0; }
        .slider { position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0; background: #ccc; border-radius: 24px; transition: .3s; }
        .slider:before { position: absolute; content: ""; height: 18px; width: 18px; left: 3px; bottom: 3px; background: white; border-radius: 50%; transition: .3s; }
        input:checked + .slider { background: #4a90d9; }
        input:checked + .slider:before { transform: translateX(20px); }
    </style>
</head>
<body>
    <div class="container">
        <h1>📈 股票监控管理</h1>

        <div class="card">
            <h2>股票管理</h2>
            <div class="form-row">
                <input type="text" id="stockCode" placeholder="股票代码 如600519">
                <input type="text" id="stockName" placeholder="股票名称">
                <button class="btn-primary" onclick="addStock()">添加股票</button>
            </div>
            <table><thead><tr><th>代码</th><th>名称</th><th>操作</th></tr></thead><tbody id="stockList"></tbody></table>
        </div>

        <div class="card">
            <h2>规则管理</h2>
            <div class="form-row">
                <input type="text" id="ruleName" placeholder="规则名称">
                <select id="ruleType">
                    <option value="price_above_ma">突破均线</option>
                    <option value="price_below_ma">跌破均线</option>
                </select>
                <select id="ruleStock"><option value="">选择股票</option></select>
                <select id="ruleKline">
                    <option value="5min">5分钟</option>
                    <option value="15min">15分钟</option>
                    <option value="30min">30分钟</option>
                    <option value="60min">60分钟</option>
                    <option value="daily" selected>日K</option>
                </select>
                <input type="number" id="rulePeriod" placeholder="MA周期" value="60" style="width:80px">
                <select id="ruleLevel">
                    <option value="info">Info</option>
                    <option value="warning" selected>Warning</option>
                    <option value="critical">Critical</option>
                </select>
                <button class="btn-primary" onclick="addRule()">添加规则</button>
            </div>
            <table><thead><tr><th>名称</th><th>类型</th><th>股票</th><th>周期</th><th>MA</th><th>级别</th><th>启用</th><th>操作</th></tr></thead><tbody id="ruleList"></tbody></table>
        </div>

        <div class="card">
            <h2>通知配置</h2>
            <div class="form-row">
                <label><input type="checkbox" id="feishuEnabled"> 飞书</label>
                <input type="text" id="feishuWebhook" placeholder="飞书 Webhook URL" style="flex:1">
            </div>
            <div class="form-row">
                <label><input type="checkbox" id="serverchanEnabled"> Server酱</label>
                <input type="text" id="serverchanKey" placeholder="SendKey" style="flex:1">
            </div>
            <div class="form-row">
                <label><input type="checkbox" id="dingtalkEnabled"> 钉钉</label>
                <input type="text" id="dingtalkWebhook" placeholder="钉钉 Webhook URL" style="flex:1">
            </div>
            <button class="btn-success" onclick="saveNotifiers()">保存通知配置</button>
        </div>
    </div>

    <script>
        let stocks = [], rules = [], ruleTypes = {};
        const api = (url, opt) => fetch(url, opt).then(r => r.json());

        async function loadRuleTypes() {
            const types = await api('/api/rule-types');
            ruleTypes = {};
            const select = document.getElementById('ruleType');
            select.innerHTML = types.map(t => {
                ruleTypes[t.type] = t.name;
                return ` + "`" + `<option value="${t.type}">${t.name}</option>` + "`" + `;
            }).join('');
        }

        async function loadStocks() {
            stocks = await api('/api/stocks');
            document.getElementById('stockList').innerHTML = stocks.map(s =>
                ` + "`" + `<tr><td>${s.code}</td><td>${s.name}</td><td><button class="btn-danger" onclick="delStock('${s.code}')">删除</button></td></tr>` + "`" + `
            ).join('');
            document.getElementById('ruleStock').innerHTML = '<option value="">全部股票</option>' +
                stocks.map(s => ` + "`" + `<option value="${s.code}">${s.name}</option>` + "`" + `).join('');
        }

        async function addStock() {
            const code = document.getElementById('stockCode').value;
            const name = document.getElementById('stockName').value;
            if (!code || !name) return alert('请填写完整');
            await api('/api/stocks', {method:'POST', body:JSON.stringify({code,name})});
            document.getElementById('stockCode').value = '';
            document.getElementById('stockName').value = '';
            loadStocks();
        }

        async function delStock(code) {
            await api('/api/stocks?code='+code, {method:'DELETE'});
            loadStocks();
        }

        const klineNames = {'5min':'5分钟','15min':'15分钟','30min':'30分钟','60min':'60分钟','daily':'日K'};

        async function loadRules() {
            rules = await api('/api/rules');
            document.getElementById('ruleList').innerHTML = rules.map(r => ` + "`" + `
                <tr>
                    <td>${r.name}</td>
                    <td>${ruleTypes[r.type] || r.type}</td>
                    <td>${r.stock_code || '全部'}</td>
                    <td>${klineNames[r.kline_type] || r.kline_type}</td>
                    <td>MA${r.period}</td>
                    <td><span class="tag tag-${r.level}">${r.level}</span></td>
                    <td><label class="switch"><input type="checkbox" ${r.enabled?'checked':''} onchange="toggleRule('${r.id}',this.checked)"><span class="slider"></span></label></td>
                    <td><button class="btn-danger" onclick="delRule('${r.id}')">删除</button></td>
                </tr>
            ` + "`" + `).join('');
        }

        async function addRule() {
            const name = document.getElementById('ruleName').value;
            if (!name) return alert('请填写规则名称');
            await api('/api/rules', {method:'POST', body:JSON.stringify({
                name, type: document.getElementById('ruleType').value, enabled:true,
                stock_code: document.getElementById('ruleStock').value,
                kline_type: document.getElementById('ruleKline').value,
                period: parseInt(document.getElementById('rulePeriod').value),
                level: document.getElementById('ruleLevel').value
            })});
            document.getElementById('ruleName').value = '';
            loadRules();
        }

        async function toggleRule(id, enabled) {
            const rule = rules.find(r => r.id === id);
            rule.enabled = enabled;
            await api('/api/rules', {method:'PUT', body:JSON.stringify(rule)});
        }

        async function delRule(id) {
            await api('/api/rules?id='+id, {method:'DELETE'});
            loadRules();
        }

        async function loadNotifiers() {
            const n = await api('/api/notifiers');
            document.getElementById('feishuEnabled').checked = n.feishu?.enabled;
            document.getElementById('feishuWebhook').value = n.feishu?.webhook || '';
            document.getElementById('serverchanEnabled').checked = n.serverchan?.enabled;
            document.getElementById('serverchanKey').value = n.serverchan?.send_key || '';
            document.getElementById('dingtalkEnabled').checked = n.dingtalk?.enabled;
            document.getElementById('dingtalkWebhook').value = n.dingtalk?.webhook || '';
        }

        async function saveNotifiers() {
            await api('/api/notifiers', {method:'PUT', body:JSON.stringify({
                feishu: {enabled: document.getElementById('feishuEnabled').checked, webhook: document.getElementById('feishuWebhook').value},
                serverchan: {enabled: document.getElementById('serverchanEnabled').checked, send_key: document.getElementById('serverchanKey').value},
                dingtalk: {enabled: document.getElementById('dingtalkEnabled').checked, webhook: document.getElementById('dingtalkWebhook').value}
            })});
            alert('保存成功');
        }

        loadRuleTypes(); loadStocks(); loadRules(); loadNotifiers();
    </script>
</body></html>`
