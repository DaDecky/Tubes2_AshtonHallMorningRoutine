"use client"

import { useState } from 'react';
import { useWebsocket } from '../../lib/websocket';

export default function Test() {

    const { connect, progress, status, error, result } = useWebsocket();
    const [algo, setAlgo] = useState('');
    const [mode, setMode] = useState('');
    const [max, setMax] = useState('');

    // connect to websocket with params
    const handleStart = () => {
        const params: Record<string, string> = {};
        if (algo) params.algo = algo;
        if (mode) params.mode = mode;
        if (max) params.max = max;
        connect(params);
    };

    return (
        <div>
            <h2>Test Websocket</h2>
            
            <div>
                <label>
                algo: 
                <input
                    type="text"
                    value={algo}
                    onChange={(e) => setAlgo(e.target.value)}
                />
                </label>
            </div>

            <div>
                <label>
                mode: 
                <input
                    type="text"
                    value={mode}
                    onChange={(e) => setMode(e.target.value)}
                />
                </label>
            </div>

            <div>
                <label>
                max-recipe: 
                <input
                    type="text"
                    value={max}
                    onChange={(e) => setMax(e.target.value)}
                />
                </label>
            </div>

            <div>
                <h3>Status : {status}</h3>
                <h3>Results :</h3>
                <pre>{ result && (JSON.stringify(result, null, 2))}</pre>
            </div>

            <div>
                <button onClick={handleStart} className=' bg-red-100 text-black rounded' disabled={status === 'running'}>
                Start Search
                </button>
                
            </div>

            {status === 'running' && (
                <div>
                <progress value={progress} max="10" />
                <span>{progress}%</span>
                </div>
            )}

            {status === 'error' && (
                <div style={{ color: 'red' }}>
                {error || 'An error occurred'}
                </div>
            )}

            {status === 'complete' && (
                <div>
                <h3>Results:</h3>
                <pre>{JSON.stringify(result, null, 2)}</pre>
                </div>
            )}

            
        </div>
    );
}