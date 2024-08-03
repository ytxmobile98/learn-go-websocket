window.onload = async () => {
    let conn;

    const msg = document.getElementById("msg");
    const log = document.getElementById("log");
    const form = document.getElementById("form");

    function appendLog(item) {
        const doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    async function newConnection() {
        conn = new WebSocket("ws://" + document.location.host + "/ws");

        let resolve, reject;
        const promise = new Promise((res, rej) => {
            resolve = res;
            reject = rej;
        });

        conn.onclose = event => {
            reject(event);
        }
        conn.onopen = () => {
            resolve(conn);
        };

        conn.onclose = () => {
            const item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);

            conn = undefined;
        };

        conn.onmessage = function (evt) {
            const messages = evt.data.split('\n');
            for (const message of messages) {
                const item = document.createElement("div");
                item.innerText = message;
                appendLog(item);
            }
        };

        return promise;
    }

    form.onsubmit = () => {
        queueMicrotask(async () => {
            if (!conn) {
                try {
                    conn = await newConnection();

                    const item = document.createElement("div");
                    item.innerHTML = "<b>Connection reopened.</b>";
                    appendLog(item);
                } catch (err) {
                    conn = undefined;
                    return;
                }
            }

            if (!msg.value) {
                return;
            }
            conn.send(msg.value);
            msg.value = "";
        });

        return false;
    };

    if (window.WebSocket) {
        conn = await newConnection();
    } else {
        const item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};