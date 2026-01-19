import { EventsOn, EventsOff, Invoke } from '@wails/runtime';

export function subscribeWatcherEvents(cb) {
  EventsOn('watcher:event', cb);
  return () => EventsOff('watcher:event', cb);
}

export async function addWatcher(entry) {
  return Invoke('AddWatcher', entry);
}
export async function updateWatcher(entry) {
  return Invoke('UpdateWatcher', entry);
}
export async function deleteWatcher(id) {
  return Invoke('DeleteWatcher', id);
}
export async function listWatchers() {
  return Invoke('ListWatchers');
}
export async function startWatcher(id) {
  return Invoke('StartWatcher', id);
}
export async function stopWatcher(id) {
  return Invoke('StopWatcher', id);
}
export async function listJobs(watcherId, status) {
  return Invoke('ListJobs', watcherId, status);
}
export async function listLogs(watcherId, level) {
  return Invoke('ListLogs', watcherId, level);
}
