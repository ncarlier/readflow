const isDev = process.env.NODE_ENV === 'development'

const Wip = ({children, placeholder}) => isDev ? <>{ children }</> : <>{ placeholder }</>

export default Wip
