import React, { useCallback, useState } from 'react'

import LinkIcon from '../../../components/LinkIcon'
import Loader from '../../../components/Loader'

const formats = [
  { value: 'html', label: 'HTML page', icon: 'code' },
  { value: 'offline', label: 'HTML page with images', icon: 'image' },
]

interface Props {
  download: (format: string) => Promise<any>
  onCancel: (e: any) => void
}

export default ({ download, onCancel }: Props) => {
  const [loading, setLoading] = useState(false)
  const handleDownloadArticle = useCallback(
    async (format: string) => {
      setLoading(true)
      try {
        await download(format)
      } finally {
        setLoading(false)
      }
    },
    [download]
  )

  if (loading) {
    return <Loader blur />
  }

  return (
    <ul>
      {formats.map((format) => (
        <li key={`format_${format.value}`}>
          <LinkIcon
            title="Download article as ..."
            icon={format.icon}
            onClick={() => handleDownloadArticle(format.value).then(onCancel)}
          >
            <span>{format.label}</span>
          </LinkIcon>
        </li>
      ))}
      <li>
        <LinkIcon title="Cancel" icon="cancel" onClick={onCancel}>
          <span>Cancel</span>
        </LinkIcon>
      </li>
    </ul>
  )
}
