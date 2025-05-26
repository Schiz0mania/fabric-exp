const { exec } = require("child_process");
const { promisify } = require("util");
const async = require("async");

const concurrency = 200;
const batchSize = 10;

console.log(`Total queries: ${concurrency}, Batch size: ${batchSize}`);

const execAsync = promisify(exec);
const durations = [];

// 关键修复：移除 callback 参数，完全基于 Promise 控制流
const queryQueue = async.queue(async (taskId) => { // 不再接收 callback 参数
    const start = Date.now();
    try {
        await execAsync("node queryAll.js");
        const duration = Date.now() - start;
        durations.push(duration);
        console.log(`Query ${taskId} completed in ${duration}ms`);
        return duration; // 通过返回值标记成功
    } catch (error) {
        console.error(`Query ${taskId} failed: ${error.message}`);
        throw error; // 通过抛出错误标记失败
    }
}, batchSize);

// 填充任务队列（保持原样）
for (let i = 1; i <= concurrency; i++) {
    queryQueue.push(i);
}

// 结果统计（增加错误计数）
let failedCount = 0;
queryQueue.drain(() => {
    const successCount = durations.length;
    const avg = successCount > 0
        ? durations.reduce((a, b) => a + b) / successCount
        : 0;
    console.log(`Average: ${avg.toFixed(2)}ms (${successCount}/${concurrency} succeeded, ${failedCount} failed)`);
});

// 错误处理（单独记录失败计数）
queryQueue.error((err) => {
    failedCount++;
    console.error("Task error:", err.message);
});