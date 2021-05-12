const isDev = process.env.NODE_ENV === 'development'

const Wip = ({ children }) => isDev ? <>{ children }</> : null

export default Wip
