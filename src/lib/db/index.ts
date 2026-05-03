import Dexie, { type EntityTable } from 'dexie';
import type { UserCard, BranchVisit } from '$lib/types';

/**
 * Client-side database for user data. Uses IndexedDB via Dexie.
 *
 * IMPORTANT: Data in this database never leaves the user's device.
 * There is no server, no sync, no telemetry. Users can export/import
 * their data as JSON for backup.
 *
 * Schema versioning: when adding fields or tables, increment the version
 * number and add a migration. Dexie handles schema upgrades automatically
 * when the user next opens the app. See:
 * https://dexie.org/docs/Tutorial/Design#database-versioning
 */
const db = new Dexie('CalibraryTracker') as Dexie & {
	cards: EntityTable<UserCard, 'id'>;
	visits: EntityTable<BranchVisit, 'id'>;
};

db.version(1).stores({
	// Index on systemId so we can quickly look up "all cards for system X"
	// and on status so we can filter active/expired
	cards: 'id, systemId, status',
	// Index on branchId so we can quickly check "has user visited branch X"
	visits: 'id, branchId, visitDate'
});

export { db };

// -- Helper functions for common operations --

export async function getAllCards(): Promise<UserCard[]> {
	return db.cards.toArray();
}

export async function getCardsBySystem(systemId: string): Promise<UserCard[]> {
	return db.cards.where('systemId').equals(systemId).toArray();
}

export async function addCard(card: UserCard): Promise<string> {
	return db.cards.add(card);
}

export async function updateCard(card: UserCard): Promise<string> {
	await db.cards.put(card);
	return card.id;
}

export async function deleteCard(id: string): Promise<void> {
	await db.cards.delete(id);
}

export async function getAllVisits(): Promise<BranchVisit[]> {
	return db.visits.toArray();
}

export async function getVisitsByBranch(branchId: string): Promise<BranchVisit[]> {
	return db.visits.where('branchId').equals(branchId).toArray();
}

export async function addVisit(visit: BranchVisit): Promise<string> {
	return db.visits.add(visit);
}

export async function deleteVisit(id: string): Promise<void> {
	await db.visits.delete(id);
}

/**
 * Export all user data as a JSON-serializable object.
 * This is the backup/restore format.
 */
export async function exportData(): Promise<{ cards: UserCard[]; visits: BranchVisit[] }> {
	const [cards, visits] = await Promise.all([db.cards.toArray(), db.visits.toArray()]);
	return { cards, visits };
}

/**
 * Import user data from a backup. Merges with existing data by ID
 * (upsert semantics — existing records with the same ID are overwritten).
 */
export async function importData(data: {
	cards?: UserCard[];
	visits?: BranchVisit[];
}): Promise<void> {
	await db.transaction('rw', db.cards, db.visits, async () => {
		if (data.cards) {
			await db.cards.bulkPut(data.cards);
		}
		if (data.visits) {
			await db.visits.bulkPut(data.visits);
		}
	});
}

/**
 * Delete ALL user data. This is irreversible.
 */
export async function clearAllData(): Promise<void> {
	await db.transaction('rw', db.cards, db.visits, async () => {
		await db.cards.clear();
		await db.visits.clear();
	});
}
