type Listener = () => void;

const listeners = new Set<Listener>();

export function onJobsChanged(listener: Listener) {
	listeners.add(listener);
	return () => listeners.delete(listener);
}

export function notifyJobsChanged() {
	listeners.forEach((listener) => listener());
}
