import { padStart, size, startsWith } from 'lodash'

/**
 * @typedef {Object} RGB
 * @property {number} r
 * @property {number} g
 * @property {number} b
 * @property {number} [a]
 */

/**
 * parse hex color to rgb object
 * @param hex
 * @return {RGB}
 */
export function parseHexColor(hex) {
    if (size(hex) < 6) {
        return { r: 0, g: 0, b: 0 }
    }
    if (startsWith(hex, '#')) {
        hex = hex.slice(1)
    }
    const bigint = parseInt(hex, 16)
    const r = (bigint >> 16) & 255
    const g = (bigint >> 8) & 255
    const b = bigint & 255
    return { r, g, b }
}