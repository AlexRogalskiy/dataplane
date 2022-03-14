import { useEffect, useRef, useState } from 'react';
import ConsoleLogHelper from '../../../Helper/logger';
import { useGlobalRunState } from '../../../pages/View/useWebSocket';

var loc = window.location, new_uri;
if (loc.protocol === "https:") {
    new_uri = "wss:";
} else {
    new_uri = "ws:";
}
new_uri += "//" + loc.host;

// console.log("websockets loc:", new_uri)
if (process.env.REACT_APP_DATAPLANE_ENV == "build"){
    new_uri += process.env.REACT_APP_WEBSOCKET_ROOMS_ENDPOINT;
}else{
    new_uri = process.env.REACT_APP_WEBSOCKET_ROOMS_ENDPOINT;
}

const websocketEndpoint = new_uri;

export default function useWebSocketLog(environmentId, run_id, node_id) {
    const [socketResponseWithUID, setSocketResponseWithUID] = useState('');
    const reconnectOnClose = useRef(true);
    const ws = useRef(null);

    // Global state
    const RunState = useGlobalRunState();

    useEffect(() => {
        if (RunState.selectedNodeStatus.get() !== 'Run') return;

        function connect() {
            ws.current = new WebSocket(`${websocketEndpoint}/${environmentId}?subject=workerlogs.${run_id}.${node_id}&id=${run_id}.${node_id}`);

            ws.current.onopen = () => ConsoleLogHelper('ws opened');
            ws.current.onclose = () => {
                // Exit if closing the connection was intentional
                if (!reconnectOnClose.current) {
                    return;
                }

                // Reconnect on close
                ConsoleLogHelper('ws trying to reconnect...');
                setTimeout(function () {
                    connect();
                }, 1000);
            };

            ws.current.onmessage = (e) => {
                const resp = JSON.parse(e.data);
                // Return if not a log message
                if (resp.run_id) return;

                let text = `${formatDate(resp.created_at)} ${resp.log} -${resp.uid}`;
                setSocketResponseWithUID(text);
            };
        }

        connect();

        return () => {
            reconnectOnClose.current = false;
            ws.current.close();
        };

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [run_id]);

    return socketResponseWithUID;
}

export function formatDate(dateString) {
    const date = new Date(dateString);
    return (
        date.getFullYear() +
        '/' +
        ('0' + (date.getMonth() + 1)) +
        '/' +
        ('0' + date.getDate()).slice(-2) +
        ' ' +
        ('0' + date.getHours()).slice(-2) +
        ':' +
        ('0' + date.getMinutes()).slice(-2) +
        ':' +
        ('0' + date.getSeconds()).slice(-2)
    );
}
