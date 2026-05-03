/**
 * Reference data types — these match the JSON schema produced by the Go scraper.
 *
 * If you change the scraper's output format, update these types to match.
 * The TypeScript compiler will then flag every place in the app that needs
 * updating. That's the point of having types.
 */

// -- Reference data (read-only, bundled at build time from data/*.json) --

export interface DigitalAccess {
	/** Whether this system has an OverDrive/Libby collection */
	hasOverdrive: boolean;
	/** URL to the system's OverDrive page (e.g., "https://saclibrary.overdrive.com") */
	overdriveUrl: string;
	/** Whether this system offers online eCard registration */
	offersEcard: boolean;
	/** URL to the eCard registration page */
	ecardUrl: string;
	/** Free-text notes about eCard limitations (e.g., "30-day temporary, must visit to convert") */
	ecardNotes: string;
}

export interface LibrarySystem {
	/** Stable slug identifier (e.g., "sacramento-public-library") */
	id: string;
	/** Full display name (e.g., "Sacramento Public Library") */
	name: string;
	/** California county (e.g., "Sacramento") */
	county: string;
	/** City of the administrative HQ */
	city: string;
	/** Street address of the administrative HQ */
	address: string;
	/** Public phone number */
	phone: string;
	/** System's public website URL */
	website: string;
	/** Contact email */
	email: string;
	/** Online catalog URL, if known */
	catalogUrl: string;
	/** Library director's name */
	directorName: string;
	/**
	 * CLSA cooperative system + subgroup (e.g., "northnet-mvls").
	 * Format: "{system}-{subgroup}" where both are slugified.
	 */
	clsaSystem: string;
	/** Latitude of the administrative HQ */
	lat: number;
	/** Longitude of the administrative HQ */
	lng: number;
	/** Digital resource access information */
	digitalAccess: DigitalAccess;
	/** ISO 8601 date when this record was last updated by the scraper */
	lastUpdated: string;
}

export type OutletType = 'central' | 'branch' | 'bookmobile' | 'other';
export type BranchStatus = 'open' | 'closed' | 'temporarily-closed';

export interface Branch {
	/** Stable slug identifier */
	id: string;
	/** ID of the parent LibrarySystem */
	systemId: string;
	/** Full display name */
	name: string;
	/** Street address */
	address: string;
	/** City */
	city: string;
	/** ZIP code */
	zipCode: string;
	/** California county */
	county: string;
	/** Phone number */
	phone: string;
	/** Branch-specific web page, if different from system website */
	website: string;
	/** Latitude */
	lat: number;
	/** Longitude */
	lng: number;
	/** Type of outlet */
	outletType: OutletType;
	/** Current operational status */
	status: BranchStatus;
	/** ISO 8601 date when this record was last updated */
	lastUpdated: string;
}

export interface CooperativeSubgroup {
	/** Slug identifier (e.g., "mvls") */
	id: string;
	/** Display name (e.g., "Mountain Valley Library System") */
	name: string;
	/** IDs of LibrarySystems that belong to this subgroup */
	memberSystemIds: string[];
}

export interface CooperativeSystem {
	/** Slug identifier (e.g., "northnet") */
	id: string;
	/** Display name (e.g., "NorthNet Library System") */
	name: string;
	/** Subgroups within this cooperative */
	subgroups: CooperativeSubgroup[];
}

export type PartnershipType = 'catalog-sharing' | 'reciprocal-borrowing' | 'other';

export interface Partnership {
	/** Slug identifier */
	id: string;
	/** Display name (e.g., "Sacramento Valley Partners") */
	name: string;
	/** Human-readable description of what this partnership means for cardholders */
	description: string;
	/** IDs of LibrarySystems in this partnership */
	memberSystemIds: string[];
	/** What kind of partnership this is */
	type: PartnershipType;
}

// -- User data (stored in IndexedDB via Dexie, never sent to a server) --

export type CardStatus = 'active' | 'expired' | 'lost' | 'unknown';

export interface UserCard {
	/** Auto-generated UUID */
	id: string;
	/** ID of the LibrarySystem this card belongs to */
	systemId: string;
	/** Card number — optional, stored locally only, never transmitted */
	cardNumber: string;
	/** PIN — optional, stored locally only, never transmitted */
	pin: string;
	/** ISO 8601 date when the card was obtained */
	obtainedDate: string;
	/** ISO 8601 date when the card expires */
	expirationDate: string;
	/** Current status of the card */
	status: CardStatus;
	/** Whether this is a digital-only eCard (vs a physical card) */
	isDigitalOnly: boolean;
	/** Free-text notes */
	notes: string;
}

export interface BranchVisit {
	/** Auto-generated UUID */
	id: string;
	/** ID of the Branch that was visited */
	branchId: string;
	/** ISO 8601 date of the visit */
	visitDate: string;
	/** Free-text notes about the visit */
	notes: string;
}

// -- Aggregate types for convenience --

export interface ReferenceData {
	systems: LibrarySystem[];
	branches: Branch[];
	cooperatives: CooperativeSystem[];
	partnerships: Partnership[];
}
