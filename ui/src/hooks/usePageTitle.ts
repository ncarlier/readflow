import { useEffect } from 'react'

export const usePageTitle = (title = 'Readflow', subtitle?: string) => {
  useEffect(() => {
    document.title = subtitle ? subtitle : title
  }, [title, subtitle])
}
