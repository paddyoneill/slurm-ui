type Listener = () => void;

const jobListeners = new Set<Listener>();

export function onJobsChanged(listener: Listener) {
	jobListeners.add(listener);
	return () => jobListeners.delete(listener);
}

export function notifyJobsChanged() {
	jobListeners.forEach((listener) => listener());
}

const notebookListeners = new Set<Listener>();

export function onNotebooksChanged(listener: Listener) {
	notebookListeners.add(listener);
	return () => notebookListeners.delete(listener);
}

export function notifyNotebooksChanged() {
	notebookListeners.forEach((listener) => listener());
}
