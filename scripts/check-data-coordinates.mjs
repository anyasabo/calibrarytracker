import { readFile } from 'node:fs/promises';

const CA_BOUNDS = {
	latMin: 32.5,
	latMax: 42.0,
	lngMin: -124.5,
	lngMax: -114.0
};

function isFiniteNumber(value) {
	return Number.isFinite(value);
}

function validateCoord(record, kind) {
	const issues = [];
	const lat = record.lat;
	const lng = record.lng;
	const id = record.id ?? '(missing-id)';
	const name = record.name ?? '(missing-name)';

	if (!isFiniteNumber(lat) || !isFiniteNumber(lng)) {
		issues.push({
			kind,
			id,
			name,
			lat,
			lng,
			reason: 'lat/lng must be finite numbers'
		});
		return issues;
	}

	// (0,0) represents "no reliable coordinate yet" and is allowed.
	if (lat === 0 && lng === 0) {
		return issues;
	}

	if (lat === 0 || lng === 0) {
		issues.push({
			kind,
			id,
			name,
			lat,
			lng,
			reason: 'lat/lng must both be zero or both be non-zero'
		});
		return issues;
	}

	if (
		lat < CA_BOUNDS.latMin ||
		lat > CA_BOUNDS.latMax ||
		lng < CA_BOUNDS.lngMin ||
		lng > CA_BOUNDS.lngMax
	) {
		issues.push({
			kind,
			id,
			name,
			lat,
			lng,
			reason: 'coordinates are outside California bounds'
		});
	}

	return issues;
}

async function readJsonArray(path) {
	const raw = await readFile(path, 'utf8');
	const parsed = JSON.parse(raw);
	if (!Array.isArray(parsed)) {
		throw new Error(`${path} must contain a JSON array`);
	}
	return parsed;
}

async function main() {
	const [branches, systems] = await Promise.all([
		readJsonArray(new URL('../data/branches.json', import.meta.url)),
		readJsonArray(new URL('../data/systems.json', import.meta.url))
	]);

	const issues = [
		...branches.flatMap((branch) => validateCoord(branch, 'branch')),
		...systems.flatMap((system) => validateCoord(system, 'system'))
	];

	const missingBranches = branches.filter((b) => b.lat === 0 && b.lng === 0).length;
	const missingSystems = systems.filter((s) => s.lat === 0 && s.lng === 0).length;

	console.log(
		`Coordinate audit: branches=${branches.length} systems=${systems.length} missing(branches)=${missingBranches} missing(systems)=${missingSystems}`
	);

	if (issues.length === 0) {
		console.log('Coordinate audit passed.');
		return;
	}

	console.error(`Coordinate audit failed: found ${issues.length} invalid records.`);
	for (const issue of issues.slice(0, 30)) {
		console.error(
			`- [${issue.kind}] ${issue.id} (${issue.name}) lat=${issue.lat} lng=${issue.lng}: ${issue.reason}`
		);
	}
	if (issues.length > 30) {
		console.error(`...and ${issues.length - 30} more`);
	}

	process.exitCode = 1;
}

main().catch((error) => {
	console.error(
		`Coordinate audit crashed: ${error instanceof Error ? error.message : String(error)}`
	);
	process.exitCode = 1;
});
