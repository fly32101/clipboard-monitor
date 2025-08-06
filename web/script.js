// 全局变量
let currentHistory = [];

// 更新状态
function updateStatus(text) {
    document.getElementById('status').textContent = '⚡ システム状態: ' + text;
}

// 渲染历史记录
function renderHistory() {
    console.log('开始渲染历史记录, 数量:', currentHistory.length);
    const container = document.getElementById('historyContainer');
    if (!container) {
        console.error('找不到 historyContainer 元素');
        return;
    }

    if (currentHistory.length === 0) {
        console.log('显示空状态');
        container.innerHTML = '<div class="empty-state"><div style="font-size: 3rem; margin-bottom: 15px; opacity: 0.7;">🌟</div><p>データなし</p><p style="font-size: 0.8rem; margin-top: 5px; opacity: 0.6;">何かをコピーして開始しましょう</p></div>';
        return;
    }

    console.log('渲染', currentHistory.length, '条记录');
    container.innerHTML = '';
    currentHistory.forEach((entry, index) => {
        console.log('渲染第', index, '条:', entry);
        const item = document.createElement('div');
        item.className = 'history-item';
        const time = new Date(entry.Timestamp || entry.timestamp);
        const timeStr = time.toLocaleTimeString('zh-CN', { hour12: false });
        let content = entry.Content || entry.content || '';
        if (content.length > 200) content = content.substring(0, 200) + '...';

        item.innerHTML = '<div class="item-time">' + timeStr + '</div><div class="item-content">' + escapeHtml(content) + '</div>';
        item.ondblclick = () => copyToClipboard(entry.Content || entry.content, index);
        container.appendChild(item);
    });
    console.log('历史记录渲染完成');
}

// HTML 转义
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 刷新历史记录
async function refreshHistory() {
    try {
        if (typeof getHistory === 'function') {
            const result = getHistory();
            console.log('getHistory 返回:', result, '类型:', typeof result);

            // 检查是否是 Promise
            if (result && typeof result.then === 'function') {
                console.log('检测到 Promise，等待结果...');
                const data = await result;
                console.log('Promise 解析结果:', data, '类型:', typeof data);
                currentHistory = Array.isArray(data) ? data : [];
            } else if (Array.isArray(result)) {
                currentHistory = result;
            } else if (result === null || result === undefined) {
                currentHistory = [];
            } else {
                console.log('尝试解析为数组...');
                currentHistory = [];
            }

            console.log('获取到历史记录:', currentHistory.length, '条');
        } else {
            console.log('getHistory 函数不可用');
            currentHistory = [];
        }
        renderHistory();
        updateStatus('リスト更新完了');
    } catch (error) {
        console.error('获取历史记录失败:', error);
        updateStatus('获取历史记录失败');
        currentHistory = [];
        renderHistory();
    }
}

// 复制到剪贴板
async function copyToClipboard(content, index) {
    try {
        if (typeof copyToClipboardGo === 'function') {
            const result = copyToClipboardGo(content);

            // 处理可能的 Promise
            let response = result;
            if (result && typeof result.then === 'function') {
                response = await result;
            }

            if (response && response.error) {
                throw new Error(response.error);
            }
        } else {
            await navigator.clipboard.writeText(content);
        }
        updateStatus('コピー完了！');
    } catch (error) {
        console.error('复制失败:', error);
        updateStatus('复制失败');
    }
}

// 清空历史记录
async function clearHistoryFunc() {
    if (confirm('本当にすべてのデータを消去しますか？')) {
        try {
            if (typeof clearHistory === 'function') {
                const result = clearHistory();
                // 处理可能的 Promise
                if (result && typeof result.then === 'function') {
                    await result;
                }
            }
            currentHistory = [];
            renderHistory();
            updateStatus('データ消去完了');
        } catch (error) {
            console.error('清空失败:', error);
            updateStatus('清空失败');
        }
    }
}

// 显示关于信息
async function showAbout() {
    let version = '⚡ クリップボード・モニター v2.0';
    try {
        if (typeof getVersionInfo === 'function') {
            version = getVersionInfo();
        }
    } catch (error) {}

    alert(version + '\n\n🌟 機能:\n- リアルタイム監視\n- ヒストリー表示\n- ダブルクリックコピー\n- アニメスタイル UI\n\n👨‍💻 作者: Fly\n🛠️ 技術: Go + WebView');
}

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    console.log('页面加载完成');
    console.log('可用函数:', typeof getHistory, typeof copyToClipboardGo, typeof clearHistory, typeof getVersionInfo);

    // 测试基本功能
    try {
        updateStatus('クリップボード監視中...');
        console.log('状态更新成功');

        refreshHistory();
        console.log('首次刷新完成');

        setInterval(refreshHistory, 2000);
        console.log('定时器设置完成');
    } catch (error) {
        console.error('初始化错误:', error);
    }
});
