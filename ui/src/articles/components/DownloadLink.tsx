import React, { useCallback, useContext, useState } from 'react'

import LinkIcon from '../../components/LinkIcon'
import Loader from '../../components/Loader'
import { MessageContext } from '../../context/MessageContext'
import fetchAPI from '../../helpers/fetchAPI'
import { Article } from '../models'

interface Props {
  article: Article
  keyboard?: boolean
}

export default ({ article }: Props) => {
  const [loading, setLoading] = useState(false)
  const { showErrorMessage } = useContext(MessageContext)

  const download = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetchAPI(`/articles/${article.id}`, null, { method: 'GET' })
      if (res.ok) {
        const data = await res.blob()
        const href = window.URL.createObjectURL(data)
        const link = document.createElement('a')
        link.href = href
        link.setAttribute('download', `${article.title}.html`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      } else {
        const err = await res.json()
        throw new Error(err.detail || res.statusText)
      }
    } catch (err) {
      showErrorMessage(err.message)
    } finally {
      setLoading(false)
    }
  }, [article, showErrorMessage])

  if (loading) {
    return <Loader center />
  }

  return (
    <LinkIcon title="Download article" onClick={download} icon="download">
      <span>Download</span>
    </LinkIcon>
  )
}
