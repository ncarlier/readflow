import { useEffect } from 'react'

export default (title = 'Readflow', subtitle?: string) => {
  useEffect(() => {
    document.title = subtitle ? subtitle : title
  }, [title, subtitle])
}
