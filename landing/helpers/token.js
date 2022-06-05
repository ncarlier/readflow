import jwt_decode from 'jwt-decode'

/**
 * Decode JSON Web Token.
 * @param {string} token JWT base64 string
 * @returns {Object} decoded JWT
 */
export const decodeToken = (token) => {
    const decoded = jwt_decode(token)
    // TODO validate signature!
    return decoded
}
