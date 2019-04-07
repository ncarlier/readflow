import { useEffect } from "react"

export default (title: string, subtitle?: string) => {
  useEffect(() => {
    document.title = subtitle ? subtitle : title
  }, [title, subtitle])
}

