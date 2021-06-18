// Avatar resolvers

export const getGravatar = (key: string, size: string) => `https://seccdn.libravatar.org/avatar/${key}?d=mp&s=${size}`

export const getRoboHash = (key: string, size: string) => `https://robohash.org/${key}?set=set4&size=${size}x${size}`
