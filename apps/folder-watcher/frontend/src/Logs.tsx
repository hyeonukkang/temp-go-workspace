import React, { useEffect, useState } from 'react';
import { listLogs } from './api';

function Logs({ watcherId }) {
  const [logs, setLogs] = useState([]);
  useEffect(() => {
    if (!watcherId) return;
    listLogs(watcherId, '').then(setLogs);
  }, [watcherId]);

  return (
    <div>
      <h3>Logs</h3>
      <table>
        <thead>
          <tr><th>Time</th><th>Level</th><th>Message</th></tr>
        </thead>
        <tbody>
          {logs.map(l => (
            <tr key={l.id}>
              <td>{new Date(l.createdAt).toLocaleString()}</td>
              <td>{l.level}</td>
              <td>{l.message}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default Logs;
