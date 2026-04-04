/**
 * Client types for monetary amounts serialized as integer minor units on the API
 * (see `backend/infra/money` and https://github.com/techforge-lat/money-as-integer).
 */

declare const minorUnitsBrand: unique symbol

/** Integer minor units for a currency (e.g. USD cents). */
export type MinorUnits = number & {
	readonly [minorUnitsBrand]: typeof minorUnitsBrand
}

const DEFAULT_DECIMAL_PLACES = 2

function assertWholeMinorUnits(n: number): asserts n is MinorUnits {
	if (!Number.isFinite(n)) {
		throw new TypeError("minor units must be a finite number")
	}
	if (!Number.isInteger(n)) {
		throw new TypeError("minor units must be an integer")
	}
	if (Math.abs(n) > Number.MAX_SAFE_INTEGER) {
		throw new RangeError("minor units exceed safe integer range")
	}
}

/** Coerces a value from the API (JSON number) into minor units. */
export function minorUnitsFromApi(n: number): MinorUnits {
	assertWholeMinorUnits(n)
	return n
}

/**
 * Parses a user-entered decimal amount into minor units using half-up rounding
 * (aligned with the default truncation mode in money-as-integer).
 */
export function parseDecimalToMinorUnits(
	raw: string,
	decimalPlaces = DEFAULT_DECIMAL_PLACES,
): MinorUnits | null {
	const normalized = raw.replace(/,/g, "").trim()
	if (normalized === "") {
		return minorUnitsFromApi(0)
	}
	const n = Number.parseFloat(normalized)
	if (!Number.isFinite(n)) {
		return null
	}
	const factor = 10 ** decimalPlaces
	const rounded = Math.round(n * factor)
	assertWholeMinorUnits(rounded)
	return rounded as MinorUnits
}

/** Formats minor units for display; `decimalPlaces` should match the currency (default 2). */
export function formatMinorUnits(
	minor: MinorUnits,
	currencyCode: string,
	decimalPlaces = DEFAULT_DECIMAL_PLACES,
): string {
	return new Intl.NumberFormat(undefined, {
		style: "currency",
		currency: currencyCode,
		minimumFractionDigits: decimalPlaces,
		maximumFractionDigits: decimalPlaces,
	}).format(minor / 10 ** decimalPlaces)
}
