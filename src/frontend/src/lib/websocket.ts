// import { useState, useCallback, useEffect } from 'react';

// export const useWebsocket = () => {
//   const [socket, setSocket] = useState<WebSocket | null>(null);
//   const [progress, setProgress] = useState(0);
//   const [status, setStatus] = useState<'idle' | 'connecting' | 'running' | 'complete' | 'error'>('idle');
//   const [error, setError] = useState<string | null>(null);
//     // eslint-disable-next-line @typescript-eslint/no-explicit-any
//   const [result, setResult] = useState<any>(null);

//   const connect = useCallback(async (params: Record<string, string> = {}) => {
//     // First check if we should use WebSocket at all
//     if (Object.keys(params).length === 0) {
//       setError('No parameters provided');
//       setStatus('error');
//       return;
//     }

//     setStatus('connecting');
//         try {
//             const queryString = new URLSearchParams(params).toString();
//             const ws = new WebSocket(`ws://localhost:8081/liveSearch?${queryString}`);

//             ws.onopen = () => {
//                 setStatus('running');
//                 setProgress(0);
//                 setError(null);
//             };

//             ws.onmessage = (event) => {
//                 const data = JSON.parse(event.data);
//                 switch (data.type) {
//                     case 'progress':
//                         setProgress(data.progress_counter);
//                         break;
//                     case 'complete':
//                         setResult(data);
//                         setStatus('complete');
//                         ws.close();
//                         break;
//                 }
//             };

//             ws.onerror = () => {
//                 setError('Connection failed');
//                 setStatus('error');
//             };

//             ws.onclose = () => {
//                 if (status !== 'complete') {
//                     setStatus('idle');
//                 }

//             };

//             setSocket(ws);
//         } 
//         // eslint-disable-next-line @typescript-eslint/no-unused-vars
//         catch (err) {
//                 setError('Failed to establish connection');
//                 setStatus('error');
//         }
//     }, [status]);

//     // Cleanup on unmount
//     useEffect(() => {
//         return () => {
//             if (socket) {
//             socket.close();
//             }
//         };
//     }, [socket]);

//     return { connect, progress, status, error, result };
// };