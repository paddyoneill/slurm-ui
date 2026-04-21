export const formatAge = (origin: Date | string): string => {
	const now = new Date();
	const creation = origin instanceof Date ? origin : new Date(origin);

	const msMinute = 1000 * 60;
	const msHour = msMinute * 60;
	const msDay = msHour * 24;

	let diff = now.getTime() - creation.getTime();

	if (diff < msMinute) return '0 minutes';

	const days = Math.floor(diff / msDay);
	diff %= msDay;

	const hours = Math.floor(diff / msHour);
	diff %= msHour;

	const minutes = Math.floor(diff / msMinute);

	const parts: string[] = [];
	if (days > 0) parts.push(`${days} ${days === 1 ? 'day' : 'days'}`);
	if (hours > 0) parts.push(`${hours} ${hours === 1 ? 'hour' : 'hours'}`);
	if (minutes > 0) parts.push(`${minutes} ${minutes === 1 ? 'minute' : 'minutes'}`);

	return parts.join(',');
};
