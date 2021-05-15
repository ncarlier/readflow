import jwt_decode from "jwt-decode"

export const decodeToken = (token) => {
    const decoded = jwt_decode(token)
    // TODO validate signature!
    return decoded
}
