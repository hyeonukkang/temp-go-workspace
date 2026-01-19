import React, { useEffect, useState } from 'react';
import { listJobs } from './api';

function Jobs({ watcherId }) {
  const [jobs, setJobs] = useState([]);
  useEffect(() => {
    if (!watcherId) return;
    listJobs(watcherId, '').then(setJobs);
  }, [watcherId]);

  return (
    <div>
      <h3>Jobs</h3>
      <table>
        <thead>
          <tr><th>ID</th><th>File</th><th>Status</th><th>Updated</th><th>Error</th></tr>
        </thead>
        <tbody>
          {jobs.map(j => (
            <tr key={j.id}>
              <td>{j.id}</td>
              <td>{j.filePath}</td>
              <td>{j.status}</td>
              <td>{new Date(j.updatedAt).toLocaleString()}</td>
              <td>{j.errorMsg}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default Jobs;
