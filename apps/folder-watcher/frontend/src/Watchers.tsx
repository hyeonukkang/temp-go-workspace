
import React, { useEffect, useState } from 'react';
import {
  listWatchers,
  addWatcher,
  updateWatcher,
  deleteWatcher,
  startWatcher,
  stopWatcher,
  subscribeWatcherEvents
} from './api';
import Logs from './Logs';
import Jobs from './Jobs';

function Watchers() {
  const [watchers, setWatchers] = useState([]);
  const [selected, setSelected] = useState(null);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    fetchList();
    const unsub = subscribeWatcherEvents((evt) => {
      setToast(evt.message);
      fetchList();
    });
    return unsub;
  }, []);

  function fetchList() {
    listWatchers().then(setWatchers);
  }

  function handleStart(id) {
    startWatcher(id);
  }
  function handleStop(id) {
    stopWatcher(id);
  }
  function handleDelete(id) {
    deleteWatcher(id).then(fetchList);
  }

  return (
    <div>
      <h2>Watchers</h2>
      {toast && <div className="toast">{toast}</div>}
      <table>
        <thead>
          <tr>
            <th>ID</th><th>Source</th><th>Output</th><th>Action</th><th>OnStart</th><th>...</th>
          </tr>
        </thead>
        <tbody>
          {watchers.map(w => (
            <tr key={w.id} onClick={() => setSelected(w)} style={{background:selected?.id===w.id?'#eef':''}}>
              <td>{w.id}</td>
              <td>{w.sourcePath}</td>
              <td>{w.outputPath}</td>
              <td>{w.action}</td>
              <td>{w.onStartBehavior}</td>
              <td>
                <button onClick={e => {e.stopPropagation(); handleStart(w.id);}}>Start</button>
                <button onClick={e => {e.stopPropagation(); handleStop(w.id);}}>Stop</button>
                <button onClick={e => {e.stopPropagation(); handleDelete(w.id);}}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {selected && <>
        <Jobs watcherId={selected.id} />
        <Logs watcherId={selected.id} />
      </>}
      {/* Add/Edit 폼 등은 추후 확장 */}
    </div>
  );
}

export default Watchers;
